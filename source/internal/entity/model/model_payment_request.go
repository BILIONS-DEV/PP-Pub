package model

import (
	"github.com/syyongx/php2go"
	"time"
)

func (PaymemtRequest) TableName() string {
	return "payment_request"
}

type PaymemtRequest struct {
	Id        int64              `json:"id" form:"id"`
	Creator   int64              `json:"creator"`
	UserId    int64              `json:"user_id"`
	Revenue   float64            `json:"revenue"`
	Amount    float64            `json:"amount"`
	Rate      float64            `json:"rate"`
	Currency  string             `json:"currency"`
	Note      string             `json:"note"`
	Type      TypePaymentRequest `json:"type"`
	Status    StatusPayment      `json:"status"`
	StartDate time.Time          `json:"start_date"`
	EndDate   time.Time          `json:"end_date"`
	CreatedAt time.Time          `gorm:"column:created_at" json:"created_at"`
}

type TypePaymentRequest int

const (
	TypeCommission TypePaymentRequest = iota + 1
	TypePrepaid
	TypeAgencyProfit     // agency_profit
	TypeAgencyShare      // agency_share
	TypeAgencyGuarantee  // agency_guarantee
	TypeRevenueAffiliate // revenue_affiliate
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

func (r PaymemtRequest) StartDateString() string {
	// yourDate, err := time.Parse("20060102", strconv.Itoa(r.StartDate))
	// if err != nil {
	// 	return ""
	// }
	return r.StartDate.Format("01/02/2006") //
}

func (r PaymemtRequest) EndDateString() string {
	// yourDate, err := time.Parse("20060102", strconv.Itoa(r.EndDate))
	// if err != nil {
	// 	return ""
	// }
	return r.EndDate.Format("01/02/2006") //
}

func (r PaymemtRequest) GetAmountPaymentFormat() string {
	return "$" + php2go.NumberFormat(r.Amount, 2, ".", ",")
}

func (r PaymemtRequest) GetRevenuePaymentFormat() string {
	return "$" + php2go.NumberFormat(r.Revenue, 2, ".", ",")
}
