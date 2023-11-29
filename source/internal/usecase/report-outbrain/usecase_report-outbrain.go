package report_outbrain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	report_outbrain "source/internal/repo/report-outbrain"
	"source/pkg/logger"
	"source/pkg/utility"
	"strings"
	"time"
)

type UseCaseReportOutBrain interface {
	GetCampaignIDs(marketerID string) (campaigns dto.ResponseCampaignOutBrain, err error)
	GetMarketerIDs() (marketerIDS []string, err error)
	GetConversionIDs(marketerID string) (conversionIDs []string, err error)
	GetSectionByHour(marketerID string, campaignID string, conversionIDS []string) (output []byte, err error)
	HandlerReport(records []model.ReportOutBrainModel) (err error)
	UpdateCampaignOutBrain(campaignID string, inputs InputUpdateCampaign) (output []byte, err error)
}

func NewReportOutBrainUC(repos *repo.Repositories) *outBrainUC {
	return &outBrainUC{
		repos: repos,
	}
}

type outBrainUC struct {
	repos *repo.Repositories
}

// Xử lý report OutBrain
func (t *outBrainUC) HandlerReport(records []model.ReportOutBrainModel) (err error) {
	//fmt.Printf("%+v \n", records)
	// Từ input build lấy messages cho druid của a Xuân
	_, recordNews, err := t.handlerInput(records)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Từ messages gửi lên vào kafka của a Xuân để update report
	//err = t.repos.ReportOutBrain.PushMessage(messages)
	//if err != nil {
	//	logger.Error(err.Error())
	//	return
	//}

	// Nếu như gửi message k lỗi tiến hành save cho DB theo record mới
	err = t.repos.ReportOutBrain.SaveSlice(recordNews)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

func (t *outBrainUC) handlerInput(records []model.ReportOutBrainModel) (messages []model.TrackingAdsMessage, recordNews []*model.ReportOutBrainModel, err error) {
	for _, record := range records {
		// Tạo một biến mới tránh trừng địa chỉ biến vòng lặp
		recordNew := record
		// Từ record tạo input kiểm tra xem record đã có trong db chưa
		isExist := t.repos.ReportOutBrain.IsExists(&report_outbrain.InputIsExists{
			MarketerID: recordNew.MarketerID,
			CampaignID: recordNew.CampaignID,
			SectionID:  recordNew.SectionID,
			Time:       recordNew.Time,
		})

		// => Xử lý message theo 2 trường hợp có và chưa có trong DB
		if !isExist {
			// Nếu chưa có thì tạo mới message và thêm list recordNew để tạo mới
			messages = append(messages, recordNew.ToMessageKafka())
			recordNews = append(recordNews, &recordNew)
		} else {
			// Nếu đã có thì tạo query lấy ra record cũ
			query := make(map[string]interface{})
			query["marketer_id"] = recordNew.MarketerID
			query["campaign_id"] = recordNew.CampaignID
			query["section_id"] = recordNew.SectionID
			query["time"] = recordNew.Time
			recordOld, err := t.repos.ReportOutBrain.FindByQuery(query)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if !recordOld.IsFound() {
				continue
			}

			// So sánh spend recordOld và recordNew
			amount := recordNew.Spend - recordOld.Spend
			// Từ recordNew build ra message và đổi lại phần chênh của amount để push lên druid cộng thêm
			message := recordNew.ToMessageKafka()
			message.Amount = amount
			messages = append(messages, message)

			// Gán lại ID của recordOld để khi save gorm nhận đây là update
			recordNew.ID = recordOld.ID
			recordNews = append(recordNews, &recordNew)

		}
	}
	return
}

// Get toàn bộ Campaign đang chạy của report theo ngày đó - date format 2022-08-17 (yyyy-mm-dd)
func (t *outBrainUC) GetCampaignIDs(marketerID string) (campaigns dto.ResponseCampaignOutBrain, err error) {

	// Lấy toàn bộ các campaign có chạy trong 2 ngày hôm nay và hôm qua
	loc, _ := time.LoadLocation("America/New_York") // Report OutBrain dùng theo time UTC -5 (nhưng check theo report thì đang trả về là UTC -4)
	from := time.Now().In(loc).AddDate(0, 0, -1)
	to := time.Now().In(loc)

	url := "https://api.outbrain.com/amplify/v0.1/reports/marketers/" + marketerID + "/campaigns"
	//=> Thêm Auth Key của để check report của mình
	values := url2.Values{}
	values.Set("from", from.Format("2006-01-02"))
	values.Set("to", to.Format("2006-01-02"))

	//=> Tạo link get Reports
	urlReport := url + "?" + values.Encode()
	fmt.Println(urlReport)

	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("OB-TOKEN-V1", t.repos.ReportOutBrain.GetToken())
	})
	// Check Response
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &campaigns)
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
	c.Wait()
	// Sau khi đợi xử lý xong nếu có lỗi thì return luôn
	if err != nil {
		return
	}
	var listCampaignID []string
	for _, result := range campaigns.Results {
		listCampaignID = append(listCampaignID, result.Metadata.ID)
	}
	records, _ := t.repos.Campaign.FindAllTrafficSourceID("outbrain", marketerID)
	if len(records) > 0 {
		for _, record := range records {
			if record.TrafficSourceID != "" && !utility.InArray(record.TrafficSourceID, listCampaignID, true) {
				campaigns.Results = append(campaigns.Results, dto.Result{Metadata: dto.Metadata{
					ID:   record.TrafficSourceID,
					Name: record.Name,
				}})
			}
		}
	}
	return
}

// Get toàn bộ Campaign đang chạy của report theo ngày đó - date format 2022-08-17 (yyyy-mm-dd)
func (t *outBrainUC) GetConversionIDs(marketerID string) (conversionIDs []string, err error) {
	type ConversionEvent struct {
		TrackingStatus       string `json:"trackingStatus"`
		Name                 string `json:"name"`
		Enabled              bool   `json:"enabled"`
		ConversionWindow     int    `json:"conversionWindow"`
		ID                   string `json:"id"`
		Type                 string `json:"type"`
		Category             string `json:"category"`
		IncludeInConversions bool   `json:"includeInConversions"`
		MarketerID           string `json:"marketerId"`
		URL                  string `json:"url,omitempty"`
		LastModified         string `json:"lastModified"`
		UsingNewTag          bool   `json:"usingNewTag,omitempty"`
		CreationTime         string `json:"creationTime"`
	}
	type Response struct {
		Count            int               `json:"count"`
		ConversionEvents []ConversionEvent `json:"conversionEvents"`
	}
	var resp Response
	url := "https://api.outbrain.com/amplify/v0.1/marketers/" + marketerID + "/conversionEvents"
	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	fmt.Println(url)
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("OB-TOKEN-V1", t.repos.ReportOutBrain.GetToken())
	})
	// Check Response
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &resp)
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
	err = c.Visit(url)
	c.Wait()
	// Sau khi đợi xử lý xong nếu có lỗi thì return luôn
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, conversionEvent := range resp.ConversionEvents {
		conversionIDs = append(conversionIDs, conversionEvent.ID)
	}
	return
}

// Get report realtime SectionByHour từ report của Outbrain
func (t *outBrainUC) GetSectionByHour(marketerID string, campaignID string, conversionIDs []string) (output []byte, err error) {
	time.Sleep(1200 * time.Millisecond)
	url := "https://api.outbrain.com/amplify/v0.1/realtime/marketers/" + marketerID + "/sectionsHourly"
	//=> Thêm Auth Key của để check report của mình
	values := url2.Values{}
	values.Set("campaignId", campaignID)
	values.Set("hours", "24")
	//fmt.Println(conversionIDs)
	values.Set("conversionIds", strings.Join(conversionIDs, ","))

	//=> Tạo link get Reports
	urlReport := url + "?" + values.Encode()
	fmt.Println(urlReport)

	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.SetRequestTimeout(2 * time.Second)
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("OB-TOKEN-V1", t.repos.ReportOutBrain.GetToken())
	})
	// Check Response
	c.OnResponse(func(r *colly.Response) {
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
	err = c.Visit(urlReport)
	c.Wait()
	return
}

// Get toàn bộ các marketerID đang chạy
func (t *outBrainUC) GetMarketerIDs() (marketerIDS []string, err error) {
	// => Tạo struct response từ API
	type Marketer struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Enabled      bool   `json:"enabled"`
		Currency     string `json:"currency"`
		CreationTime string `json:"creationTime"`
		LastModified string `json:"lastModified"`
	}
	type Response struct {
		Marketers []Marketer `json:"marketers"`
		Count     int        `json:"count"`
	}
	var res Response

	//=> Link api get toàn bộ marketer
	url := "https://api.outbrain.com/amplify/v0.1/marketers"

	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("OB-TOKEN-V1", t.repos.ReportOutBrain.GetToken())
	})
	// Check Response
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &res)
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
	err = c.Visit(url)
	c.Wait()

	// Sau khi đợi xử lý xong nếu có lỗi thì return luôn
	if err != nil {
		return
	}

	// => từ response nhận về lấy ra marketerIDS
	for _, marketer := range res.Marketers {
		marketerIDS = append(marketerIDS, marketer.ID)
	}
	return
}

// Update campaign
type InputUpdateCampaign struct {
	BlockedSites *InputUpdateCampaignBlockedSites `json:"blockedSites,omitempty"`
}

type InputUpdateCampaignBlockedSites struct {
	BlockedPublishers []InputUpdateBlocked `json:"blockedPublishers,omitempty"`
	BlockedSections   []InputUpdateBlocked `json:"blockedSections,omitempty"`
}

type InputUpdateBlocked struct {
	ID string `json:"id,omitempty"`
}

func (t *outBrainUC) UpdateCampaignOutBrain(campaignID string, inputs InputUpdateCampaign) (output []byte, err error) {
	url := "https://api.outbrain.com/amplify/v0.1/campaigns/" + campaignID
	// initialize http client
	client := &http.Client{}

	// marshal User to json
	bodyRequest, err := json.Marshal(inputs)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(bodyRequest))
	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyRequest))
	if err != nil {
		panic(err)
	}

	req.Header.Set("OB-TOKEN-V1", t.repos.ReportOutBrain.GetToken())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	output, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(output))
	return
}
