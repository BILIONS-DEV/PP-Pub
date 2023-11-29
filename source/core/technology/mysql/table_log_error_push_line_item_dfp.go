package mysql

import (
	"time"
)

type TableLogErrorPushLineItemDfp struct {
	Id          int64     `gorm:"column:id" json:"id"`
	GamId       int64     `gorm:"column:gam_id" json:"gam_id"`
	NetworkId   int64     `gorm:"column:network_id" json:"network_id"`
	NetworkName string    `gorm:"column:network_name" json:"network_name"`
	DomainId    int64     `gorm:"column:domain_id" json:"domain_id"`
	TagId       int64     `gorm:"column:tag_id" json:"tag_id"`
	LineItemId  int64     `gorm:"column:line_item_id" json:"line_item_id"`
	Value       string    `gorm:"column:value" json:"value"`
	Type        string    `gorm:"column:type" json:"type"`
	Log         string    `gorm:"column:log" json:"log"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}
