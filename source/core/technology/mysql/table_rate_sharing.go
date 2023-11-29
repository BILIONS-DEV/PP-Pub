package mysql

import (
	"time"
)

type TableRateSharingInventory struct {
	Id             int64     `gorm:"column:id" json:"id"`
	InventoryId    int64     `gorm:"column:inventory_id" json:"inventory_id"`
	RateType       string    `gorm:"column:rate_type" json:"rate_type"`
	ApplicableDate int64     `gorm:"column:applicable_date" json:"applicable_date"`
	Rate           int64     `gorm:"column:rate" json:"rate"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (TableRateSharingInventory) TableName() string {
	return Tables.RateSharingInventory
}
