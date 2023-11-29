package report_bodis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/repo"
	reportbodisRepo "source/internal/repo/report-bodis"
	reportBodisTrafficRepo "source/internal/repo/report-bodis-traffic"
	"strings"
	"time"
)

type UsecaseReportBodis interface {
	Cronjob() (err error)
	Filter(payload *dto.PayloadReportBodisFilter) (reportBodis []model.ReportBodisModel, err error)
	GetAllSubID() []string
	GetReportTraffic(from, to, page string) (output []byte, err error)
	HandlerReportTraffic(inputs []InputHandlerReportTraffic) (err error)
}

type reportBodisUsecase struct {
	repos *repo.Repositories
	Trans *lang.Translation
}

func NewReportBodisUsecase(repos *repo.Repositories, trans *lang.Translation) *reportBodisUsecase {
	return &reportBodisUsecase{repos: repos, Trans: trans}
}

func (t *reportBodisUsecase) Cronjob() (err error) {
	loc, _ := time.LoadLocation("America/New_York")
	// startDate, err := time.Parse("2006-01-02", "2022-07-21")
	// if err != nil {
	// 	return
	// }
	// endDate, err := time.Parse("2006-01-02", "2022-07-27")
	// if err != nil {
	// 	return
	// }
	return t.repos.ReportBodis.CronjobReport(&reportbodisRepo.DateRange{
		StartDate: time.Now().AddDate(0, 0, -1).In(loc),
		EndDate:   time.Now().In(loc),
	})
}

func (t *reportBodisUsecase) Filter(payload *dto.PayloadReportBodisFilter) (reportBodis []model.ReportBodisModel, err error) {
	if err = payload.Validate(); err != nil {
		return
	}
	// chuyển đổi và validate dữ liệu time
	var startDate, endDate time.Time
	if payload.StartDate != "" && payload.EndDate != "" {
		if startDate, err = time.Parse("2006-01-02", payload.StartDate); err != nil {
			return
		}
		if endDate, err = time.Parse("2006-01-02", payload.EndDate); err != nil {
			return
		}
		if startDate.After(endDate) {
			err = errors.New("the start date must be the previous day")
			return
		}
	}
	// gọi vào repo để thực hiện filter
	reportBodis, err = t.repos.ReportBodis.FindByFilter(&reportbodisRepo.FilterInput{
		Condition: payload.ToCondition(),
		OrderBy:   payload.OrderBy,
		//GroupBy:   payload.GroupBy,
		SubID:     payload.SubID,
		StartDate: startDate,
		EndDate:   endDate,
		Order:     "id DESC",
	})
	return
}

func (t *reportBodisUsecase) GetAllSubID() (SubIDs []string) {
	SubIDs = t.repos.ReportBodis.AllSubIDs()
	return
}

func (t *reportBodisUsecase) GetReportTraffic(from, to, page string) (output []byte, err error) {
	url := "https://api.bodis.com/v2/parking/search?page=1"
	//fmt.Println(url)

	//=> Dùng Colly request lên request lên API của a Xuân
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
		r.Headers.Set("Authorization", "Bearer rFgfcLqRQPSIM01ExcEMWcUfJKPK7MXV")
	})
	// Get response
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("Response:", string(r.Body))
		output = r.Body
	})
	// Check Error
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
	})
	type request struct {
		Filter struct {
			Type      string   `json:"type"`
			DateRange []string `json:"date_range"`
		} `json:"filter"`
		PerPage   string `json:"per_page"`
		Page      string `json:"page"`
		OrderBy   string `json:"order_by"`
		SortOrder string `json:"sort_order"`
	}
	requestData, _ := json.Marshal(request{
		Filter: struct {
			Type      string   `json:"type"`
			DateRange []string `json:"date_range"`
		}{
			Type:      "click",
			DateRange: []string{from, to},
		},
		PerPage:   "250",
		Page:      page,
		OrderBy:   "server_datetime",
		SortOrder: "asc",
	})
	fmt.Println("Request: ", string(requestData))
	err = c.PostRaw(url, requestData)
	if err != nil {
		return
	}

	c.Wait()

	return
}

type InputHandlerReportTraffic struct {
	Time             string
	VisitID          string
	DomainName       string
	IpAddress        string
	Type             string
	CountryID        int64
	PageQuery        string
	SubIds           string
	Clicks           int64
	EstimatedRevenue float64
}

func (t *reportBodisUsecase) HandlerReportTraffic(inputs []InputHandlerReportTraffic) (err error) {
	// Từ input build lấy messages cho druid của a Xuân
	_, records, err := t.handlerInput(inputs)
	if err != nil {
		return
	}

	// Từ messages gửi lên vào kafka của a Xuân để update report
	//err = t.repos.ReportBodisTraffic.PushMessage(messages)
	//if err != nil {
	//	return
	//}

	// Nếu như gửi message k lỗi tiến hành save cho DB theo record mới
	err = t.repos.ReportBodisTraffic.SaveSlice(records)
	if err != nil {
		return
	}
	return
}

func (t *reportBodisUsecase) handlerInput(inputs []InputHandlerReportTraffic) (messages []model.TrackingAdsMessage, records []*model.ReportBodisTrafficModel, err error) {

	// Từ inputs xử lý tổng hợp lại data clicks, revenue vào map
	for _, input := range inputs {
		//=> Parse subIds để lấy trafficSource, campaign, sectionID
		trafficSource := "unknow"
		campaign := "unknow"
		sectionID := "unknow"

		subIds := strings.Split(input.SubIds, "&")
		for _, subID := range subIds {
			if subID == "" {
				continue
			}
			//=> subID sẽ có dạng subid_1=Outbrain, subid_2=campaignID-00bae50ca86b39d8f7c051e8c232534e54 parse tiếp để lấy dữ liệu cần
			detailSubID := strings.Split(subID, "=")
			if detailSubID[0] == "subid_1" {
				trafficSource = detailSubID[1]
			} else {
				data := strings.Split(detailSubID[1], "-")
				if data[0] == "campaignID" {
					campaign = data[1]
				} else if data[0] == "selectionID" {
					sectionID = data[1]
				}
			}
		}

		// Từ các data build ra model cho DB
		record := model.ReportBodisTrafficModel{
			Time:             input.Time,
			VisitID:          input.VisitID,
			DomainName:       input.DomainName,
			IpAddress:        input.IpAddress,
			Type:             input.Type,
			CountryID:        input.CountryID,
			PageQuery:        input.PageQuery,
			TrafficSource:    trafficSource,
			Campaign:         campaign,
			SelectionID:      sectionID,
			Clicks:           input.Clicks,
			EstimatedRevenue: input.EstimatedRevenue,
		}

		// Build Inputs để check exist của report
		inputExist := reportBodisTrafficRepo.InputIsExists{
			VisitID:       input.VisitID,
			Time:          input.Time,
			TrafficSource: trafficSource,
			Campaign:      campaign,
			SelectionID:   sectionID,
		}

		// Check exist
		isExist := t.repos.ReportBodisTraffic.IsExists(&inputExist)
		if !isExist {
			// Nếu như chưa tồn tại build lại từ record thành request luôn không cần xử lý gì thêm
			message := record.ToMessageKafka()
			messages = append(messages, message)
			records = append(records, &record)
		}
	}
	return
}
