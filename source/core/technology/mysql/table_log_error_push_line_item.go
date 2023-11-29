package mysql

import (
	"time"
)

type TableLogErrorPushLineItem struct {
	Id          int64     `gorm:"column:id" json:"id"`
	GamId       int64     `gorm:"column:gam_id" json:"gam_id"`
	NetworkId   int64     `gorm:"column:network_id" json:"network_id"`
	NetworkName string    `gorm:"column:network_name" json:"network_name"`
	Type        string    `gorm:"column:type" json:"type"`
	Index       int       `gorm:"column:idx" json:"idx"`
	Command     string    `gorm:"column:command" json:"command"`
	Log         string    `gorm:"column:log" json:"log"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (TableLogErrorPushLineItem) TableName() string {
	return Tables.LogErrorPushLineItem
}
