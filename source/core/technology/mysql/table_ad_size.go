package mysql

type TableAdSize struct {
	Id     int64  `gorm:"column:id" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Width  int    `gorm:"column:width" json:"width"`
	Height int    `gorm:"column:height" json:"height"`
	Type   int    `gorm:"column:type" json:"type"`
}

func (TableAdSize) TableName() string {
	return Tables.AdSize
}
