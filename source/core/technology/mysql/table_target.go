package mysql

type TableTarget struct {
	Id          int64
	UserId      int64 `gorm:"column:user_id" json:"user_id"`
	FloorId     int64 `gorm:"column:floor_id" json:"floor_id"`
	IdentityId  int64 `gorm:"column:identity_id" json:"identity_id"`
	LineItemId  int64 `gorm:"column:line_item_id" json:"line_item_id"`
	AbTestingId int64 `gorm:"column:ab_testing_id" json:"ab_testing_id"`
	TagId       int64 `gorm:"column:tag_id" json:"tag_id"`
	AdSizeId    int64 `gorm:"column:ad_size_id" json:"ad_size_id"`
	InventoryId int64 `gorm:"column:inventory_id" json:"inventory_id"`
	AdFormatId  int64 `gorm:"column:ad_format_id" json:"ad_format_id"`
	GeoId       int64 `gorm:"column:geo_id" json:"geo_id"`
	DeviceId    int64 `gorm:"column:device_id" json:"device_id"`
}

func (TableTarget) TableName() string {
	return Tables.Target
}
