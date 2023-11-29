package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableBlocking struct {
	Id              int64                       `gorm:"column:id" json:"id"`
	UserId          int64                       `gorm:"column:user_id" json:"user_id"`
	RestrictionName string                      `gorm:"column:restriction_name" json:"restriction_name"`
	CreatedAt       time.Time                   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time                   `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt              `gorm:"column:deleted_at" json:"deleted_at"`
	Inventories     []TableInventory            `gorm:"-" as:"-"`
	Restrictions    []TableBlockingRestrictions `gorm:"-" as:"-"`
}

func (TableBlocking) TableName() string {
	return Tables.Blocking
}

func (rec *TableBlocking) GetRls() {
	// Get các rls
	var rlsBlockInventories []TableRlsBlockingInventory
	Client.Where("blocking_id = ?", rec.Id).Find(&rlsBlockInventories)
	var inventories []TableInventory
	for _, rlsBlockInventory := range rlsBlockInventories {
		var inventory TableInventory
		Client.Find(&inventory, rlsBlockInventory.InventoryId)
		inventories = append(inventories, inventory)
	}
	rec.Inventories = inventories

	var restrictions []TableBlockingRestrictions
	Client.Where(TableBlockingRestrictions{BlockingId: rec.Id}).Find(&restrictions)
	rec.Restrictions = restrictions
	return
}
