package report_tonic

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gofiber/fiber/v2"
	url2 "net/url"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	"strconv"
)

type UseCaseReportTonic interface {
	GetToken() (token string, expires int64, err error)
	GetReportFinal(token string, date string) (reports []*model.ReportTonicModel, err error)
	SaveReport(reports []*model.ReportTonicModel) (err error)
}

func NewReportTonicUC(repos *repo.Repositories) *tonicUC {
	return &tonicUC{
		userName: "42016694691571259338",
		password: "9761fde19defda1105b8ee2c9f2ecfe7d5386ca3",
		repos:    repos,
	}
}

type tonicUC struct {
	userName string
	password string
	repos    *repo.Repositories
}

func (t *tonicUC) GetToken() (token string, expires int64, err error) {
	url := "https://api.publisher.tonic.com/jwt/authenticate"
	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	type response struct {
		Token   string `json:"token"`
		Expires int64  `json:"expires"`
	}
	var res response
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Content-Type", "application/json")
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
	bodyRequest := fiber.Map{
		"consumer_key":    t.userName,
		"consumer_secret": t.password,
	}
	data, _ := json.Marshal(&bodyRequest)
	if err != nil {
		return "", 0, err
	}
	err = c.PostRaw(url, data)
	c.Wait()
	// Sau khi đợi xử lý xong nếu có lỗi thì return luôn
	if err != nil {
		return
	}
	token = res.Token
	expires = res.Expires
	return
}

func (t *tonicUC) GetReportFinal(token string, date string) (reports []*model.ReportTonicModel, err error) {
	url := "https://api.publisher.tonic.com/privileged/v3/epc/final"
	//=> Thêm Auth Key của để check report của mình
	values := url2.Values{}
	values.Set("from", date)
	values.Set("to", date)

	//=> Tạo link get Reports
	urlReport := url + "?" + values.Encode()
	fmt.Println(urlReport)

	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Authorization", "Bearer "+token)
		request.Headers.Set("Content-Type", "application/json")
	})
	// Check Response
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &reports)
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

	return
}

func (t *tonicUC) SaveReport(reports []*model.ReportTonicModel) (err error) {
	for _, report := range reports {
		query := make(map[string]interface{})
		query["date"] = report.Date
		query["campaign_id"] = report.CampaignId
		query["subid1"] = report.Subid1
		reportOld, _ := t.repos.ReportTonic.FindByQuery(query)
		if reportOld.ID != 0 {
			report.ID = reportOld.ID
			report.SectionId = reportOld.SectionId
			report.SectionName = reportOld.SectionName
			report.PublisherID = reportOld.PublisherID
			report.PublisherName = reportOld.PublisherName
			report.AdID = reportOld.AdID
			report.RedirectID = reportOld.RedirectID
		}

		if report.RedirectID == 0 {
			pixel := new(model.ReportAffPixelModel)
			queryPixel := make(map[string]interface{})
			if len(report.Subid1) > 99 {
				pixel, err = t.repos.ReportAffPixel.FindOneByQuery(queryPixel, "uid like '"+report.Subid1+"%'")
				if err != nil {
					continue
				}
			} else {
				queryPixel["uid"] = report.Subid1
				pixel, err = t.repos.ReportAffPixel.FindOneByQuery(queryPixel)
				if err != nil {
					continue
				}
			}
			report.SectionId = pixel.SectionId
			report.SectionName = pixel.SectionName
			report.PublisherID = pixel.PublisherID
			report.PublisherName = pixel.PublisherName
			report.AdID = pixel.AdID
			report.RedirectID, _ = strconv.ParseInt(pixel.RedirectID, 10, 64)
		}
		err = t.repos.ReportTonic.Save(report)
		if err != nil {
			continue
		}
	}
	return
}
