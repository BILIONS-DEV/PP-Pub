package model

import "source/core/technology/mysql"

type RlBlockingInventory struct{}

type RlsBlockingInventoryRecord struct {
	mysql.TableRlsBlockingInventory
}

func (RlsBlockingInventoryRecord) TableName() string {
	return mysql.Tables.RlBlockingInventory
}

func (t *RlBlockingInventory) CreateRl(blockingId int64, inventoryId int64) (err error) {
	recordRlBlockingInventory := RlsBlockingInventoryRecord{mysql.TableRlsBlockingInventory{
		BlockingId:  blockingId,
		InventoryId: inventoryId,
	}}
	err = mysql.Client.FirstOrCreate(&recordRlBlockingInventory, RlsBlockingInventoryRecord{mysql.TableRlsBlockingInventory{
		BlockingId:  blockingId,
		InventoryId: inventoryId,
	}}).Error
	return
}

func (t *RlBlockingInventory) GetByBlockingId(blockingId int64) (records []RlsBlockingInventoryRecord, err error) {
	err = mysql.Client.Where("blocking_id = ?", blockingId).Find(&records).Error
	return
}

func (t *RlBlockingInventory) DeleteByBlockingId(blockingId int64) (err error) {
	err = mysql.Client.Where("blocking_id = ?", blockingId).Delete(&RlsBlockingInventoryRecord{}).Error
	return
}
