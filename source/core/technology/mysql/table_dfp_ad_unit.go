package mysql

import "time"

type TableDfpAdUnit struct {
	Id          int64               `gorm:"column:id" json:"id"`
	UserId      int64               `gorm:"column:user_id" json:"user_id"`
	InventoryId int64               `gorm:"column:inventory_id" json:"inventory_id"`
	TagId       int64               `gorm:"column:tag_id" json:"tag_id"`
	NetworkId   int64               `gorm:"column:network_id" json:"network_id"`
	NetworkName string              `gorm:"column:network_name" json:"network_name"`
	Version     int64               `gorm:"column:version" json:"version"`
	AdUnitName  string              `gorm:"column:ad_unit_name" json:"ad_unit_name"`
	AdUnitId    string              `gorm:"column:ad_unit_id" json:"ad_unit_id"`
	Status      TYPEStatusDfpAdUnit `gorm:"column:status" json:"status"`
	CreatedAt   time.Time           `gorm:"column:created_at" json:"created_at"`
}

type TYPEStatusDfpAdUnit string

const (
	TYPEStatusDfpAdUnitPending    = "pending"
	TYPEStatusDfpAdUnitProcessing = "processing"
	TYPEStatusDfpAdUnitSuccess    = "success"
	TYPEStatusDfpAdUnitError      = "error"
)
