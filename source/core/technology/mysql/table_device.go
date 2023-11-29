package mysql

type TableDevice struct {
	Id     int64  `gorm:"column:id" json:"id"`
	IdName string `gorm:"column:id_name" json:"id_name"`
	Name   string `gorm:"column:name" json:"name"`
}

func (TableDevice) TableName() string {
	return Tables.Device
}
