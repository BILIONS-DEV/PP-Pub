package mysql

import "gorm.io/gorm"

type LineItemAdsenseAdSlot struct {
	Id              int64  `gorm:"column:id" json:"id"`
	LineItemId      int64  `gorm:"column:line_item_id" json:"line_item_id"`
	Size            string `gorm:"column:size" json:"size"`
	AdsenseAdSlotId string `gorm:"column:adsense_ad_slot_id" json:"adsense_ad_slot_id"`
	DeletedAt       gorm.DeletedAt
}

func (LineItemAdsenseAdSlot) TableName() string {
	return Tables.LineItemAdsenseAdSlot
}
