package mysql

import (
	"time"
)

type TableNotification struct {
	Id        int64                  `gorm:"column:id" json:"id"`
	UserId    int64                  `gorm:"column:user_id" json:"user_id"`
	Status    TypeStatusNotification `gorm:"status" json:"status"`
	Message   string                 `gorm:"column:message" json:"message"`
	Action    string                 `gorm:"column:action" json:"action"`
	Link      string                 `gorm:"column:link" json:"link"`
	CreatedAt time.Time              `gorm:"column:created_at" json:"created_at"`
}

type TypeStatusNotification int

const (
	TypeStatusNotificationNew = iota + 1
	TypeStatusNotificationWatched
)
