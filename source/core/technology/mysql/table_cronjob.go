package mysql

import (
	"time"
)

type TableCronjob struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatorID int64     `gorm:"column:creator_id" json:"creator_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Type      string    `gorm:"column:type" json:"type"`
	Data      string    `gorm:"column:data" json:"data"`
	DataMeta  string    `gorm:"column:data_meta" json:"data_meta"`
	Processed string    `gorm:"column:processed" json:"processed"`
	Status    string    `gorm:"column:status" json:"status"`
	LogError  string    `gorm:"column:log_error" json:"log_error"`
	Number    int       `gorm:"column:number" json:"number"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (TableCronjob) TableName() string {
	return Tables.Cronjob
}
