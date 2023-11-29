package mysql

import "gorm.io/gorm"

type TableBlockingRestrictions struct {
	Id              int64          `gorm:"column:id" json:"id"`
	BlockingId      int64          `gorm:"column:blocking_id" json:"blocking_id"`
	AdvertiseDomain string         `gorm:"column:advertiser_domain" json:"advertiser_domain"`
	CreativeId      string         `gorm:"column:creative_id" json:"creative_id"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableBlockingRestrictions) TableName() string {
	return Tables.BlockingRestrictions
}