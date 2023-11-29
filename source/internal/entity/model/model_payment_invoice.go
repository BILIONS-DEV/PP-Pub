package model

import (
	"source/pkg/datatable"
	"time"
)

func (PaymentInvoice) TableName() string {
	return "payment_invoice"
}

type PaymentInvoice struct {
	ID          int64         `json:"id" form:"id"`
	PayerID     int64         `json:"payer_id"`
	UserID      int64         `json:"user_id" gorm:" Foreignkey: UserName; reference: name "`
	RequestID   string        `json:"request_id"`
	Amount      float64       `json:"amount"`
	Rate        int64         `json:"rate"`
	Status      StatusPayment `json:"status"`
	Method      string        `json:"method"`
	Billing     string        `json:"billing"`
	Note        string        `json:"note"`
	Pdf         string        `json:"pdf"`
	Statement   string        `json:"statement"`
	StartDate   time.Time     `json:"start_date"`
	EndDate     time.Time     `json:"end_date"`
	PaymentDate time.Time     `json:"payment_date"`
	PaidDate    time.Time     `json:"paid_date"`
	UpdatedAt   time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

type PaymentInvoicesIndex struct {
	Invoice PaymentInvoice
	Id      string
	User    string
	Amount  string
	Method  string
	Period  string
	Status  string
	Action  string
	Info    string
}

type ParamPaymentIndex struct {
	datatable.Request
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	FromMoney   string   `query:"f_from_money" json:"f_from_money" form:"f_from_money"`
	ToMoney     string   `query:"f_to_money" json:"f_to_money" form:"f_to_money"`
	PaidDate    string   `query:"f_paid_date" form:"f_paid_date" json:"f_paid_date"`
	Month       string   `query:"f_month" form:"f_month" json:"f_month"`
	Sort        string   `query:"f_sort" form:"f_sort" json:"f_sort"`
	Permission  string   `query:"f_permission" form:"f_permission" json:"f_permission"`
	Page        string   `query:"f_page" form:"f_page" json:"f_page"`
	Size        string   `query:"f_size" form:"f_size" json:"f_size"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
}

type PageData struct {
	ListPage            []Pages
	FirstPage, LastPage string
}

type Pages struct {
	LinkPage string
	Number   int
}

const (
	PageSize = 10
	PageShow = 5
)

type StatusPayment int

const (
	StatusPaymentPending StatusPayment = iota + 1
	StatusPaymentDone
	StatusPaymentPrepaidPaid
	StatusPaymentRject
)

func (s StatusPayment) String() string {
	switch s {
	case StatusPaymentPending:
		return "Pending"
	case StatusPaymentDone:
		return "Done"
	case StatusPaymentPrepaidPaid:
		return "Prepaid Paid"
	case StatusPaymentRject:
		return "Reject"
	default:
		return ""
	}
}

type ReportAgency struct {
	GrossRevenue float64 `json:"grossRevenue" form:"grossRevenue"`
	NetRevenue   float64 `json:"netRevenue" form:"netRevenue"`
}

func CheckTimeEmpty(t time.Time) bool {
	return t == time.Time{}
}
