package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableAdTag struct {
	Id         int64          `gorm:"column:id" json:"id"`
	Name       string         `gorm:"column:name" json:"name"`
	Type       int            `gorm:"column:type" json:"ad_tag_type"`
	SizeId     int            `gorm:"column:size_id" json:"ad_tag_size"`
	Gam        string         `gorm:"column:gam" json:"gam"`
	FloorPrice float64        `gorm:"column:floor_price" json:"floor_price"`
	PassBack   string         `gorm:"column:pass_back" json:"pass_back"`
	Status     int            `gorm:"column:status" json:"status"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableAdTag) TableName() string {
	return Tables.AdTag
}
