package mysql

func (TableSetting) TableName() string {
	return Tables.Setting
}

type TableSetting struct {
	Id        int64  `gorm:"id" json:"id"`
	MetaKey   string `gorm:"column:meta_key" json:"meta_key"`
	MetaValue string `gorm:"column:meta_value" json:"meta_value"`
}
