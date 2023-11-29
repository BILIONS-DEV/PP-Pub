package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

func (BidderModel) TableName() string {
	return "bidder"
}

type BidderModel struct {
	ID               int64           `gorm:"column:id;primary_key" json:"id"`
	UserId           int64           `gorm:"column:user_id" json:"user_id"`
	BidderTemplateId int64           `gorm:"column:bidder_template_id" json:"bidder_template_id"`
	BidderCode       string          `gorm:"column:bidder_code" json:"bidder_code"`
	DisplayName      string          `gorm:"column:display_name" json:"display_name"`
	AliasName        string          `gorm:"column:alias_name;default:null" json:"alias_name"`
	ShowOnPub        string          `gorm:"column:show_on_pub;default:null" json:"show_on_pub"`
	BidderType       int             `gorm:"column:bidder_type" json:"bidder_type"`
	BidderAlias      string          `gorm:"column:bidder_alias" json:"bidder_alias"`
	AccountType      TYPEAccountType `gorm:"column:account_type" json:"account_type"`
	PubId            string          `gorm:"column:pub_id" json:"pub_id"`
	LinkedAccount    TYPEState       `gorm:"column:linked_account;default:2" json:"linked_account"`
	LinkedGam        int64           `gorm:"column:linked_gam" json:"linked_gam"`
	BidAdjustment    *float64        `gorm:"column:bid_adjustment" json:"bid_adjustment"`
	RPM              float64         `gorm:"column:rpm" json:"rpm"`
	IsChange         string          `gorm:"column:is_change" json:"is_change"`
	AdsTxt           string          `gorm:"column:ads_txt" json:"ads_txt"`
	Status           TYPEStatus      `gorm:"column:status" json:"status"`
	IsLock           TYPEIsLock      `gorm:"column:is_lock" json:"is_lock"`
	IsDefault        TYPEState       `gorm:"column:is_default" json:"is_default"`
	SellerDomain     string          `gorm:"column:seller_domain" json:"seller_domain"`
	XlsxPath         string          `gorm:"column:xlsx_path" json:"xlsx_path"`
	ScanAmz          int             `gorm:"column:scan_amz" json:"scan_amz"`
	LastScanAmz      time.Time       `gorm:"column:last_scan_amz" json:"last_scan_amz"`
	CreatedAt        time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`
}

func (b *BidderModel) IsFound() bool {
	if b.ID > 0 {
		return true
	}
	return false
}

func (b *BidderModel) ToJSON() string {
	jsonEncode, _ := json.Marshal(b)
	return string(jsonEncode)
}

type TYPEAccountType int

const (
	TYPEAccountTypeAdsense TYPEAccountType = iota + 1
	TYPEAccountTypeAdx
)

func (s TYPEAccountType) String() string {
	switch s {
	case TYPEAccountTypeAdsense:
		return "Adsense"
	case TYPEAccountTypeAdx:
		return "Adx"
	default:
		return ""
	}
}

func (s TYPEAccountType) Int() int64 {
	switch s {
	case TYPEAccountTypeAdsense:
		return 1
	case TYPEAccountTypeAdx:
		return 2
	default:
		return 0
	}
}

type TYPEIsLock int

const (
	TYPEIsLockTypeUnlock TYPEIsLock = iota + 1
	TYPEIsLockTypeLock
)

func (s TYPEIsLock) String() string {
	switch s {
	case TYPEIsLockTypeLock:
		return "lock"
	case TYPEIsLockTypeUnlock:
		return "unlock"
	default:
		return "unlock"
	}
}

func (s TYPEIsLock) Int() int {
	switch s {
	case TYPEIsLockTypeUnlock:
		return 1
	case TYPEIsLockTypeLock:
		return 2
	default:
		return 0
	}
}
