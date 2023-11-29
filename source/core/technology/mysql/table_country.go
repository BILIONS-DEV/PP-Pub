package mysql

type TableCountry struct {
	Id         int64  `gorm:"column:id" json:"id"`
	Name       string `gorm:"column:name" json:"name"`
	Code2      string `gorm:"column:code2" json:"code2"`
	Code3      string `gorm:"column:code3" json:"code3"`
	CodeNumber int    `gorm:"column:code_number" json:"code_number"`
}

func (TableCountry) TableName() string {
	return Tables.Country
}
