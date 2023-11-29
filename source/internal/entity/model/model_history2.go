package model

import (
	"time"
)

type HistoryModel struct {
	ID          int64             `gorm:"column:id"`
	CreatorID   int64             `gorm:"column:creator_id"`
	UserID      int64             `gorm:"column:user_id"`
	Title       string            `gorm:"column:title"`
	Page        string            `gorm:"column:page"`
	Object      string            `gorm:"column:object"`
	ObjectID    int64             `gorm:"column:object_id"`
	ObjectName  string            `gorm:"column:object_name"`
	ObjectType  HistoryObjectTYPE `gorm:"column:object_type"`
	DetailType  string            `gorm:"column:detail_type"`
	InventoryID int64             `gorm:"column:inventory_id" json:"inventory_id"`
	BidderID    int64             `gorm:"column:bidder_id" json:"bidder_id"`
	App         string            `gorm:"column:app"`
	OldData     string            `gorm:"column:old_data"`
	NewData     string            `gorm:"column:new_data"`
	CreatedAt   time.Time
}

type CompareData struct {
	Action  string
	Text    string
	OldData string
	NewData string
}

type HistoryDetailTypeTYPE string

const (
	DetailInventorySubmitFE  HistoryDetailTypeTYPE = "inventory_submit_fe"
	DetailInventoryConfigFE  HistoryDetailTypeTYPE = "inventory_config_fe"
	DetailInventoryConsentFE HistoryDetailTypeTYPE = "inventory_consent_fe"
	DetailInventoryAdsTxtFE  HistoryDetailTypeTYPE = "inventory_adstxt_fe"

	DetailBidderFE                 HistoryDetailTypeTYPE = "bidder_fe"
	DetailBidderAdxConnectionMcmBE HistoryDetailTypeTYPE = "bidder_adx_connection_mcm_be"
	DetailBidderBE                 HistoryDetailTypeTYPE = "bidder_be"
)

type HistoryObjectTYPE int

const (
	HistoryObjectTYPEAdd HistoryObjectTYPE = iota + 1
	HistoryObjectTYPEUpdate
	HistoryObjectTYPEDelete
)

func (t HistoryObjectTYPE) String() string {
	switch t {
	case HistoryObjectTYPEAdd:
		return "add"
	case HistoryObjectTYPEUpdate:
		return "update"
	case HistoryObjectTYPEDelete:
		return "delete"
	default:
		return ""
	}
}
