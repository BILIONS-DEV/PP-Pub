package mysql

type TableKeyValueGam struct {
	Id          int64  `gorm:"column:id" json:"id"`
	NetworkId   int64  `gorm:"column:network_id" json:"network_id"`
	AbTestId    int64  `gorm:"column:ab_test_id" json:"ab_test_id"`
	InventoryId int64  `gorm:"column:inventory_id" json:"inventory_id"`
	TagId       int64  `gorm:"column:tag_id" json:"tag_id"`
	BidderId    int64  `gorm:"column:bidder_id" json:"bidder_id"`
	Name        string `gorm:"column:name" json:"name"`
	Value       string `gorm:"column:value" json:"value"`
	KeyId       int64  `gorm:"column:key_id" json:"key_id"`
	ValueId     int64  `gorm:"column:value_id" json:"value_id"`
}

func (TableKeyValueGam) TableName() string {
	return Tables.KeyValueGam
}
