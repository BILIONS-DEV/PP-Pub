package mysql

import (
	"github.com/syyongx/php2go"
	"time"
)

type TablePaymentSubPub struct {
	Id          int64         `gorm:"column:id" json:"id" form:"id"`
	UserId      int64         `gorm:"column:user_id" json:"user_id"`
	Amount      float64       `gorm:"column:amount" json:"amount"`
	Status      StatusPayment `gorm:"column:status" json:"status"`
	Method      string        `gorm:"column:method" json:"method"`
	Billing     string        `gorm:"column:billing" json:"billing"`
	Note        string        `gorm:"column:note" json:"note"`
	StartDate   time.Time     `gorm:"column:start_date" json:"start_date"`
	EndDate     time.Time     `gorm:"column:end_date" json:"end_date"`
	PaymentDate time.Time     `gorm:"column:payment_date" json:"payment_date"`
	PaidDate    time.Time     `gorm:"column:paid_date" json:"paid_date"`
	UpdatedAt   time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

func (TablePaymentSubPub) TableName() string {
	return Tables.PaymentSubPub
}

func (row TablePaymentSubPub) PaymentDateString() string {
	// 2006/01/02 <=> y/m/d
	return row.PaymentDate.Format("2006-01-02") //
}

func (row TablePaymentSubPub) DatePaid() string {
	return row.PaidDate.Format("2006-01-02")
}

func (row TablePaymentSubPub) StartDateString() string {
	return row.StartDate.Format("2006-01-02") //
}

func (row TablePaymentSubPub) StartDateInvoiceString() string {
	return row.StartDate.Format("01-02-2006") //
}

func (row TablePaymentSubPub) EndDateString() string {
	return row.EndDate.Format("2006-01-02") //
}

func (row TablePaymentSubPub) DatePaymentYmd() string {
	return row.PaymentDate.Format("2006-01-02") //
}

func (row TablePaymentSubPub) GetAmountPayment() string {
	return "$" + php2go.NumberFormat(row.Amount, 2, ".", ",")
}
