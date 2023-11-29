package mysql

import "time"

type TableDfpCreative struct {
	Id          int64
	UserId      int64           `gorm:"column:user_id" json:"user_id"`
	NetworkId   int64           `gorm:"column:network_id" json:"network_id"`
	NetworkName string          `gorm:"column:network_name" json:"network_name"`
	Type        TypeDfpCreative `gorm:"column:type" json:"type"`
	AccountName string          `gorm:"column:account_name" json:"account_name"`
	AccountType string          `gorm:"column:account_type" json:"account_type"`
	CreativeId  string          `gorm:"column:creative_id" json:"creative_id"`
	PubId       string          `gorm:"column:pub_id" json:"pub_id"`
	AdSlotId    string          `gorm:"column:ad_slot_id" json:"ad_slot_id"`
	Size        string          `gorm:"column:size" json:"size"`
	CreatedAt   time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   time.Time       `gorm:"column:deleted_at" json:"deleted_at"`
}

type TypeDfpCreative int

const (
	TypeDfpCreativeGoogle TypeDfpCreative = iota + 1
	TypeDfpCreativePrebid
)

func (t TypeDfpCreative) String() string {
	switch t {
	case TypeDfpCreativeGoogle:
		return "Google"
	case TypeDfpCreativePrebid:
		return "Prebid"
	}
	return ""
}
