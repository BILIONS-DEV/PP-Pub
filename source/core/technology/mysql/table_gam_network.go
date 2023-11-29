package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableGamNetwork struct {
	Id               int64          `gorm:"column:id" json:"id"`
	UserId           int64          `gorm:"column:user_id" json:"user_id"`
	GamId            int64          `gorm:"column:gam_id" json:"gam_id"`
	NetworkId        int64          `gorm:"column:network_id" json:"network_id"`
	NetworkName      string         `gorm:"column:network_name" json:"network_name"`
	CurrencyCode     string         `gorm:"column:currency_code" json:"currency_code"`
	TimeZone         string         `gorm:"column:time_zone" json:"time_zone"`
	ApiAccess        TYPEApiAccess  `gorm:"column:api_access" json:"api_access"`
	Status           TYPEStatusGam  `gorm:"column:status" json:"status"`
	PushLineItem     int            `gorm:"column:push_line_item" json:"push_line_item"`
	DatePushLineItem time.Time      `gorm:"column:date_push_line_item" json:"date_push_line_item"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	User             TableUser      `gorm:"-"`
}

func (TableGamNetwork) TableName() string {
	return Tables.GamNetwork
}

type TYPEStatusGam int

const (
	StatusGamPending TYPEStatusGam = iota + 1
	StatusGamSelected
)

func (t TYPEStatusGam) String() string {
	switch t {
	case StatusGamPending:
		return "pending"
	case StatusGamSelected:
		return "selected"
	default:
		return ""
	}
}

type TYPEApiAccess int

const (
	ApiAccessEnable TYPEApiAccess = iota + 1
	ApiAccessDisable
)

func (t TYPEApiAccess) String() string {
	switch t {
	case ApiAccessEnable:
		return "enable"
	case ApiAccessDisable:
		return "disable"
	default:
		return ""
	}
}
