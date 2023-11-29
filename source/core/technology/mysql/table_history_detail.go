package mysql

type TableHistoryDetail struct {
	Id            int64  `gorm:"column:id" json:"id"`
	HistoryId     int64  `gorm:"column:history_id" json:"history_id"`
	TableN        string `gorm:"column:table_name" json:"table_name"`
	FieldName     string `gorm:"column:field_name" json:"field_name"`
	FieldValueOld string `gorm:"column:field_value_old" json:"field_value_old"`
	FieldValueNew string `gorm:"column:field_value_new" json:"field_value_new"`
}

func (TableHistoryDetail) TableName() string {
	return Tables.HistoryDetail
}
