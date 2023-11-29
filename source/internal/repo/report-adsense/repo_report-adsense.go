package report_adsense

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os/exec"
	"path"
	"runtime"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/pkg/utility"
	"strings"
	"time"
)

type RepoReportAdsense interface {
	GetAllAccounts(refreshToken, accountAdsense string, try ...int64) (data []model.AdsenseAccount, err error)
	GetReports(inputs InputGetReports, try ...int64) (data []model.AdsenseReport, err error)
	GetReports2(inputs InputGetReports, try ...int64) (data []model.AdsenseReport, err error)
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindOneByQuery(query map[string]interface{}) (record *model.ReportAdsenseModel, err error)
	FindAllByQuery(query map[string]interface{}) (records []*model.ReportAdsenseModel, err error)
	FindAllByGroup(group ...string) (records []*model.ReportAdsenseModel, err error)
	Save(record *model.ReportAdsenseModel) (err error)
	SaveSlice(records []*model.ReportAdsenseModel) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportAdsenseModel, err error)
	Filter(input *InputFilter) (totalRecord int64, records []*model.ReportAffModel, recordTotal *model.ReportAffModel, err error)
}

type reportAdsenseRepo struct {
	Db         *gorm.DB
	Cache      caching.Cache
	Kafka      *kafka.Client
	totalRetry int64
	pkgPath    string
}

func NewReportAdsenseRepo(db *gorm.DB) *reportAdsenseRepo {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("No caller information")
	}
	pkgPath := path.Dir(filename)
	return &reportAdsenseRepo{
		Db:         db,
		pkgPath:    pkgPath,
		totalRetry: 3, // Số lần retry tối đa
	}
}

func (t *reportAdsenseRepo) GetAllAccounts(refreshToken, accountAdsense string, try ...int64) (data []model.AdsenseAccount, err error) {
	// Xử lý retry lại nếu có
	var Try int64
	if len(try) > 0 {
		Try = try[0] + 1
	}
	defer func() {
		if len(try) > 0 {
			if Try < t.totalRetry && err != nil {
				time.Sleep(3 * time.Second)
				data, err = t.GetAllAccounts(refreshToken, accountAdsense, Try)
			}
		}
	}()

	c := exec.Command("php", t.pkgPath+"/php/get_all_accounts.php",
		"--refreshToken="+refreshToken,
		"--accountAdsense="+accountAdsense,
	)
	//fmt.Println(c.String())
	out, err := c.Output()
	//fmt.Printf("\n out net: %+v \n", string(out))
	if err != nil {
		return
	}
	var resp model.ResponsePHPAdsense
	err = json.Unmarshal(out, &resp)
	if err != nil {
		return
	}
	if !resp.Status {
		err = errors.New(resp.Message)
		return
	}
	bData, err := json.Marshal(resp.Data)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(bData, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

type InputGetReports struct {
	AccountAdsense string
	RefreshToken   string
	Account        string
	StartDate      string
	EndDate        string
}

func (t *reportAdsenseRepo) GetReports(inputs InputGetReports, try ...int64) (data []model.AdsenseReport, err error) {
	// Xử lý retry lại nếu có
	var Try int64
	if len(try) > 0 {
		Try = try[0] + 1
	}
	defer func() {
		if len(try) > 0 {
			if Try < t.totalRetry {
				time.Sleep(3 * time.Second)
				data, err = t.GetReports(inputs, Try)
			}
		}
	}()

	c := exec.Command("php", t.pkgPath+"/php/get_reports.php",
		"--accountAdsense="+inputs.AccountAdsense,
		"--refreshToken="+inputs.RefreshToken,
		"--account="+inputs.Account,
		"--startDate="+inputs.StartDate,
		"--endDate="+inputs.EndDate,
	)
	//fmt.Println(c.String())
	out, err := c.Output()
	//fmt.Printf("\n out net: %+v \n", string(out))
	if err != nil {
		return
	}
	var resp model.ResponsePHPAdsense
	err = json.Unmarshal(out, &resp)
	if err != nil {
		return
	}
	if !resp.Status {
		err = errors.New(resp.Message)
		return
	}
	bData, err := json.Marshal(resp.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bData, &data)
	if err != nil {
		return
	}
	return
}

func (t *reportAdsenseRepo) GetReports2(inputs InputGetReports, try ...int64) (data []model.AdsenseReport, err error) {
	// Xử lý retry lại nếu có
	var Try int64
	if len(try) > 0 {
		Try = try[0] + 1
	}
	defer func() {
		if len(try) > 0 {
			if Try < t.totalRetry {
				time.Sleep(3 * time.Second)
				data, err = t.GetReports(inputs, Try)
			}
		}
	}()
	dates := utility.GetAllDates(inputs.StartDate, inputs.EndDate)
	for _, date := range dates {
		c := exec.Command("php", t.pkgPath+"/php/get_reports2.php",
			"--accountAdsense="+inputs.AccountAdsense,
			"--refreshToken="+inputs.RefreshToken,
			"--account="+inputs.Account,
			"--startDate="+date.Format("2006-01-02"),
			"--endDate="+date.Format("2006-01-02"),
		)
		fmt.Println(c.String())
		out, err := c.Output()
		//fmt.Printf("\n out net: %+v \n", string(out))
		if err != nil {
			return nil, err
		}
		var resp model.ResponsePHPAdsense
		err = json.Unmarshal(out, &resp)
		if err != nil {
			return nil, err
		}
		if !resp.Status {
			err = errors.New(resp.Message)
			return nil, err
		}
		bData, err := json.Marshal(resp.Data)
		if err != nil {
			return nil, err
		}
		var output []model.AdsenseReport
		err = json.Unmarshal(bData, &output)
		if err != nil {
			return nil, err
		}
		data = append(data, output...)
	}
	return
}

type InputIsExists struct {
}

func (t *reportAdsenseRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportAffPixelModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportAdsenseRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportAdsenseModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportAdsenseRepo) FindOneByQuery(query map[string]interface{}) (record *model.ReportAdsenseModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportAdsenseRepo) FindAllByQuery(query map[string]interface{}) (records []*model.ReportAdsenseModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&records).Error
	return
}

func (t *reportAdsenseRepo) FindAllByGroup(group ...string) (records []*model.ReportAdsenseModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Group(strings.Join(group, ",")).
		Find(&records).Error
	return
}

func (t *reportAdsenseRepo) Save(record *model.ReportAdsenseModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportAdsenseRepo) SaveSlice(records []*model.ReportAdsenseModel) (err error) {
	if len(records) == 0 {
		return
	}
	size := 500
	var j int
	for i := 0; i < len(records); i += size {
		j += size
		if j > len(records) {
			j = len(records)
		}
		var recordsSlice []*model.ReportAdsenseModel
		recordsSlice = records[i:j]
		err = t.Db.
			//Debug().
			Save(&recordsSlice).Error
	}
	return
}

func (t *reportAdsenseRepo) FindByDayForReportAff(day string) (records []*model.ReportAdsenseModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportAdsenseModel{}).
		Where("date LIKE '" + day + "'").
		Find(&records).Error
	return
}

type InputFilter struct {
	UserID int64

	Search string
	Offset int
	Limit  int
	Order  string

	StartDate      string
	EndDate        string
	Campaigns      interface{}
	SectionID      interface{}
	TrafficSources interface{}
	GroupBy        interface{}
}

func (t *reportAdsenseRepo) Filter(input *InputFilter) (totalRecord int64, records []*model.ReportAffModel, recordTotal *model.ReportAffModel, err error) {
	err = t.Db.
		//Debug().
		Select(
			"date",
			"traffic_source",
			"section_id", "section_name",
			"campaign_id",
			"campaign_name",
			"redirect_id",
			"partner",
			"style_id",
			"SUM(system_traffic) AS system_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
			"SUM(click_adsense) AS click_adsense",
			"SUM(cost) AS cost",
			"SUM(revenue) AS revenue",
			"SUM(estimated_revenue) AS estimated_revenue",
			"SUM(pre_estimated_revenue) AS pre_estimated_revenue",
			"SUM(supply_click) AS supply_click",
			"SUM(supply_conversions) AS supply_conversions",
		).
		Where("partner = 'adsense'").
		Scopes(
			t.setFilterCondition(input),
			t.setFilterGroupBy(input.GroupBy),
			t.setFilterDate(input.StartDate, input.EndDate),
		).
		Model(&model.ReportAffModel{}).
		Count(&totalRecord).
		Order(input.Order).
		Scopes(mysql.Paginate(mysql.Deps{Offset: input.Offset, Limit: input.Limit})).
		Find(&records).Error
	err = t.Db.
		//Debug().
		Select(
			"SUM(system_traffic) AS system_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
			"SUM(click_adsense) AS click_adsense",
			"SUM(cost) AS cost",
			"SUM(revenue) AS revenue",
			"SUM(supply_click) AS supply_click",
			"SUM(supply_conversions) AS supply_conversions",
			"SUM(estimated_revenue) AS estimated_revenue",
			"SUM(pre_estimated_revenue) AS pre_estimated_revenue",
		).
		Scopes(
			t.setFilterCondition(input),
			t.setFilterDate(input.StartDate, input.EndDate),
		).
		Where("partner = 'adsense'").
		Model(&model.ReportAffModel{}).
		Find(&recordTotal).Error
	return
}

func (t *reportAdsenseRepo) setFilterGroupBy(groupBy interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if groupBy == nil {
			return db
		}
		switch groupBy.(type) {
		case string:
			if groupBy != "" {
				group := ""
				if fmt.Sprintf("%v", groupBy) == "campaign" {
					group = "campaign_id"
				} else {
					group = fmt.Sprintf("%v", groupBy)
				}
				return db.Group(group)
			}
		case []interface{}:
			listGroup := groupBy.([]interface{})
			var groups []string
			for _, group := range listGroup {
				if fmt.Sprintf("%v", group) == "campaign" {
					groups = append(groups, fmt.Sprintf("%v", "campaign_id"))
				} else {
					groups = append(groups, fmt.Sprintf("%v", group))
				}
			}
			return db.Group(strings.Join(groups, ","))
		}
		return db
	}
}

func (t *reportAdsenseRepo) setFilterCondition(input *InputFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var condition = make(map[string]interface{})
		//=> convert input sang dữ liệu condition để query
		if input.Campaigns != nil {
			condition["campaign_id"] = input.Campaigns
		}
		if input.TrafficSources != nil {
			condition["traffic_source"] = input.TrafficSources
		}
		if input.SectionID != nil {
			condition["section_id"] = input.SectionID
		}
		return db.Where(condition)
	}
}

func (t *reportAdsenseRepo) setFilterDate(startDate, endDate string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var listDate []string
		layout := "2006-01-02"
		start, _ := time.Parse(layout, startDate)
		listDate = append(listDate, start.Format(layout))
		if startDate != endDate {
			for {
				// Add thêm 1 ngày
				start = start.AddDate(0, 0, 1)
				listDate = append(listDate, start.Format(layout))
				if start.Format(layout) == endDate {
					break
				}
			}
		}
		return db.Where("date in ?", listDate)
	}
}
