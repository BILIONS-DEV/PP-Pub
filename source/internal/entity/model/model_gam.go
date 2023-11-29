package model

import (
	"gorm.io/gorm"
	"time"
)

type GamModel struct {
	ID           int64             `gorm:"column:id"`
	UserID       int64             `gorm:"column:user_id"`
	RefreshToken string            `gorm:"column:refresh_token"`
	TokenType    string            `gorm:"column:token_type"`
	GgEmail      string            `gorm:"column:gg_email"`
	GgUserID     int64             `gorm:"column:gg_user_id"`
	GgUserRole   string            `gorm:"column:gg_user_role"`
	IsDisabled   TYPEGamDisabled   `gorm:"column:is_disabled"`
	CreatedAt    time.Time         `gorm:"column:created_at"`
	DeletedAt    gorm.DeletedAt    `gorm:"column:deleted_at"`
	GamNetworks  []GamNetworkModel `gorm:"foreignKey:GamID;references:ID"`
}

func (GamModel) TableName() string {
	return "gam"
}

type TYPEGamDisabled int

const (
	GamEnable TYPEGamDisabled = iota + 1
	GamDisabled
)

func (t TYPEGamDisabled) String() string {
	switch t {
	case GamEnable:
		return "enable"
	case GamDisabled:
		return "disabled"
	default:
		return ""
	}
}
