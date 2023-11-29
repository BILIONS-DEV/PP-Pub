package model

func (CountryModel) TableName() string {
	return "country"
}

type CountryModel struct {
	Id         int64  `gorm:"column:id" json:"id"`
	Name       string `gorm:"column:name" json:"name"`
	Code2      string `gorm:"column:code2" json:"code2"`
	Code3      string `gorm:"column:code3" json:"code3"`
	CodeNumber int    `gorm:"column:code_number" json:"code_number"`
}

func (CountryOutbrainModel) TableName() string {
	return "country_outbrain"
}

type CountryOutbrainModel struct {
	ID        int64  `gorm:"column:id" json:"id"`
	CountryID string `gorm:"column:country_id" json:"country_id"`
	Country   string `gorm:"column:country" json:"country"`
	Code2     string `gorm:"column:code2" json:"code2"`
	Code3     string `gorm:"column:code3" json:"code3"`
}
