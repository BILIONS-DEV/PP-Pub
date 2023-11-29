package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

func (KeyModel) TableName() string {
	return "custom_key"
}

type KeyModel struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64          `gorm:"column:user_id"`
	KeyName   string         `gorm:"column:key_name"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	Value     []ValueModel   `gorm:"foreignKey:KeyID;references:ID"`
}

func (t *KeyModel) Validate() (err error) {
	return
}

func (t *KeyModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}

func (t *KeyModel) ToJSON() string {
	jsonEncode, _ := json.Marshal(t)
	return string(jsonEncode)
}
