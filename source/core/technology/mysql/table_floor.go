package mysql

import (
	"math"
	"time"

	"gorm.io/gorm"
)

type TableFloor struct {
	Id          int64          `gorm:"column:id" json:"id"`
	Name        string         `gorm:"column:name" json:"name"`
	Description string         `gorm:"column:description" json:"description"`
	Status      int            `gorm:"column:status" json:"status"`
	FloorType   TYPEFloor      `gorm:"column:floor_type" json:"floor_type"`
	FloorValue  float64        `gorm:"column:floor_value" json:"floor_value"`
	Priority    int            `gorm:"column:priority" json:"priority"`
	UserId      int64          `gorm:"column:user_id" json:"user_id"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	Targets     []TableTarget  `gorm:"-"`
}

type TYPEFloor int

const (
	TYPEFloorDynamicFloor TYPEFloor = iota + 1
	TYPEFloorHardPriceFloor
	TYPEFloorPricingRule
)

func (TableFloor) TableName() string {
	return Tables.Floor
}

func (rec *TableFloor) GetRls() {
	// Get các rls
	var targets []TableTarget
	Client.Where(TableTarget{FloorId: rec.Id}).Find(&targets)
	rec.Targets = targets
	return
}

type TablePricingRulesGoogle struct {
	Id            int       `gorm:"id" json:"id"`
	NetworkId     int64     `gorm:"network_id" json:"network_id"`
	CustomValue   string    `gorm:"custom_value" json:"custom_value"`
	Name          string    `gorm:"name" json:"name"`
	Ecpm          float64   `gorm:"ecpm" json:"ecpm"`
	PricingRuleId int64     `gorm:"pricing_rule_id" json:"pricing_rule_id"`
	Active        bool      `gorm:"active" json:"active"`
	CreatedAt     time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (t *TablePricingRulesGoogle) EcpmFixed() float64 {
	return math.Round(t.Ecpm*100) / 100
}

type TablePricingRulesJobs struct {
	Id          int       `gorm:"id" json:"id"`
	GroupId     string    `gorm:"group_id" json:"group_id"`
	NetworkId   int64     `gorm:"network_id" json:"network_id"`
	CustomValue string    `gorm:"id" json:"custom_value"`
	Ecpm        float64   `gorm:"ecpm" json:"ecpm"`
	Status      string    `gorm:"status" json:"status"`
	Log         string    `gorm:"log" json:"log"`
	CreatedAt   time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt   time.Time `gorm:"deleted_at,omitempty" json:"deleted_at,omitempty"`

	Action string `json:"action" gorm:"-"`
	Count  int    `json:"count" gorm:"count"`
}

func (t *TablePricingRulesJobs) BuildAction() {
	t.Action = `<div data-group="` + t.GroupId + `" + data-status="` + t.Status + `" class="btn btn-dark-100 bg-gray-200 btn-icon btn-sm rounded-circle mx-1 view-more-floor-jobs" data-bs-toggle="modal" data-bs-target="#floor-price-jobs"><span data-bs-toggle="tooltip" title="View">
		<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" class="bi bi-eye" fill="currentColor">
		<path fill-rule="evenodd" d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8zM1.173 8a13.134 13.134 0 0 0 1.66 2.043C4.12 11.332 5.88 12.5 8 12.5c2.12 0 3.879-1.168 5.168-2.457A13.134 13.134 0 0 0 14.828 8a13.133 13.133 0 0 0-1.66-2.043C11.879 4.668 10.119 3.5 8 3.5c-2.12 0-3.879 1.168-5.168 2.457A13.133 13.133 0 0 0 1.172 8z"></path>
		<path fill-rule="evenodd" d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5zM4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0z"></path>
		</svg></span>`
}

func (t *TablePricingRulesJobs) AutoStatus() {
	class := "font-weight-bold"

	switch t.Status {
	case "error":
		class += " text-danger"
	case "success":
		class += " text-success"
	default:
	}
	t.Status = `<span class="` + class + `">` + t.Status + `</span>`
}
