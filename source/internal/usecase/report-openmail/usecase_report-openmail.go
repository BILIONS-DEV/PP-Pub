package report_openmail

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/mitchellh/mapstructure"
	url2 "net/url"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	"strconv"
	"strings"
	"time"
)

type UseCaseReportOpenMail interface {
	GetReports(url string, query ...string) (reports []map[string]interface{}, err error)
	GetReportHourly(query ...string) (records []*model.ReportOpenMailSubIdModel, err error)
	HandlerReportSubID(records []*model.ReportOpenMailSubIdModel) (err error)
}

func NewReportOpenMailUC(repos *repo.Repositories) *openMailUC {
	return &openMailUC{repos: repos, authKey: "ahgDbcMLaeY8gSONa7am"}
}

type openMailUC struct {
	authKey string
	repos   *repo.Repositories
}

func (t *openMailUC) GetReports(url string, query ...string) (reports []map[string]interface{}, err error) {
	//=> Gán query vào values của url
	queryString := strings.Join(query, "&")
	values, err := url2.ParseQuery(queryString)
	if err != nil {
		return
	}

	//=> Thêm Auth Key của để check report của mình
	values.Set("auth_key", t.authKey)

	//=> Tạo link get Reports
	urlReport := url + "?" + values.Encode()
	fmt.Println(urlReport)

	//=> Tạo response động cho dữ liệu trả về
	var responses [][]interface{}

	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	// Check Response
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &responses)
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
	err = c.Visit(urlReport)
	if err != nil {
		return
	}
	c.Wait()

	//=> Nếu response rỗng thì bỏ qua
	if len(responses) == 0 {
		err = errors.New("no response")
		return
	}

	//=> Tạo mảng chứa các key
	var keys []string

	//=> Response nhận về sẽ là một mảng json tiến hành for để xử lý
	for indexResponse, res := range responses {
		if indexResponse == 0 {
			//=> Trong response nhận về sẽ là một mảng json trong đó mảng đầu tiên là để tên định danh cho các dữ liệu ở vị trí tương ứng
			for _, key := range res {
				keys = append(keys, fmt.Sprintf("%v", key))
			}
		} else {
			//=> Các mảng tiếp theo sẽ là dữ liệu tương ứng với vị trí key đã định danh ở mảng đầu
			report := make(map[string]interface{})
			for indexValue, value := range res {
				if keys != nil && len(keys) > indexValue {
					report[keys[indexValue]] = value
				}
			}
			reports = append(reports, report)
		}
	}
	return
}

func (t *openMailUC) GetReportHourly(query ...string) (records []*model.ReportOpenMailSubIdModel, err error) {
	response := make(map[string]interface{})
	url := "https://reports.openmail.com/v2/subid_estimated_hourly.json"
	//=> Gán query vào values của url
	queryString := strings.Join(query, "&")
	values, err := url2.ParseQuery(queryString)
	if err != nil {
		return
	}

	//=> Thêm Auth Key của để check report của mình
	values.Set("auth_key", t.authKey)

	//=> Tạo link get Reports
	urlReport := url + "?" + values.Encode()
	fmt.Println(urlReport)

	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	// Get reports từ response
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &response)
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
	err = c.Visit(urlReport)
	if err != nil {
		return
	}
	c.Wait()

	//=> Nếu reports rỗng thì trả về lỗi no reports
	if len(response) == 0 {
		err = errors.New("no response")
		return
	}

	// Chuyển từ output map sang struct dto của report
	var reports dto.ReportOpenMailSubIdEstimatedHourly
	err = mapstructure.Decode(response, &reports)
	// fmt.Printf("%+v \n", reports)

	// Chuyển từ dto sang model để truyền vào usecase xử lý
	for _, subID := range reports.SubIds {
		record := subID.ToModel()
		records = append(records, &record)
	}
	return
}

func (t *openMailUC) HandlerReportSubID(records []*model.ReportOpenMailSubIdModel) (err error) {
	// Chạy vòng lặp xử lý cho từng record
	for _, record := range records {
		// Check exist
		// Nếu như record đã tồn tại tạo query để lấy recordOld
		query := make(map[string]interface{})
		query["utc_hour"] = record.UtcHour
		query["campaign"] = record.Campaign
		query["sub_id"] = record.SubID

		recordOld, err := t.repos.ReportOpenMailSubID.FindByQuery(query)
		if err != nil {
			continue
		}

		if recordOld.ID != 0 {
			//=> Gán ID cho recordNew để khi save update
			record.ID = recordOld.ID
			record.SectionId = recordOld.SectionId
			record.SectionName = recordOld.SectionName
			record.PublisherID = recordOld.PublisherID
			record.PublisherName = recordOld.PublisherName
			record.AdID = recordOld.AdID
			record.RedirectID = recordOld.RedirectID
		}

		if record.RedirectID == 0 {
			splitSubID := strings.Split(record.SubID, ":")
			if len(splitSubID) < 2 {
				// err = errors.New("subID malformed")
				continue
			}
			subID1 := splitSubID[0] // subID1 trafficSource_campaignID || uid
			if subID1 == "uid" {
				uid := splitSubID[1]
				// Nếu subID1 là clickID thì subID 2 chứa clickID để tìm map với bảng report_aff_pixel
				query := make(map[string]interface{})
				query["uid"] = uid
				reportPixel, err := t.repos.ReportAffPixel.FindOneByQuery(query)
				if err != nil {
					continue
				}
				redirectID, _ := strconv.ParseInt(reportPixel.RedirectID, 10, 64)
				record.SectionId = reportPixel.SectionId
				record.SectionName = reportPixel.SectionName
				record.PublisherID = reportPixel.PublisherID
				record.PublisherName = reportPixel.PublisherName
				record.AdID = reportPixel.AdID
				record.RedirectID = redirectID
			}
		}
	}

	// Gửi messages lên kafka topic của a Xuân để update vào druid report
	//err = t.repos.ReportOpenMailSubID.PushMessage(messages)
	//if err != nil {
	//	return
	//}

	// Lưu lại record mới vào DB
	err = t.repos.ReportOpenMailSubID.SaveSlice(records)
	if err != nil {
		return
	}

	return
}

func (t *openMailUC) buildMessageExist(record *model.ReportOpenMailSubIdModel) (message model.TrackingAdsMessage, err error) {
	// Nếu như record đã tồn tại tạo query để lấy recordOld
	query := make(map[string]interface{})
	query["utc_hour"] = record.UtcHour
	query["campaign"] = record.Campaign
	query["sub_id"] = record.SubID

	recordOld, err := t.repos.ReportOpenMailSubID.FindByQuery(query)
	if err != nil {
		return
	}
	//=> Gán ID cho recordNew để khi save update
	record.ID = recordOld.ID

	// Tạo request gửi lên cho api của a Xuân
	message.Time = time.Now().UTC().Format("2006-01-02T15:04:05")
	// Từ subID lấy được trafficSource, campaign, sectionID
	//=> SubID nhận vào sẽ có dạng trafficSource_campaignID:sectionID
	splitSubID := strings.Split(record.SubID, "_")
	if len(splitSubID) < 2 {
		err = errors.New("sub_id malformed")
		return
	}
	trafficSource := splitSubID[0]
	splitSubID2 := strings.Split(splitSubID[1], ":")
	if len(splitSubID2) < 2 {
		err = errors.New("sub_id malformed")
		return
	}
	campaign := splitSubID2[0]
	sectionID := splitSubID2[1]

	message.TrafficSource = trafficSource
	message.Campaign = campaign
	message.SelectionId = sectionID

	/*
		=> Vì druid nhận dữ liệu sẽ cộng tiếp vào các row đang có vậy nên sẽ so sánh record mới nhận được và recordOld
		Lấy phần chênh lệch giữa recordOld và recordNew gửi lên api của a Xuân để cộng vào report trên Druid
	*/
	isChange := false
	if recordOld.Clicks != record.Clicks {
		isChange = true
		clicks := record.Clicks - recordOld.Clicks
		message.Click = int(clicks)
	}
	//if recordOld.Searches != record.Searches {
	//	isChange = true
	//	searches := record.Searches - recordOld.Searches
	//	message.Impressions = searches
	//}
	if recordOld.EstimatedRevenue != record.EstimatedRevenue {
		isChange = true
		estimatedRevenue := record.EstimatedRevenue - recordOld.EstimatedRevenue
		message.GrossRevenue = estimatedRevenue
	}

	if !isChange {
		err = errors.New("no change")
		return
	}
	return
}

func (t *openMailUC) buildMessageNoExist(record *model.ReportOpenMailSubIdModel) (message model.TrackingAdsMessage, err error) {
	// Tạo request gửi lên cho api của a Xuân
	message.Time = time.Now().UTC().Format("2006-01-02T15:04:05")
	// Từ subID lấy được trafficSource, campaign, sectionID
	splitSubID := strings.Split(record.SubID, "_")
	if len(splitSubID) < 2 {
		err = errors.New("sub_id malformed")
		return
	}
	trafficSource := splitSubID[0]
	splitSubID2 := strings.Split(splitSubID[1], ":")
	if len(splitSubID2) < 2 {
		err = errors.New("sub_id malformed")
		return
	}
	campaign := splitSubID2[0]
	sectionID := splitSubID2[1]

	message.TrafficSource = trafficSource
	message.Campaign = campaign
	message.SelectionId = sectionID
	message.Click = int(record.Clicks)
	//message.Searches = record.Searches
	message.GrossRevenue = record.EstimatedRevenue
	return
}
