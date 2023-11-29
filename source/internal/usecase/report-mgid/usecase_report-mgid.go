package report_mgid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/google/go-querystring/query"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	"source/pkg/logger"
	"strconv"
	"time"
)

type UseCaseReportMgid interface {
	GetReport(campaignID string, date string) (records []*model.ReportMgidModel, err error)
	HandlerReport(records []*model.ReportMgidModel) (err error)
}

func NewReportMgidUC(repos *repo.Repositories) *mgidUC {
	return &mgidUC{
		repos: repos,
	}
}

const (
	URLMGIDGetToken                = "https://api.mgid.com/v1/auth/token"
	URLMGIDGetReportCampaignBySite = "https://api.mgid.com/v1/goodhits/"
)

type mgidUC struct {
	repos *repo.Repositories
}

func (t *mgidUC) getToken() (token string, err error) {
	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
		IdAuth       int    `json:"idAuth"`
	}

	var res response

	type request struct {
		Email    string `url:"email"`
		Password string `url:"password"`
		Force    string `url:"force"`
	}

	reqs := request{
		Email:    "ken@pubpower.io",
		Password: "Linh25101997@",
		Force:    "0",
	}

	params, _ := query.Values(reqs)
	c := &http.Client{}

	//fmt.Printf("%+v \n", params)
	req, _ := http.NewRequest("POST", URLMGIDGetToken, bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
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
	//if len(resp.InsertedID) == 0 {
	//	err = errors.New("submit fail")
	//}
	// Check the response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		return
	}

	token = "Bearer " + res.Token
	return
}

// Xử lý report Mgid
func (t *mgidUC) HandlerReport(records []*model.ReportMgidModel) (err error) {
	err = t.repos.ReportMgid.SaveSlice(records)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

// Get report
func (t *mgidUC) GetReport(campaignID string, date string) (records []*model.ReportMgidModel, err error) {
	// Lấy token
	token, err := t.getToken()
	if err != nil {
		return
	}

	url := URLMGIDGetReportCampaignBySite + "campaigns/" + campaignID + "/quality-analysis/uid"
	//=> Thêm Auth Key của để check report của mình
	values := url2.Values{}
	values.Set("dateInterval", "interval")
	values.Set("startDate", date)
	values.Set("endDate", date)

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
		request.Headers.Set("Authorization", token)
		request.Headers.Set("Cache-Control", "no-cache")
	})
	// Check Response
	var output []byte
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
	type report struct {
		Clicks        int64   `json:"clicks"`
		Spent         float64 `json:"spent"`
		Cpc           float64 `json:"cpc"`
		QualityFactor int64   `json:"qualityFactor"`
	}
	var response map[string]map[string]map[string]report
	err = json.Unmarshal(output, &response)
	if err != nil {
		return
	}
	for campaignID, mapInterval := range response {
		for _, mapReportBySite := range mapInterval {
			for sectionID, report := range mapReportBySite {
				records = append(records, &model.ReportMgidModel{
					Time:       date,
					CampaignID: campaignID,
					SectionID:  sectionID,
					Spent:      report.Spent,
					Clicks:     report.Clicks,
					CPC:        report.Cpc,
				})
			}
		}
	}
	return
}
