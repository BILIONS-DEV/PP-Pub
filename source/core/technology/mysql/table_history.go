package mysql

import (
	"time"
)

type TableHistory struct {
	Id          int64          `gorm:"column:id" json:"id"`
	CreatorId   int64          `gorm:"column:creator_id" json:"creator_id"`
	Page        string         `gorm:"column:page" json:"page"`
	Title       string         `gorm:"column:title" json:"title"`
	Object      string         `gorm:"column:object" json:"object"`
	ObjectId    int64          `gorm:"column:object_id" json:"object_id"`
	ObjectName  string         `gorm:"column:object_name" json:"object_name"`
	ObjectType  TYPEObjectType `gorm:"column:object_type" json:"object_type"`
	DetailType  string         `gorm:"column:detail_type" json:"detail_type"`
	App         string         `gorm:"column:app" json:"app"`
	UserId      int64          `gorm:"column:user_id" json:"user_id"`
	InventoryId int64          `gorm:"column:inventory_id" json:"inventory_id"`
	BidderId    int64          `gorm:"column:bidder_id" json:"bidder_id"`
	OldData     string         `gorm:"column:old_data" json:"old_data"`
	NewData     string         `gorm:"column:new_data" json:"new_data"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (TableHistory) TableName() string {
	return Tables.History
}

type TYPEObjectType int

const (
	TYPEObjectTypeAdd TYPEObjectType = iota + 1
	TYPEObjectTypeUpdate
	TYPEObjectTypeDel
)
