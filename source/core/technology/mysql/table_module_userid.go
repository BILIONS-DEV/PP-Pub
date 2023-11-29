package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableModuleUserId struct {
	Id                   int64  `gorm:"column:id" json:"id"`
	UserId               int64  `gorm:"column:user_id" json:"user_id"`
	Name                 string `gorm:"column:name" json:"name"`
	PrebidModuleFilename string `gorm:"column:prebid_module_filename" json:"prebid_module_filename"`
	Params               string `gorm:"params" json:"params" form:"params"`
	Storage              string `gorm:"storage" json:"storage" form:"storage"`
	//Status    TYPEStatus     `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableModuleUserId) TableName() string {
	return Tables.ModuleUserId
}

type ParamModuleUserId struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template"`
}

type StorageModuleUserId struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template"`
}
