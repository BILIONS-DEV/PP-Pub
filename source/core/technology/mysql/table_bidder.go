package mysql

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TableBidder struct {
	Id               int64                     `gorm:"column:id;primary_key" json:"id"`
	UserId           int64                     `gorm:"column:user_id" json:"user_id"`
	BidderTemplateId int64                     `gorm:"column:bidder_template_id" json:"bidder_template_id"`
	BidderCode       string                    `gorm:"column:bidder_code" json:"bidder_code"`
	DisplayName      string                    `gorm:"column:display_name" json:"display_name"`
	AliasName        string                    `gorm:"column:alias_name;default:null" json:"alias_name"`
	ShowOnPub        string                    `gorm:"column:show_on_pub;default:null" json:"show_on_pub"`
	BidderType       int                       `gorm:"column:bidder_type" json:"bidder_type"`
	BidderAlias      string                    `gorm:"column:bidder_alias" json:"bidder_alias"`
	AccountType      TYPEAccountType           `gorm:"column:account_type" json:"account_type"`
	PubId            string                    `gorm:"column:pub_id" json:"pub_id"`
	LinkedAccount    TypeOnOff                 `gorm:"column:linked_account;default:2" json:"linked_account"`
	LinkedGam        int64                     `gorm:"column:linked_gam" json:"linked_gam"`
	BidAdjustment    *float64                  `gorm:"column:bid_adjustment" json:"bid_adjustment"`
	RPM              float64                   `gorm:"column:rpm" json:"rpm"`
	IsChange         string                    `gorm:"column:is_change" json:"is_change"`
	StatusAdsTxt     string                    `gorm:"column:status_adstxt" json:"status_adstxt"`
	AdsTxt           string                    `gorm:"column:ads_txt" json:"ads_txt"`
	Status           TYPEStatus                `gorm:"column:status" json:"status"`
	IsLock           TYPEIsLock                `gorm:"column:is_lock" json:"is_lock"`
	IsDefault        TypeOnOff                 `gorm:"column:is_default" json:"is_default"`
	SellerDomain     string                    `gorm:"column:seller_domain" json:"seller_domain"`
	XlsxPath         string                    `gorm:"column:xlsx_path" json:"xlsx_path"`
	ScanAmz          int                       `gorm:"column:scan_amz" json:"scan_amz"`
	LastScanAmz      time.Time                 `gorm:"column:last_scan_amz" json:"last_scan_amz"`
	CreatedAt        time.Time                 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time                 `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt            `gorm:"column:deleted_at" json:"deleted_at"`
	MediaTypes       []TableRlsBidderMediaType `gorm:"foreignKey:BidderId;references:Id" as:"-"`
	Params           []TableBidderParams       `gorm:"foreignKey:BidderId;references:Id" as:"-"`
	GAM              TableGamNetwork           `gorm:"-" as:"-"`
	RlsConnectionMCM TableRlsConnectionMCM     `gorm:"-" as:"-"`
	ListStatusAdsTxt []string                  `gorm:"-" as:"-"`
}

func (rec *TableBidder) GetById(id int64) {
	Client.Preload(clause.Associations).Preload("MediaTypes").Preload("Params").Find(&rec, id)

	var gamNetwork TableGamNetwork
	Client.Find(&gamNetwork, rec.LinkedGam)
	rec.GAM = gamNetwork

}

func (TableBidder) TableName() string {
	return Tables.Bidder
}

type TYPEAccountType int

const (
	TYPEAccountTypeAdsense TYPEAccountType = iota + 1
	TYPEAccountTypeAdxMCM
	TYPEAccountTypeAdxMA
)

func (s TYPEAccountType) IsAdx() bool {
	switch s {
	case TYPEAccountTypeAdxMCM:
		return true
	case TYPEAccountTypeAdxMA:
		return true
	default:
		return false
	}
}

func (s TYPEAccountType) String() string {
	switch s {
	case TYPEAccountTypeAdsense:
		return "Adsense"
	case TYPEAccountTypeAdxMCM:
		return "Adx"
	case TYPEAccountTypeAdxMA:
		return "Adx MA"
	default:
		return ""
	}
}

func (s TYPEAccountType) StringObject() string {
	switch s {
	case TYPEAccountTypeAdsense:
		return "Adsense"
	case TYPEAccountTypeAdxMCM:
		return "Adx"
	case TYPEAccountTypeAdxMA:
		return "Adx"
	default:
		return ""
	}
}

func (s TYPEAccountType) Int() int64 {
	switch s {
	case TYPEAccountTypeAdsense:
		return 1
	case TYPEAccountTypeAdxMCM:
		return 2
	case TYPEAccountTypeAdxMA:
		return 3
	default:
		return 0
	}
}

type TYPEStatusAdsTxt string

const (
	TYPEStatusAdsTxtAll       TYPEStatusAdsTxt = "all"
	TYPEStatusAdsTxtApprove   TYPEStatusAdsTxt = "approve"
	TYPEStatusAdsTxtPending   TYPEStatusAdsTxt = "pending"
	TYPEStatusAdsTxtQueue     TYPEStatusAdsTxt = "queue"
	TYPEStatusAdsTxtSubmitted TYPEStatusAdsTxt = "submitted"
	TYPEStatusAdsTxtNotfound  TYPEStatusAdsTxt = "notfound"
	TYPEStatusAdsTxtReject    TYPEStatusAdsTxt = "reject"
	TYPEStatusAdsTxtNotUser   TYPEStatusAdsTxt = "not_user"
)

func (t TYPEStatusAdsTxt) StatusBidder() (stt []TYPERlsBidderSystemInventoryStatus) {
	switch t {
	case TYPEStatusAdsTxtAll:
		return
	case TYPEStatusAdsTxtApprove:
		stt = append(stt, RlsBidderSystemInventoryApproved)
		stt = append(stt, RlsBidderSystemInventoryApprovedClient)
		stt = append(stt, RlsBidderSystemInventoryApprovedS2S)
		return
	case TYPEStatusAdsTxtPending:
		stt = append(stt, RlsBidderSystemInventoryPending)
		return
	case TYPEStatusAdsTxtQueue:
		stt = append(stt, RlsBidderSystemInventoryQueue)
		return
	case TYPEStatusAdsTxtSubmitted:
		stt = append(stt, RlsBidderSystemInventorySubmitted)
		return
	case TYPEStatusAdsTxtNotfound:
		stt = append(stt, RlsBidderSystemInventoryNotfound)
		return
	case TYPEStatusAdsTxtReject:
		stt = append(stt, RlsBidderSystemInventoryRejected)
		return
	case TYPEStatusAdsTxtNotUser:
		stt = append(stt, RlsBidderSystemInventoryNotUse)
		return
	}
	return
}
