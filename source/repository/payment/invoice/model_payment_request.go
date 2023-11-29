package invoice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/syyongx/php2go"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"source/apps/backend/input"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/telegram"
	"source/pkg/utility"
	"strconv"
	"strings"
	"sync"
	"time"
	// "time"
)

const (
	TIME_ZONE_DEFAULT    = "America/New_York"
	URlQuantumdex        = "https://be.quantumdex.io/"
	URLReValueimpression = "https://re.valueimpression.com/"
	URLAppValue          = "https://app.valueimpression.com/"
)

type PaymentRequest struct{}

type PaymentRequestRecord struct {
	mysql.TablePaymentRequest
}

func (PaymentRequestRecord) TableName() string {
	return mysql.Tables.PaymentRequest
}

func (t *PaymentRequest) ScanPaymentPublisher() (errs []ajax.Error) {
	errs = t.ScanRequestPublisher()
	if len(errs) > 0 {
		telegram.SendErrorsTuan(errs, "Scan Request Publisher")
		return
	}
	errs = new(PaymentInvoice).ScanInvoicePublisher()
	if len(errs) > 0 {
		telegram.SendErrorsTuan(errs, "Scan Invoice Publisher")
		return
	}
	err := new(PaymentInvoice).BuildInvoicePDF()
	if err != nil {
		telegram.SendErrorTuan(err.Error(), "Build Invoice PDF")
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	return
}

func (t *PaymentRequest) ScanPaymentAgency() (errs []ajax.Error) {

	errs = t.ScanRequestAgency()
	if len(errs) > 0 {
		telegram.SendErrorsTuan(errs, "Scan Request Agency")
		return
	}

	errs = new(PaymentInvoice).ScanInvoiceAgency()
	if len(errs) > 0 {
		telegram.SendErrorsTuan(errs, "Scan Invoice Agency")
		return
	}
	return
}

func (t *PaymentRequest) ScanRequestAgency() (errs []ajax.Error) {
	agencies, err := getAllAgency()
	// agencies, err := getUsersById(146)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "There are no agencies!",
		})
		return
	}
	for _, agency := range agencies {
		// fmt.Printf("%+v\n", agency.Email)

		// lấy khoảng thời gian tạo request
		periodRequest := PeriodRequest{}
		loc, _ := time.LoadLocation("America/New_York")
		now := time.Now().In(loc)
		year, month, _ := now.Date()
		firstDayOfLastMonth := time.Date(year, month-1, 1, 0, 0, 0, 0, now.Location())
		endDayOfLastMonth := time.Date(year, month, 0, 0, 0, 0, 0, now.Location())

		periodRequest.StartDate = firstDayOfLastMonth.Format("2006-01-02")
		periodRequest.EndDate = endDayOfLastMonth.Format("2006-01-02")
		// periodRequest.StartDate = "2022-08-01"
		// periodRequest.EndDate = "2022-08-31"

		// check exist request in periodRequest
		exist := t.CheckExistRequestAgencyInMonth(agency.Id, endDayOfLastMonth.Format("2006-01"))
		if exist.Id != 0 {
			continue
		}
		// get report for user in time
		report, errss := t.ReportAgencyInTime(agency, periodRequest.StartDate, periodRequest.EndDate)
		// fmt.Printf("%+v\n", report)
		// fmt.Println("--------------------------------------------------------------------------------")

		if len(errss) > 0 {
			errs = append(errs, errss...)
			continue
		}
		if report.Commission < 0.01 {
			continue
		}

		var paymentRequest PaymentRequestRecord
		paymentRequest.Revenue = report.GrossRevenue - report.RevenuePublisher - report.RevenueAffiliate
		paymentRequest.Amount = report.Commission
		paymentRequest.Rate = report.Rate * 100
		paymentRequest.Type = mysql.TypeAgencyProfit
		paymentRequest.Status = mysql.StatusPaymentPending
		paymentRequest.UserId = agency.Id
		paymentRequest.StartDate, _ = time.ParseInLocation("2006-01-02", periodRequest.StartDate, loc)
		paymentRequest.EndDate, _ = time.ParseInLocation("2006-01-02", periodRequest.EndDate, loc)
		paymentRequest.Note = "$" + php2go.NumberFormat(paymentRequest.Revenue, 2, ".", ",") + " (Rate: " + php2go.NumberFormat(report.Rate*100, 0, ".", ",") + "%)"
		err = mysql.Client.Create(&paymentRequest).Error
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      agency.Email,
				Message: err.Error(),
			})
		}

		if agency.Email != "davidwright@valueimpression.com" {
			// share manager receive agency
			var RequestShare PaymentRequestRecord
			RequestShare = paymentRequest
			RequestShare.Id = 0
			RequestShare.Amount = (report.GrossRevenue - report.RevenuePublisher - report.RevenueAffiliate) * (0.15 - report.Rate)
			RequestShare.Type = mysql.TypeAgencyShare
			RequestShare.UserId = 137 // id của davidwright@valueimpression.com
			RequestShare.Note = "$" + php2go.NumberFormat(RequestShare.Revenue, 2, ".", ",") + " (Rate: " + php2go.NumberFormat(15-report.Rate*100, 0, ".", ",") + "%) " + agency.Email
			err = mysql.Client.Create(&RequestShare).Error
		}
	}
	return
}

type ReportAgency struct {
	GrossRevenue     float64 `json:"grossRevenue" form:"grossRevenue"`
	NetRevenue       float64 `json:"netRevenue" form:"netRevenue"`
	RevenuePublisher float64 `json:"revenuePublisher" form:"revenuePublisher"`
	RevenueAffiliate float64 `json:"revenueAffiliate" form:"revenueAffiliate"`
	Rate             float64 `json:"rate" form:"rate"`
	Commission       float64 `json:"commission" form:"commission"`
}

func (t *PaymentRequest) ReportAgencyInTime(agency UserRecord, StartDate, EndDate string) (report ReportAgency, errs []ajax.Error) {
	query := "SELECT SUM(grossRevenue) as grossRevenue, SUM(netRevenue) as netRevenue FROM \"pw_bidder\" WHERE owner = 'apac' AND accManageId = " + strconv.FormatInt(agency.Id, 10) + " and FLOOR(__time TO DAY) >= '" + StartDate + "' AND FLOOR(__time TO DAY) <= '" + EndDate + "'"
	body, err := druid(query)
	if err != nil {
		errs = append(errs, ajax.Error{
			Message: err.Error(),
		})
	}
	var reports []ReportAgency
	err = json.Unmarshal([]byte(body), &reports)
	if err != nil {
		errs = append(errs, ajax.Error{
			Message: err.Error(),
		})
	}

	if len(reports) == 0 {
		return
	}
	report = reports[0]

	publishers, err := PublishersForAgency(agency.Id)
	// publishers, err := getUsersById(116)
	if err != nil {
		errs = append(errs, ajax.Error{
			Message: err.Error(),
		})
		return
	}

	for _, publisher := range publishers {
		var reportsPublisher []input.ReportPublisher
		reportsPublisher, err = t.ReportPublisherInTime(publisher.Id, StartDate, EndDate)
		if err != nil {
			errs = append(errs, ajax.Error{
				Message: err.Error(),
			})
		}
		if len(reportsPublisher) == 0 {
			continue
		}

		for _, reportPublisher := range reportsPublisher {
			if reportPublisher.Revenue == "" || reportPublisher.Revenue == "0" || reportPublisher.Revenue == "0.0" || reportPublisher.Rate == "0" || reportPublisher.Rate == "" {
				continue
			}
			var revenue float64
			revenue, err = strconv.ParseFloat(reportPublisher.Revenue, 64)
			if err != nil {
				errs = append(errs, ajax.Error{
					Message: err.Error(),
				})
			}
			var rate float64
			rate, err = strconv.ParseFloat(reportPublisher.Rate, 64)
			if err != nil {
				errs = append(errs, ajax.Error{
					Message: err.Error(),
				})
			}
			report.RevenuePublisher = report.RevenuePublisher + (revenue * rate / 100)
		}
	}

	date, _ := time.Parse("2006-01-02", EndDate)
	report.Rate, err = t.RateSharingByAgency(agency, date.Format("200601"))
	if err != nil {
		if err != nil {
			errs = append(errs, ajax.Error{
				Message: err.Error(),
			})
		}
	}
	report.Commission = (report.GrossRevenue - report.RevenuePublisher - report.RevenueAffiliate) * report.Rate
	return
}

type ResponseRate struct {
	Status  bool    `json:"status"`
	Message string  `json:"message"`
	Rate    float64 `json:"rate"`
}

func (t *PaymentRequest) RateSharingByAgency(agency UserRecord, month string) (rate float64, err error) {

	url := URLReValueimpression + "api/payment/GetRateAgency"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("email", agency.Email)
	_ = writer.WriteField("month", month)
	err = writer.Close()
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return
	}
	req.Header.Add("Password", "bli@123")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var responseRate ResponseRate
	err = json.Unmarshal(body, &responseRate)
	if err != nil {
		return
	}
	if !responseRate.Status {
		err = errors.New(responseRate.Message)
		return
	}
	rate = responseRate.Rate
	return

	//if agency.Email == "davidwright@valueimpression.com" {
	//	rate = 0.15 // 15%
	//	return
	//}
	//if agency.Email == "l.watson@valueimpression.com" {
	//	rate = 0.1 // 10%
	//	return
	//}
	//type Result struct {
	//	Revenue float64 `json:"revenue"`
	//}
	//// get revenue bên App (App.valueimporession.com)
	//responseApp, err := t.APIGetRevenueOfAgencyInMonthFormApp(agency.Email, month)
	//if err != nil {
	//	return
	//}
	//var reportApp Result
	//json.Unmarshal(responseApp, &reportApp)
	//// get revenue bên Exchange (be.quantumdex.io)
	//responseExchange, err := t.APIGetRevenueOfAgencyInMonthFormExchange(agency.Email, month)
	//if err != nil {
	//	return
	//}
	//var reportExchange Result
	//json.Unmarshal(responseExchange, &reportExchange)
	//// get revenue bên VLI (re.valueimpression.com)
	//responseVLI, err := t.APIGetRevenueOfAgencyInMonthFormVLI(agency.Email, month)
	//if err != nil {
	//	return
	//}
	//var reportVLI Result
	//json.Unmarshal(responseVLI, &reportVLI)
	//
	//totalRevenue := PubpowerRevenue + reportApp.Revenue + reportExchange.Revenue + reportVLI.Revenue
	//if totalRevenue <= 50000 {
	//	rate = float64(5) / float64(100)
	//}
	//if totalRevenue > 50000 && totalRevenue <= 100000 {
	//	rate = float64(6) / float64(100)
	//}
	//if totalRevenue > 100000 && totalRevenue <= 200000 {
	//	rate = float64(7) / float64(100)
	//}
	//if totalRevenue > 200000 && totalRevenue <= 500000 {
	//	rate = float64(8) / float64(100)
	//}
	//if totalRevenue > 500000 && totalRevenue <= 1000000 {
	//	rate = float64(9) / float64(100)
	//}
	//if totalRevenue > 1000000 {
	//	rate = float64(10) / float64(100)
	//}
	return
}

func (t *PaymentRequest) APIGetRevenueOfAgencyInMonthFormExchange(mail, month string) (body []byte, err error) {
	url := URlQuantumdex + "api/report/APIRevenueOfAgencyInMonth"
	payload := strings.NewReader("email=" + mail + "&month=" + month + "&token=8793b0c946610af48daefb07740c9a9a")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	return
}

func (t *PaymentRequest) APIGetRevenueOfAgencyInMonthFormApp(mail, month string) (body []byte, err error) {
	url := URLAppValue + "api/report/APIRevenueOfAgencyInMonth"
	payload := strings.NewReader("email=" + mail + "&month=" + month + "&token=8793b0c946610af48daefb07740c9a9a")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	return
}

func (t *PaymentRequest) APIGetRevenueOfAgencyInMonthFormVLI(mail, month string) (body []byte, err error) {
	url := URLReValueimpression + "api/report/APIRevenueOfAgencyInMonth"
	payload := strings.NewReader("email=" + mail + "&month=" + month + "&token=8793b0c946610af48daefb07740c9a9a")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	return
}

func (t *PaymentRequest) CheckExistRequestAgencyInMonth(agencyId int64, month string) (request PaymentRequestRecord) {
	mysql.Client.Where("user_id = ? and end_date like ? and type = 3", agencyId, month+"%").Find(&request)
	return
}

func (t *PaymentRequest) ScanRequestPublisher() (errs []ajax.Error) {
	publishers, err := getAllPublisher()
	// publishers, err := getUsersById(8)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "There are no publishers!",
		})
		return
	}

	wg := &sync.WaitGroup{}
	concurrency := 10
	limitGoRoutine := make(chan bool, concurrency)
	Errors := make(chan string, len(publishers))

	for _, publisher := range publishers {
		limitGoRoutine <- true
		wg.Add(1)

		go func(publisher UserRecord, limitGoRoutine chan bool, wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
				<-limitGoRoutine
			}()
			loc, _ := time.LoadLocation("America/New_York")
			// lấy khoảng thời gian tạo request
			periodRequest := t.getPeriodPayment(publisher)
			if periodRequest.EndDate != "" && periodRequest.StartDate != "" {
				// get report for user in time
				reports, err := t.ReportPublisherInTime(publisher.Id, periodRequest.StartDate, periodRequest.EndDate)
				if err != nil {
					Errors <- err.Error()
					return
				}
				var revenuePublisher float64
				var revenueTotal float64
				for _, report := range reports {
					revenue, err := strconv.ParseFloat(report.Revenue, 64)
					if err != nil {
						Errors <- err.Error()
						return
					}
					rate, err := strconv.ParseFloat(report.Rate, 64)
					if err != nil {
						Errors <- err.Error()
						return
					}
					revenueTotal = revenueTotal + revenue
					revenuePublisher = revenuePublisher + (revenue*rate)/100
				}

				flagGuarantee := false
				GuaranteFrom, _ := strconv.Atoi(php2go.StrReplace("-", "", publisher.GuaranteeFrom, 10)) // thời gian bắt đầu áp dụng guarantee
				StartDate, _ := strconv.Atoi(php2go.StrReplace("-", "", periodRequest.StartDate, 10))
				EndDate, _ := strconv.Atoi(php2go.StrReplace("-", "", periodRequest.EndDate, 10))
				// xem user có set guarantee(invoice sàn) nếu có check xem nếu total revenue thấp hơn guarantee thì cho total revenue = guarantee
				if publisher.PaymentTerm != 4 && publisher.Guarantee != 0 && (revenuePublisher < float64(publisher.Guarantee) || publisher.GuaranteeCeiling == "on") {
					if publisher.GuaranteePeriods == "once" && GuaranteFrom >= StartDate && GuaranteFrom <= EndDate { // nếu guarantee chỉ áp dụng 1 lần
						flagGuarantee = true
					} else if publisher.GuaranteePeriods == "" && GuaranteFrom <= EndDate { // nếu guarantee áp dụng vô hạn
						flagGuarantee = true
					}
				}
				if revenuePublisher <= 0 {
					return
				}

				if flagGuarantee {
					// create request for publisher
					StartDateRequest, _ := time.ParseInLocation("2006-01-02", periodRequest.StartDate, loc)
					EndDateRequest, _ := time.ParseInLocation("2006-01-02", periodRequest.EndDate, loc)
					err := mysql.Client.Create(&PaymentRequestRecord{mysql.TablePaymentRequest{
						UserId:    publisher.Id,
						Revenue:   float64(publisher.Guarantee),
						Amount:    float64(publisher.Guarantee),
						Type:      1,
						Note:      "Guarantee $" + php2go.NumberFormat(float64(publisher.Guarantee), 2, ".", ","),
						Status:    1,
						StartDate: StartDateRequest,
						EndDate:   EndDateRequest,
					}}).Error
					if err != nil {
						Errors <- err.Error()
						return
					}

					// nếu guarantee chỉ áp dụng 1 lần và payment request đã được tạo thì xóa guarantee đi để k chạy guarantee lần nào nữa
					if publisher.GuaranteePeriods == "once" && GuaranteFrom >= StartDate && GuaranteFrom <= EndDate { // nếu guarantee chỉ áp dụng 1 lần
						RemoveGuaranteeForUser(publisher.Id)
					}
					// create request for agency (agency sẽ chịu phần publisher nhận được khi $user->guarantee > $total_revenue)
					if revenuePublisher < float64(publisher.Guarantee) {
						var Request PaymentRequestRecord
						Request.Amount = float64(publisher.Guarantee) - revenuePublisher
						Request.Note = "Guarantee $" + php2go.NumberFormat(float64(publisher.Guarantee), 2, ".", ",") + " - " + "Commission $" + php2go.NumberFormat(revenuePublisher, 2, ".", ",") + " Pub:" + publisher.Email
						startDate, _ := time.ParseInLocation("2006-01-02", periodRequest.StartDate, loc)
						Request.StartDate = startDate
						endDate, _ := time.ParseInLocation("2006-01-02", periodRequest.EndDate, loc)
						Request.EndDate = endDate
						t.CreateRequestGuaranteeForAgency(publisher.Presenter, Request)
					}
				} else {
					// create request
					for _, report := range reports {
						revenue, _ := strconv.ParseFloat(report.Revenue, 64)
						rate, _ := strconv.ParseFloat(report.Rate, 64)
						StartDateRequest, _ := time.ParseInLocation("2006-01-02", report.StartDate, loc)
						EndDateRequest, _ := time.ParseInLocation("2006-01-02", report.EndDate, loc)
						err := mysql.Client.Create(&PaymentRequestRecord{mysql.TablePaymentRequest{
							UserId:    publisher.Id,
							Revenue:   revenue,
							Amount:    revenue * rate / 100,
							Rate:      rate,
							Type:      1,
							Status:    1,
							StartDate: StartDateRequest,
							EndDate:   EndDateRequest,
						}}).Error

						if err != nil {
							Errors <- err.Error()
							return
						}
					}
				}
			}
		}(publisher, limitGoRoutine, wg)
	}

	close(Errors)
	wg.Wait()
	for Error := range Errors {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: Error,
		})
	}
	return
}

func (t *PaymentRequest) CreateRequestGuaranteeForAgency(agencyId int64, row PaymentRequestRecord) (err error) {
	if row.Amount < 0.01 || agencyId == 0 {
		return
	}
	row.Creator = 0
	row.UserId = agencyId
	row.Type = mysql.TypeAgencyGuarantee
	row.Status = mysql.StatusPaymentPending
	err = mysql.Client.Create(&row).Error
	return
}

type PeriodRequest struct {
	StartDate string
	EndDate   string
}

func (t *PaymentRequest) getPeriodPayment(publisher UserRecord) (PeriodRequest PeriodRequest) {
	var input input.GetStartDateUser
	input.UserId = publisher.Id
	startDate, err := t.GetStartDateRequest(input)
	if err != nil || startDate == "" {
		return
	}

	// set time zone America/New_York
	loc, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(loc)
	startTimeRequest, _ := time.ParseInLocation("2006-01-02", startDate, loc)
	endTimeRequest := time.Time{}.In(loc)

	switch publisher.PaymentTerm {
	case 1: // month
		year, month, _ := now.Date()
		endTimeRequest = time.Date(year, month, 0, 0, 0, 0, 0, now.Location())
		break
	case 2: // week
		endTimeRequest = startTimeRequest.AddDate(0, 0, +6)
		break
	case 3: // bi-weekly
		endTimeRequest = startTimeRequest.AddDate(0, 0, +13)
		break
	default: // other k tạo request => tạo bằng tay ở backend
		return
	}
	if endTimeRequest.Unix() < startTimeRequest.Unix() || endTimeRequest.Unix() > time.Now().Unix() {
		return
	}

	PeriodRequest.StartDate = startTimeRequest.Format("2006-01-02")
	PeriodRequest.EndDate = endTimeRequest.Format("2006-01-02")
	return
}

func (t *PaymentRequest) Create(publisher UserRecord, report input.ReportPublisher) (row PaymentRequestRecord, err error) {
	row = t.makeRow(publisher, report)
	err = mysql.Client.Create(&row).Error

	return
}

func (t *PaymentRequest) makeRow(publisher UserRecord, record input.ReportPublisher) (row PaymentRequestRecord) {
	loc, _ := time.LoadLocation("America/New_York")
	row.Creator = 0
	row.UserId = publisher.Id
	revenue, _ := strconv.ParseFloat(record.Revenue, 64)
	row.Amount = revenue
	row.Currency = record.Currency
	row.Note = record.Note
	row.Type = mysql.TypeCommission
	row.Status = mysql.StatusPaymentPending
	startDate, _ := time.ParseInLocation("2006-01-02", record.StartDate, loc)
	row.StartDate = startDate.In(loc)
	endDate, _ := time.ParseInLocation("2006-01-02", record.EndDate, loc)
	row.EndDate = endDate.In(loc)
	return
}

func (t *PaymentRequest) ReportPublisherInTime(userId int64, StartDate, EndDate string) (reports []input.ReportPublisher, err error) {
	// loc, _ := time.LoadLocation("America/New_York")
	publisher, err := GetUserById(userId)
	if err != nil || publisher.Email == "" {
		return
	}
	// nếu Publisher đc sync từ VLI
	if publisher.SystemSync == 1 {
		report, __err := t.ApiGetReportPublisher(userId, StartDate, EndDate)
		if __err != nil {
			return reports, __err
		}
		report.Rate = "100"
		reports = append(reports, report)
		return
	}

	// get revenue share for user
	revenueShares, err := GetRevenueSharesForUser(userId, StartDate, EndDate)
	startDateReport := "" // thời gian bắt đầu get report
	endDateReport := ""   // thời gian cuối cùng get report
	for index, RevenueShare := range revenueShares {
		if len(revenueShares) == 1 && RevenueShare.Date.Format("2006-01-02") <= StartDate {
			// nếu revenue share đc set thời gian trước hoặc = StartDate => 1 payment_request từ StartDate -> EndDate
			report, __err := t.ApiGetReportPublisher(userId, StartDate, EndDate)
			if __err != nil {
				return reports, __err
			}
			report.Rate = strconv.FormatInt(RevenueShare.Rate, 10)
			reports = append(reports, report)
		} else {
			var report input.ReportPublisher

			DateRevenueShare := RevenueShare.Date.Format("2006-01-02")
			LastDateRevenueShare := RevenueShare.Date.AddDate(0, 0, -1).Format("2006-01-02")

			if DateRevenueShare <= StartDate {
				if len(revenueShares)-1 == index {
					report, err = t.ApiGetReportPublisher(userId, StartDate, EndDate)
					if err != nil {
						return
					}
					if report.Revenue != "" && report.Revenue != "0" && report.Revenue != "0.0" {
						report.Rate = strconv.FormatInt(RevenueShare.Rate, 10)
					}
					reports = append(reports, report)
				}
				continue
			} else if DateRevenueShare > StartDate && DateRevenueShare < EndDate {
				if startDateReport == "" {
					startDateReport = StartDate
					endDateReport = LastDateRevenueShare
				} else {
					_startDateReport, _ := time.Parse("2006-01-02", endDateReport)
					startDateReport = _startDateReport.AddDate(0, 0, 1).Format("2006-01-02")
					endDateReport = LastDateRevenueShare
				}
				report, err = t.ApiGetReportPublisher(userId, startDateReport, endDateReport)
				if err != nil {
					return
				}
				if report.Revenue != "" && report.Revenue != "0" && report.Revenue != "0.0" {
					RevenueShareOld, errrr := GetRevenueShareOldForUser(userId, RevenueShare.Date.Format("2006-01-02"))
					if errrr != nil {
						err = errrr
						return
					}
					report.Rate = strconv.FormatInt(RevenueShareOld.Rate, 10)
				}
				reports = append(reports, report)

				// nếu là revenue share mới nhất ==> get report từ DateRevenueShare -> EndDate
				if len(revenueShares)-1 == index {
					_startDateReport, _ := time.Parse("2006-01-02", endDateReport)
					startDateReport = _startDateReport.AddDate(0, 0, 1).Format("2006-01-02")
					endDateReport = EndDate

					report, err = t.ApiGetReportPublisher(userId, startDateReport, endDateReport)
					if err != nil {
						return
					}

					if report.Revenue != "" && report.Revenue != "0" && report.Revenue != "0.0" {
						report.Rate = strconv.FormatInt(RevenueShare.Rate, 10)
					}
					reports = append(reports, report)
				}
			} else if DateRevenueShare == EndDate {
				if endDateReport != "" {
					_startDateReport, _ := time.Parse("2006-01-02", endDateReport)
					startDateReport = _startDateReport.AddDate(0, 0, 1).Format("2006-01-02")
					endDateReport = LastDateRevenueShare
				} else {
					startDateReport = StartDate
					endDateReport = LastDateRevenueShare
				}
				// nếu revenue share đc set thời gian trước hoặc = StartDate => 1 payment_request từ StartDate -> EndDate
				report, err = t.ApiGetReportPublisher(userId, startDateReport, endDateReport)
				if err != nil {
					return
				}
				if report.Revenue != "" && report.Revenue != "0" && report.Revenue != "0.0" {
					RevenueShareOld, errrr := GetRevenueShareOldForUser(userId, RevenueShare.Date.Format("2006-01-02"))
					if errrr != nil {
						err = errrr
						return
					}
					report.Rate = strconv.FormatInt(RevenueShareOld.Rate, 10)
				}
				reports = append(reports, report)

				// nếu là revenue share mới nhất ==> get report từ DateRevenueShare -> EndDate
				if len(revenueShares)-1 == index {
					startDateReport = EndDate
					endDateReport = EndDate
					report, err = t.ApiGetReportPublisher(userId, startDateReport, endDateReport)
					if err != nil {
						return
					}
					if report.Revenue != "" && report.Revenue != "0" && report.Revenue != "0.0" {
						report.Rate = strconv.FormatInt(RevenueShare.Rate, 10)
					}
					reports = append(reports, report)
				}
			}
		}
	}
	return
}

// func (t *PaymentRequest) ApiGetReportPublisher(userId int64, StartDate, EndDate string) (output input.ReportPublisher, err error) {
//
// 	return
// }

func (t *PaymentRequest) ApiGetReportPublisher(userId int64, StartDate, EndDate string) (output input.ReportPublisher, err error) {
	user, err := GetUserById(userId)
	if err != nil {
		return
	}
	url := "https://apps.valueimpression.com/report/api-report/publisher/" + user.LoginToken + "?startDate=" + StartDate + "&endDate=" + EndDate
	var body []byte
	body, err = t.RequestUrl(url)
	if err != nil {
		return
	}
	json.Unmarshal(body, &output)
	output.StartDate = StartDate
	output.EndDate = EndDate
	if output.Revenue == "" {
		output.Revenue = "0"
	}
	if output.Rate == "" {
		output.Rate = "0"
	}
	return
}

func (t *PaymentRequest) RequestUrl(url string) (body []byte, err error) {
	c := colly.NewCollector(
	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	// colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	)
	// set Header to Request
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("apiTknPwv1", "Pwpo1undAPI")
	})
	// On every a element which has href attribute call callback
	c.OnResponse(func(r *colly.Response) {
		body = r.Body
	})
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
	})

	errVisit := c.Visit(url)
	if err == nil {
		err = errVisit
	}
	c.Wait()
	return
}

func (t *PaymentRequest) Validate(record input.InputPaymentRequestAdd) (errs []ajax.Error) {
	if record.UserId == 0 {
		errs = append(errs, ajax.Error{
			Id:      "user_id",
			Message: "User is required",
		})
	}
	if utility.ValidateString(record.StartDate) == "" {
		errs = append(errs, ajax.Error{
			Id:      "start_date",
			Message: "Start Date is required",
		})
	}
	if utility.ValidateString(record.EndDate) == "" {
		errs = append(errs, ajax.Error{
			Id:      "end_date",
			Message: "End Date is required",
		})
	}
	startDate, _ := time.Parse("2006-01-02", record.StartDate)
	start, _ := strconv.Atoi(startDate.Format("20060102"))
	endDate, _ := time.Parse("2006-01-02", record.EndDate)
	end, _ := strconv.Atoi(endDate.Format("20060102"))
	if end < start {
		errs = append(errs, ajax.Error{
			Id:      "end_date",
			Message: "End Date cannot be less than Start Date",
		})
	}

	newRequest, err := t.GetNewRequestByUserId(record.UserId)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	if newRequest.Id != 0 {
		if newRequest.EndDate.Unix() >= startDate.Unix() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: "Start Date được tạo request trước đó rồi",
			})
		}
	}

	return
}

func (t *PaymentRequest) GetAll() (Payments []PaymentRequestRecord) {
	mysql.Client.Find(&Payments)
	return
}

func (t *PaymentRequest) GetById(id int64) (record PaymentRequestRecord, err error) {
	err = mysql.Client.Where("id = ?", id).Find(&record).Error
	if record.Id == 0 {
		err = errors.New("Record not found")
	}
	return
}

type APIFirstDate struct {
	UserId    int64  `json:"userId"`
	FirstDate string `json:"firstDate"`
}

func (t *PaymentRequest) GetStartDateRequest(input input.GetStartDateUser) (starDate string, err error) {
	// nếu đã có request, lấy ngày kế tiếp ngày request đc tạo cuối cùng
	var request PaymentRequestRecord
	request, err = t.GetNewRequestByUserId(input.UserId)
	if err != nil {
		return
	}
	if request.Id != 0 {
		// date, _ := time.Parse("20060102", strconv.Itoa(request.EndDate))
		tomorrow := request.EndDate.AddDate(0, 0, +1)
		starDate = tomorrow.Format("2006-01-02")
		return
	}

	// nếu k có request, check xem có invoice không?
	var invoice PaymentInvoiceRecord
	invoice, err = new(PaymentInvoice).GetNewInvoiceByUserId(input.UserId)
	if invoice.Id != 0 {
		// date, _ := time.Parse("20060102", strconv.Itoa(invoice.EndDate))
		tomorrow := invoice.EndDate.AddDate(0, 0, +1)
		starDate = tomorrow.Format("2006-01-02")
		return
	}

	// nếu chưa có request thì lấy ngày đầu tiên có report
	url := "https://apps.valueimpression.com/report/api-report/date/?userId=" + strconv.FormatInt(input.UserId, 10)
	var body []byte
	body, err = t.RequestUrl(url)
	// xxx := make(map[string]interface{})
	firstDate := APIFirstDate{}
	json.Unmarshal(body, &firstDate)
	if err != nil {
		return
	}
	if firstDate.FirstDate == "" {
		err = fmt.Errorf("Users have no reports")
	}
	starDate = firstDate.FirstDate
	return
}

// get request mới nhất của user
func (t *PaymentRequest) GetNewRequestByUserId(userId int64) (record PaymentRequestRecord, err error) {
	err = mysql.Client.Where("user_id = ?", userId).Order("end_date DESC").Find(&record).Error
	return
}

// get all request commission pending
func (t *PaymentRequest) GetAllRequestCommissionPending() (records []PaymentRequestRecord, err error) {
	err = mysql.Client.Where("type = ? and status = ?", mysql.TypeCommission, mysql.StatusPaymentPending).Order("start_date ASC").Find(&records).Error
	return
}

// get all request prepaid có status = paid (đã thanh toán) của publisher
func (t *PaymentRequest) GetAllRequestPrepaidPaid(userId int64) (records []PaymentRequestRecord, err error) {
	err = mysql.Client.Where("status = ? and type = ? and user_id = ?", mysql.StatusPaymentPending, mysql.TypePrepaid, userId).Order("start_date ASC").Find(&records).Error
	return
}

// change status list requestId = done
func (t *PaymentRequest) HandleRequestsDone(ids interface{}) (err error) {
	err = mysql.Client.Model(&PaymentRequestRecord{}).Where("id in ?", ids).Update("status", 2).Error
	return
}

// get all request pending for agency
func (t *PaymentRequest) GetAllRequestPendingForAgency(agencyId int64) (records []PaymentRequestRecord, err error) {
	err = mysql.Client.Where("status = ? and user_id = ?", mysql.StatusPaymentPending, agencyId).Order("start_date ASC").Find(&records).Error
	return
}

// get request for invoice
func (t *PaymentRequest) GetAllRequestForInvoice(invoice PaymentInvoiceRecord) (records []PaymentRequestRecord, err error) {
	if invoice.RequestId == "" {
		return
	}
	requestIds := strings.Split(invoice.RequestId, ",")
	if len(requestIds) == 0 {
		return
	}
	err = mysql.Client.
		Where("id in ?", requestIds).
		Order("start_date ASC").
		Find(&records).Error
	return
}

type UserRecord struct {
	mysql.TableUser
}

func StringToSliceInt64(string, sep string) (slice []int64, err error) {
	arrayString := strings.Split(string, sep)
	if len(arrayString) == 0 {
		return
	}
	for _, value := range arrayString {
		int64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return slice, err
		}
		slice = append(slice, int64)
	}
	return
}

func GetUserById(id int64) (record UserRecord, err error) {
	err = mysql.Client.Where("id = ?", id).Find(&record).Error
	return
}

func getUsersById(id int64) (record []UserRecord, err error) {
	err = mysql.Client.Where("id = ?", id).Find(&record).Error
	return
}

func getAllPublisher() (records []UserRecord, err error) {
	err = mysql.Client.Where("permission in (?,?)", mysql.UserPermissionMember, mysql.UserPermissionManagedService).Find(&records).Error
	return
}

func getAllAgency() (records []UserRecord, err error) {
	err = mysql.Client.Where("permission = ?", mysql.UserPermissionSale).Find(&records).Error
	return
}

type UserBillingRecord struct {
	mysql.TableUserBilling
}

func (UserBillingRecord) TableName() string {
	return mysql.Tables.UserBilling
}

func GetBillingByUserId(userId int64) (record UserBillingRecord, err error) {
	err = mysql.Client.Model(&UserBillingRecord{}).Where("user_id = ?", userId).Find(&record).Error
	return
}

type RevenueShareRecord struct {
	mysql.TableRevenueShare
}

func (RevenueShareRecord) TableName() string {
	return mysql.Tables.RevenueShare
}

// get RevenueShare của publisher trc ngày startDate và trong thời gian StartDate => EndDate
func GetRevenueSharesForUser(userId int64, StartDate, EndDate string) (records []RevenueShareRecord, err error) {
	var record RevenueShareRecord
	record, err = GetRevenueShareOldForUser(userId, StartDate)
	if err != nil {
		return
	}
	if record.Id > 0 {
		records = append(records, record)
	}

	var _records []RevenueShareRecord
	err = mysql.Client.Model(&RevenueShareRecord{}).
		Where("user_id = ? and date >= ? and date <= ?", userId, StartDate, EndDate).
		Order("date ASC, id ASC").
		Find(&_records).Error

	records = append(records, _records...)
	return
}

func GetRevenueShareForUser(userId int64) (record RevenueShareRecord, err error) {
	user, _ := GetUserById(userId)
	err = mysql.Client.Model(&RevenueShareRecord{}).
		Where("user_id = ?", userId).
		Order("date DESC, id DESC").
		Find(&record).Error

	if err != nil {
		return
	}
	if record.Id == 0 {
		err = CreateRevenueShareForUser(userId, user.CreatedAt.AddDate(0, 0, 1).Format("2006-01-02"))
		if err != nil {
			return
		}
		err = mysql.Client.Model(&RevenueShareRecord{}).
			Where("user_id = ?", userId).
			Order("date DESC, id DESC").
			Find(&record).Error
	}

	return
}

func GetRevenueShareOldForUser(userId int64, dateNew string) (record RevenueShareRecord, err error) {
	err = mysql.Client.Model(&RevenueShareRecord{}).
		Where("user_id = ? and date < ?", userId, dateNew).
		Order("date DESC, id DESC").
		Find(&record).Error

	if err != nil {
		return
	}
	if record.Id == 0 {
		err = CreateRevenueShareForUser(userId, dateNew)
		if err != nil {
			return
		}
		err = mysql.Client.Model(&RevenueShareRecord{}).
			Where("user_id = ? and date < ?", userId, dateNew).
			Order("date ASC, id ASC").
			Find(&record).Error
	}
	return
}

func CreateRevenueShareForUser(userId int64, dateNew string) (err error) {

	Publisher, err := GetUserById(userId)
	if err != nil {
		return
	}
	if Publisher.CreatedAt.Format("2006-01-02") >= dateNew {
		return
	}
	var record RevenueShareRecord
	record.Rate = mysql.RevenueShareDefault
	record.UserId = userId
	record.Date = Publisher.CreatedAt
	mysql.Client.Model(&RevenueShareRecord{}).Create(&record)
	return
}

type NotificationRecord struct {
	mysql.TableNotification
}

func (NotificationRecord) TableName() string {
	return mysql.Tables.Notification
}

func CreateNotification(Message string) (err error) {
	var record NotificationRecord
	record.Status = 1
	record.UserId = 7
	record.Message = Message
	err = mysql.Client.Model(&NotificationRecord{}).Create(&record).Error
	return
}

func RemoveGuaranteeForUser(userID int64) {
	user, _ := GetUserById(userID)
	user.Guarantee = 0
	user.GuaranteeCeiling = ""
	user.GuaranteeFrom = ""
	user.GuaranteePeriods = ""
	mysql.Client.Where("id = ?", userID).Model(&UserRecord{}).Save(&user)
}

// func druid(query string, output BidderStatus) (Output BidderStatus, errs []ajax.Error) {
func druid(query string) (body []byte, err error) {
	URL := "https://query.vliplatform.com/druid/v2/sql/"
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("AUTHORIZATION", "Basic cmVhZG1hbjpMb3hpVVNKISEoODIhQA==")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})

	// On every a element which has href attribute call callback
	c.OnResponse(func(r *colly.Response) {
		body = r.Body
	})
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
			return
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
				return
			}
		}
	})
	message := map[string]interface{}{
		"query": query,
		"context": map[string]interface{}{
			"useCache":      "true",
			"populateCache": "true",
			"sqlTimeZone":   TIME_ZONE_DEFAULT,
		},
	}
	payload, err := json.Marshal(message)
	if err != nil {
		return
	}
	errVisit := c.PostRaw(URL, payload)
	if err == nil {
		err = errVisit
	}
	c.Wait()

	return
}

func PublishersForAgency(agencyId int64) (Publishers []UserRecord, err error) {
	err = mysql.Client.Model(&UserRecord{}).Where("status = 1 and presenter = ?", agencyId).Find(&Publishers).Error
	return
}
