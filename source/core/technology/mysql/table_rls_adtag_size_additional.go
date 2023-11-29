package mysql

import "gorm.io/gorm"

type TableRlsAdTagSizeAdditional struct {
	Id        int64          `gorm:"column:id" json:"id"`
	AdTagId   int64          `gorm:"column:ad_tag_id" json:"ad_tag_id"`
	Device    int64          `gorm:"column:device" json:"device"`
	AdSizeId  int64          `gorm:"column:ad_size_id" json:"ad_size_id"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableRlsAdTagSizeAdditional) TableName() string {
	return Tables.RlAdTagSizeAdditional
}
