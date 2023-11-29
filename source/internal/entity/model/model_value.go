package model

import (
	"gorm.io/gorm"
)

func (ValueModel) TableName() string {
	return "custom_value"
}

type ValueModel struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	KeyID     int64          `gorm:"column:key_id"`
	Value     string         `gorm:"column:value"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
