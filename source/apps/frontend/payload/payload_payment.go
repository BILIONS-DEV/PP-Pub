package payload

import (
	"source/pkg/datatable"
	"source/core/technology/mysql"
)

type ParamPaymentIndex struct {
	InputPaymentFilter
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
}

type InputPaymentFilter struct {
	datatable.Request
	PostData *PostDataPaymentRequestFilter `query:"postData"`
}

type InputPaymentRequestFilter struct {
	datatable.Request
	PostData *PostDataPaymentRequestFilter `query:"postData"`
}

type PostDataPaymentRequestFilter struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	// datatable input
	Page   int `query:"page" json:"page" form:"page"`
	Limit  int `query:"limit" json:"limit" form:"limit"`
	Start  int `query:"start" json:"start" form:"start"`
	Length int `query:"length" json:"length" form:"length"`
}

type InputPaymentRequestAdd struct {
	UserId    int64  `form:"user_id" json:"user_id"`
	Note      string `form:"note" json:"note"`
	StartDate string `form:"start_date" json:"start_date"`
	EndDate   string `form:"end_date" json:"end_date"`
}

type InputPaymentInvoiceUpdate struct {
	Id       int64            `form:"id" json:"id"`
	Status   mysql.StatusPayment `form:"status" json:"status"`
	Note     string           `form:"note" json:"note"`
	PaidDate string           `form:"paid_date" json:"paid_date"`
	// PayerId  int64  `form:"payer_id" json:"payer_id"`
	// Method   string `form:"method" json:"method"`
	// Billing  string `form:"billing" json:"billing"`
}

type ReportPublisher struct {
	Id        string `json:"id"`
	Revenue   string `json:"revenue"`
	Currency  string `json:"currency"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type InputPaymentRequestDelete struct {
	Id int64
}

type GetStartDateUser struct {
	UserId int64 `form:"user_id" json:"user_id"`
}
