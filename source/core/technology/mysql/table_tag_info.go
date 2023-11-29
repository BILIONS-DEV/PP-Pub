package mysql

import (
	"time"
)

type TableTagInfo struct {
	TagId         int64     `gorm:"column:id" json:"tagId"`
	AdFormat      string    `gorm:"column:adFormat" json:"adFormat"`
	InventoryId   int64     `gorm:"column:inventoryId" json:"inventoryId"`
	InventoryType string    `gorm:"column:inventoryType" json:"inventoryType"`
	PubId         int64     `gorm:"column:pubId" json:"pubId"`
	AccManageId   int64     `gorm:"column:accManageId" json:"accManageId"`
	CreatedAt     time.Time `gorm:"column:createTime" json:"createTime"`
	UpdatedAt     time.Time `gorm:"column:updateTime" json:"updateTime"`
	SyncStatus    string    `gorm:"column:syncStatus" json:"syncStatus"`
}

func (TableTagInfo) TableName() string {
	return Tables.TagInfo
}
