package mysql2

type TableBlockingRestriction struct {
	BlockingID          int64  `gorm:"column:blocking_id;"`
	InventoryID         int64  `gorm:"column:inventory_id"`
	InventoryAdvertiser string `gorm:"column:inventory_advertiser"`
	CreativeID          string `gorm:"column:creative_id"`
}

func (TableBlockingRestriction) TableName() string {
	return Tables.BlockingRestriction
}
