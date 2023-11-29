package report_taboola

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	"strconv"
	"time"
)

type UseCaseReportTaboola interface {
	GetToken(clientID, clientSecret string) (token string, expires int64, err error)
	GetAllAccountID(token string) (accountIDs []string, err error)
	GetReportCampaignSummary(accountTaboola string, token string, startDate, endDate string, accountID string) (reports []*model.ReportTaboolaModel, err error)
	SaveReport(reports []*model.ReportTaboolaModel) (err error)
}

func NewReportTaboolaUC(repos *repo.Repositories) *taboolaUC {
	return &taboolaUC{
		repos: repos,
	}
}

type taboolaUC struct {
	repos *repo.Repositories
}

func (t *taboolaUC) GetToken(clientID, clientSecret string) (token string, expires int64, err error) {
	url := "https://backstage.taboola.com/backstage/oauth/token"
	fmt.Println(url)

	type response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	var res response

	type request struct {
		ClientID     string `url:"client_id"`
		ClientSecret string `url:"client_secret"`
		GrantType    string `url:"grant_type"`
	}

	reqs := request{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	params, _ := query.Values(reqs)
	c := &http.Client{}

	fmt.Printf("%+v \n", params)
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(params.Encode()))
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

	fmt.Printf("%+v \n", string(bodyBytes))

	_ = json.Unmarshal(bodyBytes, &res)
	//if len(resp.InsertedID) == 0 {
	//	err = errors.New("submit fail")
	//}
	// Check the response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
	}

	token = res.AccessToken
	expires = time.Now().Unix() + res.ExpiresIn
	return
}

func (t *taboolaUC) GetAllAccountID(token string) (accountIDs []string, err error) {
	url := "https://backstage.taboola.com/backstage/api/1.0/users/current/allowed-accounts"
	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	type response struct {
		Results []struct {
			Id              int      `json:"id"`
			Name            string   `json:"name"`
			AccountId       string   `json:"account_id"`
			PartnerTypes    []string `json:"partner_types"`
			Type            string   `json:"type"`
			CampaignTypes   []string `json:"campaign_types"`
			Currency        string   `json:"currency"`
			TimeZoneName    string   `json:"time_zone_name"`
			DefaultPlatform string   `json:"default_platform"`
			IsActive        bool     `json:"is_active"`
			Language        string   `json:"language"`
			Country         string   `json:"country"`
		} `json:"results"`
		Metadata struct {
			Total             int           `json:"total"`
			Count             int           `json:"count"`
			StaticFields      []interface{} `json:"static_fields"`
			StaticTotalFields []interface{} `json:"static_total_fields"`
			DynamicFields     interface{}   `json:"dynamic_fields"`
			StartDate         interface{}   `json:"start_date"`
			EndDate           interface{}   `json:"end_date"`
		} `json:"metadata"`
	}

	var res response
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Authorization", "Bearer "+token)
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
	if err != nil {
		return
	}
	err = c.Visit(url)
	c.Wait()
	// Sau khi đợi xử lý xong nếu có lỗi thì return luôn
	if err != nil {
		return
	}
	for _, result := range res.Results {
		accountIDs = append(accountIDs, result.AccountId)
	}
	return
}

func (t *taboolaUC) GetReportCampaignSummary(accountTaboola string, token string, startDate, endDate string, accountID string) (reports []*model.ReportTaboolaModel, err error) {
	url := "https://backstage.taboola.com/backstage/api/1.0/" + accountID + "/reports/campaign-summary/dimensions/campaign_site_day_breakdown"
	//=> Thêm Auth Key của để check report của mình
	values := url2.Values{}
	values.Set("start_date", startDate)
	values.Set("end_date", endDate)
	type response struct {
		LastUsedRawdataUpdateTime            string `json:"last-used-rawdata-update-time"`
		LastUsedRawdataUpdateTimeGmtMillisec int64  `json:"last-used-rawdata-update-time-gmt-millisec"`
		Timezone                             string `json:"timezone"`
		Results                              []struct {
			Date               string  `json:"date"`
			Site               string  `json:"site"`
			SiteName           string  `json:"site_name"`
			SiteId             int64   `json:"site_id"`
			Campaign           string  `json:"campaign"`
			CampaignName       string  `json:"campaign_name"`
			Clicks             int64   `json:"clicks"`
			Impressions        int64   `json:"impressions"`
			VisibleImpressions int64   `json:"visible_impressions"`
			Spent              float64 `json:"spent"`
			ConversionsValue   float64 `json:"conversions_value"`
			Roas               float64 `json:"roas"`
			Ctr                float64 `json:"ctr"`
			Vctr               float64 `json:"vctr"`
			Cpm                float64 `json:"cpm"`
			Vcpm               float64 `json:"vcpm"`
			Cpc                float64 `json:"cpc"`
			Cpa                float64 `json:"cpa"`
			CpaActionsNum      int64   `json:"cpa_actions_num"`
			CpaConversionRate  float64 `json:"cpa_conversion_rate"`
			BlockingLevel      string  `json:"blocking_level"`
			Currency           string  `json:"currency"`
		} `json:"results"`
		RecordCount int `json:"recordCount"`
		Metadata    struct {
			Total        int64  `json:"total"`
			Count        int64  `json:"count"`
			StartDate    string `json:"start_date"`
			EndDate      string `json:"end_date"`
			StaticFields []struct {
				Id       string  `json:"id"`
				Format   *string `json:"format"`
				DataType string  `json:"data_type"`
			} `json:"static_fields"`
		} `json:"metadata"`
	}
	//=> Tạo link get Reports
	urlReport := url + "?" + values.Encode()
	fmt.Println(urlReport)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, urlReport, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var resp response
	err = json.Unmarshal(output, &resp)
	if err != nil {
		return
	}
	for _, result := range resp.Results {
		report := model.ReportTaboolaModel{
			Account:            accountTaboola,
			Date:               result.Date,
			CampaignID:         result.Campaign,
			CampaignName:       result.CampaignName,
			SiteID:             strconv.FormatInt(result.SiteId, 10),
			Site:               result.Site,
			SiteName:           result.SiteName,
			Clicks:             result.Clicks,
			Impressions:        result.Impressions,
			VisibleImpressions: result.VisibleImpressions,
			Spent:              result.Spent,
			ConversionsValue:   result.ConversionsValue,
			Roas:               result.Roas,
			Ctr:                result.Ctr,
			Vctr:               result.Vctr,
			Cpm:                result.Cpm,
			Vcpm:               result.Vcpm,
			Cpc:                result.Cpc,
			Cpa:                result.Cpa,
			CpaActionsNum:      result.CpaActionsNum,
			CpaConversionRate:  result.CpaConversionRate,
			BlockingLevel:      result.BlockingLevel,
			Currency:           result.Currency,
		}
		reports = append(reports, &report)
	}
	return
}

func (t *taboolaUC) SaveReport(reports []*model.ReportTaboolaModel) (err error) {
	return t.repos.ReportTaboola.SaveSlice(reports)
}
