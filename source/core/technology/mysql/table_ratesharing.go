package mysql

import (
)

type TableRateSharing struct {
	Id        int64     `gorm:"column:id" json:"id"`
	UserId    int64     `gorm:"column:user_id" json:"user_id"`
	Rate      int64     `gorm:"column:rate" json:"rate"`
	Date      string    `gorm:"column:date" json:"date"`
}

func (TableRateSharing) TableName() string {
	return Tables.RateSharing
}
