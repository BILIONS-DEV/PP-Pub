package model

type KeyValueGamModel struct {
	ID          int64  `gorm:"column:id" json:"id"`
	NetworkID   int64  `gorm:"column:network_id" json:"network_id"`
	AbTestID    int64  `gorm:"column:ab_test_id" json:"ab_test_id"`
	InventoryID int64  `gorm:"column:inventory_id" json:"inventory_id"`
	TagID       int64  `gorm:"column:tag_id" json:"tag_id"`
	BidderID    int64  `gorm:"column:bidder_id" json:"bidder_id"`
	Name        string `gorm:"column:name" json:"name"`
	Value       string `gorm:"column:value" json:"value"`
	KeyID       int64  `gorm:"column:key_id" json:"key_id"`
	ValueID     int64  `gorm:"column:value_id" json:"value_id"`
}

func (KeyValueGamModel) TableName() string {
	return "key_value_gam"
}

func (a *KeyValueGamModel) IsFound() bool {
	if a.ID > 0 {
		return true
	}
	return false
}
