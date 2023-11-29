package mysql

import (
	"gorm.io/gorm"
	"time"
)

type LogErrorWorker struct {
	Id        int64          `gorm:"column:id" json:"id"`
	Function  string         `gorm:"column:function" json:"function"`
	Path      string         `gorm:"column:" json:"path"`
	Line      string         `gorm:"column:" json:"line"`
	Message   string         `gorm:"column:" json:"message"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
