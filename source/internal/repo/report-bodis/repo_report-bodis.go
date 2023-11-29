package report_bodis

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"source/infrastructure/caching"
	"strconv"
	"strings"
	"time"

	// "source/infrastructure/mysql"
	"source/internal/entity/model"
)

type RepoReportBodis interface {
	// DB() *gorm.DB
	CronjobReport(DateRange *DateRange) (err error)
	FindByFilter(input *FilterInput) (reportBodis []model.ReportBodisModel, err error)
	AllSubIDs() (reportBodis []string)
	// Save(record *model.CronjobModel) (err error)
}

type reportBodisRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewReportBodisRepo(db *gorm.DB) *reportBodisRepo {
	return &reportBodisRepo{Db: db}
}

type FilterInput struct {
	// filter
	Condition map[string]interface{}
	StartDate time.Time
	EndDate   time.Time
	Search    string
	GroupBy   string
	OrderBy   string
	SubID     string
	// order and pagination
	Order string
	Page  int
}

func (t *reportBodisRepo) FindByFilter(input *FilterInput) (reportBodis []model.ReportBodisModel, err error) {
	t.Db.
		Debug().
		// Where("user_id = ? AND app = 'FE' AND creator_id = ?", userID, userID).
		Where(input.Condition).
		Scopes(
			t.setFilterQueryOrderBy(input.OrderBy),
			t.setFilterQueryGroupBy(input.GroupBy),
			t.setFilterQueryDate(input.StartDate, input.EndDate),
		).
		Model(&reportBodis).
		Order(input.Order).
		// Scopes(mysql.Paginate(mysql.Deps{Page: input.Page})).
		Select("*").Find(&reportBodis)
	return
}

func (t *reportBodisRepo) setFilterQueryOrderBy(OrderBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if OrderBy == "" {
			return db
		}
		return db.Order(OrderBy + " DESC")
	}
}

func (t *reportBodisRepo) setFilterQueryGroupBy(GroupBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if GroupBy == "" {
			return db
		}
		return db.Group(GroupBy)
	}
}

func (t *reportBodisRepo) setFilterQueryDate(startDate, endDate time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if startDate.IsZero() || endDate.IsZero() {
			return db
		}
		return db.Where("date >= ? AND date <= ?", startDate, endDate.AddDate(0, 0, 1))
	}
}

type PayloadCronjob struct {
	ReportType string              `json:"report_type"`
	Page       string              `json:"page"`
	PerPage    string              `json:"per_page"`
	SortOrder  string              `json:"sort_order"`
	Filter     map[string][]string `json:"filter"`
}

type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
type ResponseReport struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Metrics Metrics     `json:"metrics"`
}
type Metrics struct {
	current_page int                  `json:"current_page"`
	Data         []model.ResponseData `json:"data"`
}

func (t *reportBodisRepo) CronjobReport(DateRange *DateRange) (err error) {
	Dates := t.ArrayDate(DateRange)
	if len(Dates) == 0 {
		return
	}
	for _, date := range Dates {
		// chỉ sync report today, yesterday
		if date.Format("2006-01-02") > time.Now().Format("2006-01-02") {
			continue
		}
		// if date.Format("2006-01-02") < time.Now().AddDate(0, 0, -1).Format("2006-01-02") {
		// 	fmt.Printf("%+v\n", date.Format("2006-01-02"))
		// 	continue
		// }

		url := "https://api.bodis.com/v2/reports/search"
		method := "POST"
		// loc, _ := time.LoadLocation("America/New_York")
		payloadCronjob := PayloadCronjob{
			ReportType: "subid",
			Page:       "1",
			PerPage:    "250",
			SortOrder:  "desc",
			Filter: map[string][]string{
				"date_range": []string{
					date.Format("2006-01-02"),
					date.Format("2006-01-02"),
				},
			},
		}
		payloadByte, err := json.Marshal(payloadCronjob)
		payload := strings.NewReader(string(payloadByte))
		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return err
		}
		req.Header.Add("Authorization", "Bearer rFgfcLqRQPSIM01ExcEMWcUfJKPK7MXV")
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		var response ResponseReport
		err = json.Unmarshal(body, &response)
		if err != nil {
			return err
		}
		if response.Message != "" && response.Errors != nil {
			return errors.New(response.Message)
		}
		err = t.Save(response.Metrics.Data, date)
		if err != nil {
			return err
		}
	}
	time.Sleep(15 * time.Second)
	return
}

func (t *reportBodisRepo) Save(records []model.ResponseData, date time.Time) (err error) {
	loc, _ := time.LoadLocation("America/New_York")
	for _, value := range records {
		// check is exist date
		var report model.ReportBodisModel
		err = t.Db.
			Model(&model.ReportBodisModel{}).
			Where("subid = ? and date = ?", value.Subid, date.In(loc).Format("2006-01-02")).
			Find(&report).Error
		if err != nil {
			return
		}
		if report.ID != 0 {
			// update
			row := t.makeRow(value)
			row.Date = date.In(loc)
			row.ID = report.ID
			err = t.Db.Updates(&row).Error
		} else {
			// insert
			row := t.makeRow(value)
			row.Date = date.In(loc)
			err = t.Db.Create(&row).Error
		}
	}
	fmt.Println(" ")
	return
}

func (t *reportBodisRepo) makeRow(res model.ResponseData) (row model.ReportBodisModel) {
	row.Subid = fmt.Sprintf("%v", res.Subid)
	switch res.Subid.(type) {
	case string:
	case int:
		row.Subid = fmt.Sprintf("%v", res.Subid)
		break
	case float64:
		row.Subid = strconv.FormatInt(int64(res.Subid.(float64)), 10)
	}
	row.Visits = res.Visits
	row.LandingPageVisits = res.LandingPageVisits
	row.Clicks = res.Clicks
	row.CreditedRevenue = res.CreditedRevenue
	row.LandingPageCreditedRevenue = res.LandingPageCreditedRevenue
	row.ZeroclickCreditedRevenue = res.ZeroclickCreditedRevenue
	row.Ctr = res.Ctr
	row.Epc = res.Epc
	row.Rpm = res.Rpm
	row.LandingPageRpm = res.LandingPageRpm
	row.ZeroclickRpm = res.ZeroclickRpm
	row.ClicksSpamRatio = res.ClicksSpamRatio
	row.IsFinalized = res.IsFinalized
	return
}

func (t *reportBodisRepo) ArrayDate(dateRange *DateRange) (dates []time.Time) {
	if dateRange.StartDate.IsZero() || dateRange.EndDate.IsZero() {
		return
	}
	for i := 0; dateRange.StartDate.AddDate(0, 0, i).Format("20060102") <= dateRange.EndDate.Format("20060102"); i++ {
		dates = append(dates, dateRange.StartDate.AddDate(0, 0, i))
	}
	return
}

func (t *reportBodisRepo) AllSubIDs() (SubIDs []string) {
	var reports []model.ReportBodisModel
	t.Db.Select("subid").
		Model(&model.ReportBodisModel{}).
		Group("subid").
		Find(&reports)
	for _, value := range reports {
		SubIDs = append(SubIDs, value.Subid)
	}
	return
}
