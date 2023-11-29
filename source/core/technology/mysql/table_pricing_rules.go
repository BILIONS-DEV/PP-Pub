package mysql

import (
	"time"
)

func (TablePricingRules) TableName() string {
	return Tables.PricingRule
}

type TablePricingRules struct {
	Id            int64     `gorm:"column:id" json:"id"`
	NetworkId     int64     `gorm:"column:network_id" json:"network_id"`
	CustomValue   string    `gorm:"column:custom_value" json:"custom_value"`
	Name          string    `gorm:"column:name" json:"name"`
	Countries     string    `gorm:"column:countries" json:"countries"`
	Devices       string    `gorm:"column:devices" json:"devices"`
	Tags          string    `gorm:"column:tags" json:"tags"`
	Ecpm          float64   `gorm:"column:ecpm" json:"ecpm"`
	PricingRuleId int64     `gorm:"column:pricing_rule_id" json:"pricing_rule_id"`
	Active        bool      `gorm:"column:active" json:"active"`
	Status        string    `gorm:"column:status" json:"status"`
	IdFloor       int64     `gorm:"column:id_floor" json:"id_floor"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
