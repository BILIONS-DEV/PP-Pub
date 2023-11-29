package mysql

type TableRlsBlockingInventory struct {
	Id          int64 `gorm:"column:id" json:"id"`
	BlockingId  int64 `gorm:"column:blocking_id" json:"blocking_id"`
	InventoryId int64 `gorm:"column:inventory_id" json:"inventory_id"`
}

func (TableRlsBlockingInventory) TableName() string {
	return Tables.RlBlockingInventory
}