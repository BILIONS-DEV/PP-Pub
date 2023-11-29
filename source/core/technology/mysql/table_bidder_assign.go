package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableBidderAssign struct {
	Id            int64                  `gorm:"column:id" json:"id"`
	UserId        int64                  `gorm:"column:user_id" json:"user_id"`
	BidderId      int64                  `gorm:"column:bidder_id" json:"bidder_id"`
	BidAdjustment int                    `gorm:"column:bid_adjustment" json:"bid_adjustment"`
	RefreshAd     int                    `gorm:"column:refresh_ad" json:"refresh_ad"`
	AdFormat      ArrayInt64             `gorm:"column:ad_format;type:json" json:"ad_format"`
	Status        TYPEBidderAssignStatus `gorm:"column:status" json:"status"`
	GeoIpOption   TYPEOption             `gorm:"column:geo_ip_option" json:"geo_ip_option"`
	GeoIpList     ArrayString            `gorm:"column:geo_ip_list;type:json" json:"geo_ip_list"`
	AdsTxt        string                 `gorm:"column:ads_txt" json:"ads_txt"`
	CreatedAt     time.Time              `gorm:"column:created_at" json:"created_at"`
	DeletedAt     gorm.DeletedAt         `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableBidderAssign) TableName() string {
	return Tables.BidderAssign
}

type TYPEBidderAssignStatus int

const (
	BidderAssignStatusOn = iota + 1
	BidderAssignStatusOff
)

func (t TYPEBidderAssignStatus) String() string {
	switch t {
	case BidderAssignStatusOn:
		return "ON"
	case BidderAssignStatusOff:
		return "OFF"
	default:
		return ""
	}
}

func (t TYPEBidderAssignStatus) Standardized() TYPEBidderAssignStatus {
	if t == 1 {
		return BidderAssignStatusOn
	}
	return BidderAssignStatusOff
}

// IsFound Check User Exists
//
// return:
func (rec TableBidderAssign) IsFound() (flag bool) {
	if rec.Id > 0 {
		flag = true
	}
	return
}
