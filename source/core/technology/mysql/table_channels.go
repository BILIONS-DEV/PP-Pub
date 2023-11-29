package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableChannels struct {
	Id          int64                  `gorm:"column:id" json:"id"`
	UserId      int64                  `gorm:"column:user_id" json:"user_id"`
	Status      TYPEStatus             `gorm:"column:status" json:"status"`
	Name        string                 `gorm:"column:name" json:"name"`
	Description string                 `gorm:"column:description" json:"description"`
	Category    int64                  `gorm:"column:category" json:"category"`
	Language    int64                  `gorm:"column:language" json:"language"`
	CreatedAt   time.Time              `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time              `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `gorm:"column:deleted_at" json:"deleted_at"`
	Keywords    []TableChannelsKeyword `gorm:"-"`
}

func (TableChannels) TableName() string {
	return Tables.Channels
}

func (rec *TableChannels) GetRls() {
	var keywords []TableChannelsKeyword
	Client.Where("channels_id = ?", rec.Id).Find(&keywords)
	rec.Keywords = keywords
}
