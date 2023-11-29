package mysql2

type TableTarget struct {
	ID          int64 `gorm:"column:id"`
	UserID      int64 `gorm:"column:user_id"`
	LineItemID  int64 `gorm:"column:line_item_id"`
	FloorID     int64 `gorm:"column:floor_id"`
	IdentityID  int64 `gorm:"column:identity_id"`
	AbTestingID int64 `gorm:"column:ab_testing_id"`
	InventoryID int64 `gorm:"column:inventory_id"`
	TagID       int64 `gorm:"column:tag_id"`
	SizeID      int   `gorm:"column:ad_size_id"`
	FormatID    int   `gorm:"column:ad_format_id"`
	GeoID       int   `gorm:"column:geo_id"`
	DeviceID    int   `gorm:"column:device_id"`
}

func (TableTarget) TableName() string {
	return Tables.Target
}
