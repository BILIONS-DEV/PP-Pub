package mysql

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type TableLineItemV2 struct {
	Id              int64                       `json:"id" form:"id"`
	UserId          int64                       `json:"user_id"`
	Name            string                      `json:"name"`
	Description     string                      `json:"description"`
	ServerType      TYPEServerType              `json:"server_type"`
	Status          TYPEOnOff                   `json:"status"`
	IsLock          TYPEIsLock                  `json:"is_lock"`
	AutoCreate      TYPEOnOff                   `json:"auto_create"`
	ConnectionType  TYPEConnectionType          `json:"connection_type"`
	GamLineItemType TYPEGamLineItemType         `json:"gam_line_item_type"`
	LinkedGam       int64                       `json:"linked_gam"`
	Rate            int                         `json:"rate"`
	VastUrl         string                      `json:"vast_url"`
	AdTag           string                      `json:"ad_tag"`
	Priority        int                         `json:"priority"`
	StartDate       sql.NullTime                `json:"start_date"`
	EndDate         sql.NullTime                `json:"end_date"`
	LineItemType    int                         `json:"line_item_type"`
	ApdInventory    int64                       `json:"apd_inventory"`
	CreatedAt       time.Time                   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time                   `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt              `gorm:"column:deleted_at" json:"deleted_at"`
	Targets         []TableTarget               `gorm:"-" as:"-"`
	BidderInfo      []TableLineItemBidderInfoV2 `gorm:"-" as:"-"`
	BidderGoogle    TableBidder                 `gorm:"-" as:"-"`
	AdsenseAdSlots  []LineItemAdsenseAdSlot     `gorm:"-" as:"-"`
}

func (TableLineItemV2) TableName() string {
	return Tables.LineItemV2
}

func (rec *TableLineItemV2) GetRls() {
	// Get các rls của line item
	var targets []TableTarget
	Client.Where(TableTarget{LineItemId: rec.Id}).Find(&targets)
	rec.Targets = targets
	var bidderInfos []TableLineItemBidderInfoV2
	Client.Where("line_item_id = ?", rec.Id).Find(&bidderInfos)
	for k, bidderInfo := range bidderInfos {
		var params []TableLineItemBidderParamsV2
		Client.Model(&TableLineItemBidderParamsV2{}).Where("line_item_bidder_id = ?", bidderInfo.Id).Find(&params)
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
func (rec TableLineItemV2) IsFound() (flag bool) {
	if rec.Id > 0 {
		flag = true
	}
	return
}
