package model

import "time"

func (History) TableName() string {
	return "history"
}

type History struct {
	ID          int64       `gorm:"column:id"`
	CreatorID   int64       `gorm:"column:creator_id"`
	Title       string      `gorm:"column:title"`
	Page        string      `gorm:"column:page"`
	Object      string      `gorm:"column:object"`
	ObjectID    int64       `gorm:"column:object_id"`
	ObjectName  string      `gorm:"column:object_name"`
	ObjectType  HistoryType `gorm:"column:object_type"`
	DetailType  string      `gorm:"column:detail_type"`
	App         string      `gorm:"column:app"`
	UserID      int64       `gorm:"column:user_id"`
	InventoryID int64       `gorm:"column:inventory_id"`
	BidderID    int64       `gorm:"column:bidder_id"`
	OldData     string      `gorm:"column:old_data"`
	NewData     string      `gorm:"column:new_data"`
	CreatedAt   time.Time   `gorm:"column:created_at"`
}

type HistoryType int

const (
	HistoryAdd HistoryType = iota + 1
	HistoryUpdate
)

func (t HistoryType) String() string {
	switch t {
	case HistoryAdd:
		return "add"
	case HistoryUpdate:
		return "update"
	}
	return ""
}
