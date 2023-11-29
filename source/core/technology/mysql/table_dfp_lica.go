package mysql

import "time"

type TableDfpLica struct {
	Id         int64
	UserId     int64     `gorm:"column:user_id" json:"user_id"`
	LineItemId string    `gorm:"column:line_item_id" json:"line_item_id"`
	CreativeId string    `gorm:"column:creative_id" json:"creative_id"`
	Size       string    `gorm:"column:size" json:"size"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
