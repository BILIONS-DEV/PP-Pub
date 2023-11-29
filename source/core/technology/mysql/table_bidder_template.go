package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableBidderTemplate struct {
	Id            int64                       `gorm:"column:id" json:"id"`
	UserId        int64                       `gorm:"column:user_id" json:"user_id"`
	Name          string                      `gorm:"column:name" json:"name"`
	PrebidModule  string                      `gorm:"column:prebid_module" json:"prebid_module"`
	BidderCode    string                      `gorm:"column:bidder_code" json:"bidder_code"`
	BidderAlias   string                      `gorm:"column:bidder_alias" json:"bidder_alias"`
	BidAdjustment float32                     `gorm:"column:bid_adjustment" json:"bid_adjustment"`
	DisplayName   string                      `gorm:"column:display_name" json:"display_name"`
	BidderType    int                         `gorm:"column:bidder_type" json:"bidder_type"`
	Status        TYPEStatus                  `gorm:"column:status" json:"status"`
	Logo          string                      `gorm:"column:logo" json:"logo"`
	CreatedAt     time.Time                   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time                   `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt              `gorm:"column:deleted_at" json:"deleted_at"`
	MediaTypes    []TableRlsBidderMediaType   `gorm:"-"`
	Params        []TableBidderTemplateParams `gorm:"-"`
}

func (TableBidderTemplate) TableName() string {
	return Tables.BidderTemplate
}

func (rec *TableBidderTemplate) GetRls() {
	// Get các rls
	var mediaTypes []TableRlsBidderMediaType
	Client.Where("bidder_template_id = ?", rec.Id).Find(&mediaTypes)
	rec.MediaTypes = mediaTypes

	var params []TableBidderTemplateParams
	Client.Where("bidder_template_id = ?", rec.Id).Find(&params)
	rec.Params = params
	return
}
