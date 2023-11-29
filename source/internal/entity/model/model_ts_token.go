package model

import (
	"time"
)

func (TsTokenModel) TableName() string {
	return "ts_token"
}

type TsTokenModel struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TrafficSource string    `gorm:"column:traffic_source" json:"traffic_source"`
	Token         string    `gorm:"column:token" json:"token"`
	Date          time.Time `gorm:"column:date" json:"date"`
}
