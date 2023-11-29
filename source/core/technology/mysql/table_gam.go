package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableGam struct {
	Id           int64           `gorm:"column:id" json:"id"`
	UserId       int64           `gorm:"column:user_id" json:"user_id"`
	RefreshToken string          `gorm:"column:refresh_token" json:"refresh_token"`
	TokenType    string          `gorm:"column:token_type" json:"token_type"`
	GgEmail      string          `gorm:"column:gg_email" json:"gg_email"`
	GgUserId     int64           `gorm:"column:gg_user_id" json:"gg_user_id"`
	GgUserRole   string          `gorm:"column:gg_user_role" json:"gg_user_role"`
	IsDisabled   TYPEGamDisabled `gorm:"column:is_disabled" json:"is_disabled"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableGam) TableName() string {
	return Tables.Gam
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
