package mysql

type TableFloorTargeting struct {
	ID          int64   `gorm:"column:id;primaryKey;autoIncrement;type:int(11)" json:"id"`
	NetworkId   int64   `gorm:"column:network_id;type:bigint(20)" json:"network_id"`
	CustomValue string  `gorm:"column:custom_value;type:varchar(100)" json:"custom_value"`
	CountryCode string  `gorm:"column:countryCode;type:varchar(50)" json:"countryCode"`
	Device      string  `gorm:"column:device;type:varchar(50)" json:"device"`
	TagId       int64   `gorm:"column:tag_id;type:int(11)" json:"tagId"`
	Floor       float64 `gorm:"column:floor;type:double" json:"floor"`
}

func (TableFloorTargeting) TableName() string {
	return Tables.FloorTargeting
}
