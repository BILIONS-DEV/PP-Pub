package mysql

import (
	"time"
)

type TableRevenueShare struct {
	Id     int64     `gorm:"column:id" json:"id"`
	UserId int64     `gorm:"column:user_id" json:"user_id"`
	Rate   int64     `gorm:"column:rate" json:"rate"`
	Date   time.Time `gorm:"column:date" json:"date"`
}

func (TableRevenueShare) TableName() string {
	return Tables.RevenueShare
}
