package mysql

import (
	"github.com/syyongx/php2go"
	"time"
)

type TablePaymentInvoice struct {
	Id          int64         `gorm:"column:id" json:"id" form:"id"`
	PayerId     int64         `gorm:"column:payer_id" son:"payer_id"`
	UserId      int64         `gorm:"column:user_id" json:"user_id"`
	RequestId   string        `gorm:"column:request_id" json:"request_id"`
	Amount      float64       `gorm:"column:amount" json:"amount"`
	Rate        int64         `gorm:"column:rate" json:"rate"`
	Status      StatusPayment `gorm:"column:status" json:"status"`
	Method      string        `gorm:"column:method" json:"method"`
	Billing     string        `gorm:"column:billing" json:"billing"`
	Note        string        `gorm:"column:note" json:"note"`
	Pdf         string        `gorm:"column:pdf" json:"pdf"`
	Statement   string        `gorm:"column:statement" json:"statement"`
	StartDate   time.Time     `gorm:"column:start_date" json:"start_date"`
	EndDate     time.Time     `gorm:"column:end_date" json:"end_date"`
	PaymentDate time.Time     `gorm:"column:payment_date" json:"payment_date"`
	PaidDate    time.Time     `gorm:"column:paid_date" json:"paid_date"`
	UpdatedAt   time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

func (row TablePaymentInvoice) PaymentDateString() string {
	// 2006/01/02 <=> y/m/d
	return row.PaymentDate.Format("01/02/2006") //
}

func (row TablePaymentInvoice) DatePaid() string {
	return row.PaidDate.Format("01/02/2006")
}

func (row TablePaymentInvoice) StartDateString() string {
	return row.StartDate.Format("01/02/2006") //
}

func (row TablePaymentInvoice) StartDateInvoiceString() string {
	return row.StartDate.Format("01-02-2006") //
}

func (row TablePaymentInvoice) EndDateString() string {
	return row.EndDate.Format("01/02/2006") //
}

func (row TablePaymentInvoice) DatePaymentYmd() string {
	return row.PaymentDate.Format("2006-01-02") //
}

func (row TablePaymentInvoice) GetAmountPayment() string {
	return "$" + php2go.NumberFormat(row.Amount, 2, ".", ",")
}
