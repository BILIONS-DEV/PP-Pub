package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func (TableRule) TableName() string {
	return Tables.Rule
}

type TableRule struct {
	ID           int64             `gorm:"column:id" json:"id"`
	UserID       int64             `gorm:"column:user_id" json:"user_id"`
	Name         string            `gorm:"column:name" json:"name"`
	Status       TYPEStatusOnOff   `gorm:"column:status" json:"status"`
	Type         TYPERuleType      `gorm:"column:type" json:"type"`
	Description  string            `gorm:"column:description" json:"description"`
	CreatedAt    time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time         `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `gorm:"column:deleted_at" json:"deleted_at"`
	BlockedPages []RuleBlockedPage `gorm:"foreignKey:RuleID;references:ID"`
}

func (rec *TableRule) GetById(id, userId int64) (err error) {
	err = Client.Preload(clause.Associations).Preload("BlockedPages").Where("user_id = ?", userId).Find(&rec, id).Error
	return
}

func (rec *TableRule) GetRls() {
	var blockedPages []RuleBlockedPage
	Client.Debug().Where("rule_id = ?", rec.ID).Find(&blockedPages)
	rec.BlockedPages = blockedPages
}

type TYPEStatusOnOff int

const (
	TYPEStatusOn TYPEStatusOnOff = iota + 1
	TYPEStatusOff
)

func (t TYPEStatusOnOff) String() string {
	switch t {
	case TYPEStatusOn:
		return "on"
	case TYPEStatusOff:
		return "off"
	default:
		return ""
	}
}

type TYPERuleType int

const (
	TYPERuleTypeFloor TYPERuleType = iota + 1
	TYPERuleTypeBlocking
	TYPERuleTypeBlockedPage
)

func (t TYPERuleType) String() string {
	switch t {
	case TYPERuleTypeFloor:
		return "floor"
	case TYPERuleTypeBlocking:
		return "blocking"
	case TYPERuleTypeBlockedPage:
		return "blocked-page"
	default:
		return ""
	}
}

func (t *TableRule) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}
