package mysql

type TableCurrency struct {
	Id                    int64  `gorm:"column:id" json:"id"`
	Name                  string `gorm:"column:name" json:"name"`
	Code                  string `gorm:"column:code" json:"code"`
	GranularityMultiplier int    `gorm:"column:granularity_multiplier" json:"granularity_multiplier"`
	Symbol                string `gorm:"column:symbol" json:"symbol"`
}

func (TableCurrency) TableName() string {
	return Tables.Currency
}
