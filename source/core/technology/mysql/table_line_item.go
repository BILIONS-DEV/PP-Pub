package mysql

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type TableLineItem struct {
	Id              int64                     `json:"id" form:"id"`
	UserId          int64                     `json:"user_id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	ServerType      TYPEServerType            `json:"server_type"`
	Status          TYPEOnOff                 `json:"status"`
	IsLock          TYPEIsLock                `json:"is_lock"`
	AutoCreate      TYPEOnOff                 `json:"auto_create"`
	ConnectionType  TYPEConnectionType        `json:"connection_type"`
	GamLineItemType TYPEGamLineItemType       `json:"gam_line_item_type"`
	LinkedGam       int64                     `json:"linked_gam"`
	Rate            int                       `json:"rate"`
	VastUrl         string                    `json:"vast_url"`
	AdTag           string                    `json:"ad_tag"`
	Priority        int                       `json:"priority"`
	StartDate       sql.NullTime              `json:"start_date"`
	EndDate         sql.NullTime              `json:"end_date"`
	LineItemType    int                       `json:"line_item_type"`
	ApdInventory    int64                     `json:"apd_inventory"`
	CreatedAt       time.Time                 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time                 `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt            `gorm:"column:deleted_at" json:"deleted_at"`
	Targets         []TableTarget             `gorm:"-" as:"-"`
	BidderInfo      []TableLineItemBidderInfo `gorm:"-" as:"-"`
	BidderGoogle    TableBidder               `gorm:"-" as:"-"`
	AdsenseAdSlots  []LineItemAdsenseAdSlot   `gorm:"-" as:"-"`
}

func (TableLineItem) TableName() string {
	return Tables.LineItem
}

func (rec *TableLineItem) GetRls() {
	// Get các rls của line item
	var targets []TableTarget
	Client.Where(TableTarget{LineItemId: rec.Id}).Find(&targets)
	rec.Targets = targets
	var bidderInfos []TableLineItemBidderInfo
	Client.Where("line_item_id = ?", rec.Id).Find(&bidderInfos)
	for k, bidderInfo := range bidderInfos {
		var params []TableLineItemBidderParams
		Client.Model(&TableLineItemBidderParams{}).Where("line_item_bidder_id = ?", bidderInfo.Id).Find(&params)
		bidderInfo.Params = params
		bidderInfos[k] = bidderInfo
	}
	rec.BidderInfo = bidderInfos
	var lineItemAccount TableLineItemAccount
	Client.Where("line_item_id = ?", rec.Id).Last(&lineItemAccount)
	var bidderGoogle TableBidder
	Client.Where("id = ?", lineItemAccount.BidderId).First(&bidderGoogle)
	rec.BidderGoogle = bidderGoogle
	var adsenseAdSlots []LineItemAdsenseAdSlot
	Client.Where(LineItemAdsenseAdSlot{LineItemId: rec.Id}).Find(&adsenseAdSlots)
	rec.AdsenseAdSlots = adsenseAdSlots
	return
}

// IsFound Check User Exists
//
// return:
func (rec TableLineItem) IsFound() (flag bool) {
	if rec.Id > 0 {
		flag = true
	}
	return
}

type TYPEBidderStatus int
type TYPEBidderType int
type TYPEServerType int

const (
	TYPEBidderStatusOFF TYPEBidderStatus = iota
	TYPEBidderStatusReady
	TYPEBidderStatusInactive
	TYPEBidderStatusPaused
	TYPEBidderStatusDelivering
	TYPEBidderStatusCompleted
)

func (s TYPEBidderStatus) String() string {
	switch s {
	case TYPEBidderStatusOFF:
		return "OFF"
	case TYPEBidderStatusReady:
		return "Ready"
	case TYPEBidderStatusInactive:
		return "Inactive"
	case TYPEBidderStatusPaused:
		return "Paused"
	case TYPEBidderStatusDelivering:
		return "Delivering"
	case TYPEBidderStatusCompleted:
		return "Completed"
	default:
		return "OFF"
	}
}

const (
	TYPEBidderTypePrebidClient TYPEBidderType = iota + 1
	TYPEBidderTypePrebidServer
	TYPEBidderTypeVAST
	TYPEBidderTypeAdTag
)

func (s TYPEBidderType) String() string {
	switch s {
	case TYPEBidderTypePrebidClient:
		return "Prebid Client"
	case TYPEBidderTypePrebidServer:
		return "Prebid Server (S2S)"
	case TYPEBidderTypeVAST:
		return "VAST"
	case TYPEBidderTypeAdTag:
		return "Ad Tag"
	default:
		return ""
	}
}
func (s TYPEBidderType) Int() int {
	switch s {
	case TYPEBidderTypePrebidClient:
		return 1
	case TYPEBidderTypePrebidServer:
		return 2
	case TYPEBidderTypeVAST:
		return 3
	case TYPEBidderTypeAdTag:
		return 4
	default:
		return 0
	}
}

const (
	TYPEServerTypePrebid TYPEServerType = iota + 1
	TYPEServerTypeGoogle
)

func (s TYPEServerType) String() string {
	switch s {
	case TYPEServerTypePrebid:
		return "Prebid"
	case TYPEServerTypeGoogle:
		return "Google"
	default:
		return "Prebid"
	}
}

func (s TYPEServerType) Int() int {
	switch s {
	case TYPEServerTypePrebid:
		return 1
	case TYPEServerTypeGoogle:
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

type TYPEConnectionType int

const (
	TYPEConnectionTypeAdUnits TYPEConnectionType = iota + 1
	TYPEConnectionTypeLineItems
	TYPEConnectionTypeMCM
)

func (s TYPEConnectionType) String() string {
	switch s {
	case TYPEConnectionTypeAdUnits:
		return "Ad units"
	case TYPEConnectionTypeLineItems:
		return "Line items"
	case TYPEConnectionTypeMCM:
		return "MCM"
	default:
		return ""
	}
}

func (s TYPEConnectionType) Int() int {
	switch s {
	case TYPEConnectionTypeAdUnits:
		return 1
	case TYPEConnectionTypeLineItems:
		return 2
	case TYPEConnectionTypeMCM:
		return 3
	default:
		return 0
	}
}

type TYPEGamLineItemType int

const (
	TYPEGamLineItemTypeDisplay TYPEGamLineItemType = iota + 1
	TYPEGamLineItemTypeVideo
)

func (s TYPEGamLineItemType) String() string {
	switch s {
	case TYPEGamLineItemTypeDisplay:
		return "Display"
	case TYPEGamLineItemTypeVideo:
		return "Video"
	default:
		return ""
	}
}

func (s TYPEGamLineItemType) Int() int {
	switch s {
	case TYPEGamLineItemTypeDisplay:
		return 1
	case TYPEGamLineItemTypeVideo:
		return 2
	default:
		return 0
	}
}
