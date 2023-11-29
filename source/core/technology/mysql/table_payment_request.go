package mysql

import (
	"github.com/syyongx/php2go"
	"time"
)

type TablePaymentRequest struct {
	Id        int64              `gorm:"column:id" json:"id" form:"id"`
	Creator   int64              `gorm:"column:creator" json:"creator"`
	UserId    int64              `gorm:"column:user_id" json:"user_id"`
	Revenue   float64            `gorm:"column:revenue" json:"revenue"`
	Amount    float64            `gorm:"column:amount" json:"amount"`
	Rate      float64            `gorm:"column:rate" json:"rate"`
	Currency  string             `gorm:"column:currency" json:"currency"`
	Note      string             `gorm:"column:note" json:"note"`
	Type      TypePaymentRequest `gorm:"column:type" json:"type"`
	Status    StatusPayment      `gorm:"column:status" json:"status"`
	StartDate time.Time          `gorm:"column:start_date" json:"start_date"`
	EndDate   time.Time          `gorm:"column:end_date" json:"end_date"`
	CreatedAt time.Time          `gorm:"column:created_at" json:"created_at"`
}

type TypePaymentRequest int
type StatusPayment int

const (
	TypeCommission TypePaymentRequest = iota + 1
	TypePrepaid
	TypeAgencyProfit     // agency_profit
	TypeAgencyShare      // agency_share
	TypeAgencyGuarantee  // agency_guarantee
	TypeRevenueAffiliate // revenue_affiliate
)

const (
	StatusPaymentPending StatusPayment = iota + 1
	StatusPaymentDone
	StatusPaymentPrepaidPaid
	StatusPaymentRject
)

func (s TypePaymentRequest) String() string {
	switch s {
	case TypeCommission:
		return "Commission"
	case TypePrepaid:
		return "Prepaid"
	case TypeAgencyProfit:
		return "Agency Commission"
	case TypeAgencyShare:
		return "Agency Share"
	case TypeAgencyGuarantee:
		return "Guarantee"
	case TypeRevenueAffiliate:
		return "Revenue Affiliate"
	default:
		return ""
	}
}

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

func (r TablePaymentRequest) StartDateString() string {
	return r.StartDate.Format("01/02/2006") //
}

func (r TablePaymentRequest) EndDateString() string {
	return r.EndDate.Format("01/02/2006") //
}

func (r TablePaymentRequest) GetAmountPaymentFormat() string {
	return "$" + php2go.NumberFormat(r.Amount, 2, ".", ",")
}

func (r TablePaymentRequest) GetRevenuePaymentFormat() string {
	return "$" + php2go.NumberFormat(r.Revenue, 2, ".", ",")
}
