package report_pocpoc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"log"
	"net/http"
	"source/internal/entity/model"
	"source/internal/repo"
	"source/pkg/logger"
	"time"
)

type UseCaseReportPocPoc interface {
	GetReport(startDate string, endDate string) (records []*model.ReportPocPocModel, err error)
	HandlerReport(records []*model.ReportPocPocModel) (err error)
}

func NewReportPocPocUC(repos *repo.Repositories) *pocpocUC {
	return &pocpocUC{
		repos: repos,
	}
}

const (
	URLPocPocGetToken  = "https://api.pocpoc.io/account/login/admin"
	URLPocPocGetReport = "https://api.pocpoc.io/dashboard/overview"
)

type pocpocUC struct {
	repos      *repo.Repositories
	token      string
	expireTime time.Time
}

func (t *pocpocUC) getToken() (token string) {
	if t.expireTime.Sub(time.Now()).Milliseconds() < 0 || t.token == "" {
		_ = t.postGetToken()
	}
	return t.token
}

func (t *pocpocUC) postGetToken() (err error) {
	type response struct {
		Status bool `json:"status"`
		Data   struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	var res response

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	reqs := request{
		Email:    "khanhnh@bil.vn",
		Password: "12345@bil.vn",
	}
	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Content-Type", "application/json")
	})
	// Check Response
	c.OnResponse(func(r *colly.Response) {
		//fmt.Printf("%+v \n", string(r.Body))
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
	data, err := json.Marshal(&reqs)
	if err != nil {
		return
	}
	err = c.PostRaw(URLPocPocGetToken, data)
	c.Wait()
	// Sau khi đợi xử lý xong nếu có lỗi thì return luôn
	if err != nil {
		return
	}
	t.token = res.Data.AccessToken
	t.expireTime = time.Now().Add(9 * time.Minute)
	return
}

// Xử lý report PocPoc
func (t *pocpocUC) HandlerReport(records []*model.ReportPocPocModel) (err error) {
	size := 500
	var j int
	for i := 0; i < len(records); i += size {
		j += size
		if j > len(records) {
			j = len(records)
		}
		err = t.repos.ReportPocPoc.SaveSlice(records[i:j])
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
	return
}

// Get report
func (t *pocpocUC) GetReport(startDate string, endDate string) (records []*model.ReportPocPocModel, err error) {
	// Lấy token
	token := t.getToken()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//fmt.Println(token)
	type date struct {
		Label string `json:"label"`
	}
	type placement struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	type campaign struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	type dataPlacement struct {
		Date       date      `json:"date"`
		Placement  placement `json:"placement"`
		Campaign   campaign  `json:"campaign"`
		Revenue    float64   `json:"revenue"`
		Impression int64     `json:"impression"`
		Click      int64     `json:"click"`
		Spend      float64   `json:"spend"`
		Ecpm       float64   `json:"ecpm"`
		Ctr        float64   `json:"ctr"`
		Cpc        float64   `json:"cpc"`
		Viewable   float64   `json:"viewable"`
	}
	type data struct {
		Data []dataPlacement `json:"data"`
	}

	type response struct {
		Status bool `json:"status"`
		Data   data `json:"data"`
	}

	var res response

	type request struct {
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
		TimeInterval string `json:"time_interval"`
		Timezone     string `json:"timezone"`
		Currency     string `json:"currency"`
		Filters      struct {
			Inventory     []string `json:"inventory"`
			Placement     []string `json:"placement"`
			AdFormat      []string `json:"ad_format"`
			Device        []string `json:"device"`
			Country       []string `json:"country"`
			Region        []string `json:"region"`
			Publisher     []string `json:"publisher"`
			Advertiser    []string `json:"advertiser"`
			Campaign      []string `json:"campaign"`
			AdGroup       []string `json:"ad_group"`
			Ad            []string `json:"ad"`
			AdType        []string `json:"ad_type"`
			BidType       []string `json:"bid_type"`
			DataSource    []string `json:"data_source"`
			Conversion    []string `json:"conversion"`
			InventoryType []string `json:"inventory_type"`
		} `json:"filters"`
		Groups     []string    `json:"groups"`
		Metrics    []string    `json:"metrics"`
		Limit      int         `json:"limit"`
		Sort       string      `json:"sort"`
		Comparison interface{} `json:"comparison"`
	}
	c := &http.Client{}

	reqs := request{
		StartDate:    startDate,
		EndDate:      endDate,
		TimeInterval: "day",
		Timezone:     "UTC",
		Currency:     "USD",
		Filters: struct {
			Inventory     []string `json:"inventory"`
			Placement     []string `json:"placement"`
			AdFormat      []string `json:"ad_format"`
			Device        []string `json:"device"`
			Country       []string `json:"country"`
			Region        []string `json:"region"`
			Publisher     []string `json:"publisher"`
			Advertiser    []string `json:"advertiser"`
			Campaign      []string `json:"campaign"`
			AdGroup       []string `json:"ad_group"`
			Ad            []string `json:"ad"`
			AdType        []string `json:"ad_type"`
			BidType       []string `json:"bid_type"`
			DataSource    []string `json:"data_source"`
			Conversion    []string `json:"conversion"`
			InventoryType []string `json:"inventory_type"`
		}{
			Inventory:     []string{},
			Placement:     []string{},
			AdFormat:      []string{},
			Device:        []string{},
			Country:       []string{},
			Region:        []string{},
			Publisher:     []string{},
			Advertiser:    []string{"15"},
			Campaign:      []string{},
			AdGroup:       []string{},
			Ad:            []string{},
			AdType:        []string{},
			BidType:       []string{},
			DataSource:    []string{},
			Conversion:    []string{},
			InventoryType: []string{},
		},
		Groups: []string{
			"date",
			"placement",
			"campaign",
		},
		Metrics: []string{
			"revenue",
			"impression",
			"click",
			"spend",
			"ecpm",
			"ctr",
			"cpc",
			"viewable",
			"rpm",
			"fillrate",
		},
		Limit:      0,
		Sort:       "revenue_desc",
		Comparison: nil,
	}
	jsonStr, err := json.Marshal(reqs)
	if err != nil {
		return
	}
	fmt.Println(URLPocPocGetReport)
	req, _ := http.NewRequest("POST", URLPocPocGetReport, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", token)
	// Submit the request
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%+v \n", string(bodyBytes))

	_ = json.Unmarshal(bodyBytes, &res)
	// Check the response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		logger.Error(err.Error())
		fmt.Printf("%+v \n", string(bodyBytes))
		return
	}

	if !res.Status {
		err = fmt.Errorf("bad status: %s", resp.Status)
		logger.Error(err.Error())
		fmt.Printf("%+v \n", string(bodyBytes))
		return
	}

	for _, data := range res.Data.Data {
		records = append(records, &model.ReportPocPocModel{
			Time:         data.Date.Label,
			CampaignID:   data.Campaign.Id,
			CampaignName: data.Campaign.Name,
			SectionID:    data.Placement.Id,
			SectionName:  data.Placement.Name,
			Clicks:       data.Click,
			Spent:        data.Spend,
			CPC:          data.Cpc,
		})
	}
	return
}
