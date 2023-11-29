package invoice

import (
	"source/core/technology/mysql"
	"source/pkg/ajax"

	// "source/repository/payment/invoice/pdf"
	"strconv"
	"strings"
	"time"
)

type PaymentInvoice struct{}

type PaymentInvoiceRecord struct {
	mysql.TablePaymentInvoice
}

func (PaymentInvoiceRecord) TableName() string {
	return mysql.Tables.PaymentInvoice
}

func (t *PaymentInvoice) ScanInvoicePublisher() (errs []ajax.Error) {
	requests, err := new(PaymentRequest).GetAllRequestCommissionPending()
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	if len(requests) == 0 {
		return
	}
	UserRequests := make(map[int64][]PaymentRequestRecord)
	for _, value := range requests {
		UserRequests[value.UserId] = append(UserRequests[value.UserId], value)
	}
	if len(UserRequests) == 0 {
		return
	}
	for userId, userRequests := range UserRequests {
		user, err := GetUserById(userId)
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
			continue
		}

		startDate := time.Time{}
		endDate := time.Time{}
		totalPrepaid := 0.0
		totalRequestAmount := 0.0
		requestIds := []int64{}
		notes := ""
		for _, request := range userRequests {
			if startDate.Unix() <= 0 || startDate.Unix() > request.StartDate.Unix() {
				startDate = request.StartDate
			}
			if endDate.Unix() < request.EndDate.Unix() {
				endDate = request.EndDate
			}
			totalRequestAmount = totalRequestAmount + request.Amount
			requestIds = append(requestIds, request.Id)
			if request.Note != "" {
				if notes == "" {
					notes = request.Note
				} else {
					notes = notes + "\n" + request.Note
				}
			}
		}
		requestsPrepaid, err_ := new(PaymentRequest).GetAllRequestPrepaidPaid(user.Id)
		if err_ != nil {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err_.Error(),
			})
			continue
		}
		if len(requestsPrepaid) > 0 {
			for _, prepaid := range requestsPrepaid {
				totalPrepaid = totalPrepaid + prepaid.Amount
				requestIds = append(requestIds, prepaid.Id)
			}
		}
		if totalRequestAmount-totalPrepaid <= 0 {
			continue
		}

		// revenueShare, err := GetRevenueShareForUser(user.Id)
		// if err != nil {
		// 	errs = append(errs, ajax.Error{
		// 		Id:      "",
		// 		Message: err.Error(),
		// 	})
		// 	continue
		// }
		// if user.SystemSync == 1 && (user.Permission == 1 || user.Permission == 6) {
		// 	revenueShare.Rate = 100
		// }
		// amount := (totalRequestAmount - totalPrepaid) * float64(revenueShare.Rate) / float64(100)
		amount := totalRequestAmount - totalPrepaid
		if amount < 0.01 {
			continue
		}

		var row PaymentInvoiceRecord
		row.UserId = user.Id
		row.RequestId = Int64ToString(requestIds)
		row.Amount = amount
		// row.Rate = revenueShare.Rate
		row.Status = mysql.StatusPaymentPending
		row.Note = notes
		row.StartDate = startDate
		row.EndDate = endDate
		row.PaymentDate = endDate.AddDate(0, 0, (user.PaymentNet + 1))
		err = mysql.Client.Create(&row).Error
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
			continue
		}

		// change status list requestId = done
		err = new(PaymentRequest).HandleRequestsDone(requestIds)
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
			continue
		}
	}

	return
}

func (t *PaymentInvoice) ScanInvoiceAgency() (errs []ajax.Error) {
	agencies, err := getAllAgency()
	// agencies, err := getUsersById(137)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "There are no agencies!",
		})
		return
	}
	for _, agency := range agencies {
		// get all request pending for agency
		requests, err := new(PaymentRequest).GetAllRequestPendingForAgency(agency.Id)
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "request",
				Message: err.Error(),
			})
		}
		if len(requests) == 0 {
			continue
		}

		var Invoice PaymentInvoiceRecord

		var agencyProfit float64     // AgencyProfit
		var agencyShare float64      // commission of manager receive agency
		var revenueAffiliate float64 // commission of manager receive agency
		var agencyGuarantee float64  // agencyGuarantee = Guarantee - revenuePublisher: guarantee agency phải chịu khi publisher đc nhận số tiền guarantee mà commission của pub < guarantee
		for _, value := range requests {
			switch value.Type {
			case mysql.TypeAgencyProfit:
				agencyProfit = agencyProfit + value.Amount
				Invoice.Rate = int64(value.Rate)
				break
			case mysql.TypeAgencyShare:
				agencyShare = agencyShare + value.Amount
				break
			case mysql.TypeAgencyGuarantee:
				agencyGuarantee = agencyGuarantee + value.Amount
				break
			case mysql.TypeRevenueAffiliate:
				revenueAffiliate = revenueAffiliate + value.Amount
				break
			}
			if Invoice.StartDate.IsZero() == true {
				Invoice.StartDate = value.StartDate
			} else if Invoice.StartDate.Format("2006-01-02") > value.StartDate.Format("2006-01-02") {
				Invoice.StartDate = value.StartDate
			}
			if Invoice.EndDate.IsZero() == true {
				Invoice.EndDate = value.EndDate
			} else if Invoice.EndDate.Format("2006-01-02") > value.EndDate.Format("2006-01-02") {
				Invoice.EndDate = value.EndDate
			}
			if Invoice.RequestId == "" {
				Invoice.RequestId = strconv.FormatInt(value.Id, 10)
			} else {
				Invoice.RequestId = Invoice.RequestId + "," + strconv.FormatInt(value.Id, 10)
			}
		}

		commission := agencyProfit + agencyShare - revenueAffiliate - agencyGuarantee
		if commission < 0.01 {
			continue
		}
		if Invoice.Rate == 0 {
			if agency.Email == "davidwright@valueimpression.com" {
				Invoice.Rate = 15
			}
			if agency.Email == "l.watson@valueimpression.com" {
				Invoice.Rate = 10
			}
		}
		Invoice.Amount = commission
		Invoice.UserId = agency.Id
		Invoice.Status = mysql.StatusPaymentPending
		Invoice.PaymentDate = Invoice.EndDate.AddDate(0, 0, 60)
		err = mysql.Client.Create(&Invoice).Error
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "create invoice",
				Message: err.Error(),
			})
			continue
		}
		requestIds, err := StringToSliceInt64(Invoice.RequestId, ",")
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "string convert []int64",
				Message: err.Error(),
			})
			return
		}
		err = new(PaymentRequest).HandleRequestsDone(requestIds)
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "handle request done",
				Message: err.Error(),
			})
			return
		}
	}
	return
}

// func (t *PaymentInvoice) BuildInvoicePDF() (err error) {
// 	_, filename, _, ok := runtime.Caller(0)
// 	if !ok {
// 		log.Fatalln("No caller information")
// 	}
// 	pkgPath := path.Dir(filename)
//
// 	c := exec.Command("php", pkgPath+"/pdfinvoice.php")
// 	out, err := c.Output()
// 	if err != nil {
// 		return
// 	}
//
// 	type Response struct {
// 		Status  bool   `json:"status"`
// 		Message string `json:"message"`
// 	}
// 	var response Response
// 	err = json.Unmarshal(out, &response)
// 	if err == nil && response.Status == false && response.Message != "" {
// 		CreateNotification(response.Message)
// 	}
// 	return
// }

func (t *PaymentInvoice) BuildInvoicePDF() (err error) {
	// get invoices not pdf
	invoices, err := t.GetInvoicesNotPDF()
	if len(invoices) == 0 || err != nil {
		return
	}
	for _, value := range invoices {
		var user UserRecord
		user, err = GetUserById(value.UserId)
		if err != nil || user.Id == 0 {
			return
		}
		var billing UserBillingRecord
		if billing, err = GetBillingByUserId(user.Id); err != nil {
			return
		}
		var paymentRequest []PaymentRequestRecord
		if paymentRequest, err = new(PaymentRequest).GetAllRequestForInvoice(value); err != nil {
			return
		}
		var pathPDF string
		if pathPDF, err = GeneratePDF(value, paymentRequest, user, billing); err != nil {
			return
		}
		if value.Status == 2 {
			err = t.UpdateLinkInvoicePDF(pathPDF, value.Id)
		} else if value.Status == 1 {
			err = t.UpdateLinkStatementPDF(pathPDF, value.Id)
		}
		if err != nil {
			return
		}
	}

	return
}

func (t *PaymentInvoice) GetInvoicesNotPDF() (records []PaymentInvoiceRecord, err error) {
	err = mysql.Client.
		Where("(status = 2 AND pdf = \"\") OR (status = 1 AND statement = \"\")").
		Find(&records).Error
	return
}

func (t *PaymentInvoice) GetById(id int64) (record PaymentInvoiceRecord, err error) {
	err = mysql.Client.Where("id = ?", id).Find(&record).Error
	return
}

func (t *PaymentInvoice) GetByUserId(userId int64) (record []PaymentInvoiceRecord, err error) {
	err = mysql.Client.Where("user_id = ?", userId).Find(&record).Error
	return
}

// get payment mới nhất
func (t *PaymentInvoice) GetNewInvoiceByUserId(userId int64) (record PaymentInvoiceRecord, err error) {
	err = mysql.Client.Where("user_id = ?", userId).Order("end_date DESC").Find(&record).Error
	return
}

func (t *PaymentInvoice) UpdateLinkInvoicePDF(pdf string, invoiceID int64) (err error) {
	err = mysql.Client.Model(&PaymentInvoiceRecord{}).Where("id = ?", invoiceID).Update("pdf", pdf).Error
	return
}

func (t *PaymentInvoice) UpdateLinkStatementPDF(pdf string, invoiceID int64) (err error) {
	err = mysql.Client.Model(&PaymentInvoiceRecord{}).Where("id = ?", invoiceID).Update("statement", pdf).Error
	return
}

func Int64ToString(a []int64) string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(b, ",")
}
