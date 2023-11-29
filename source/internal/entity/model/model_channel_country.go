package model

import (
	"database/sql"
)

func (ChannelCountryModel) TableName() string {
	return "channel_country"
}

type ChannelCountryModel struct {
	ID         int64        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Channel    string       `gorm:"column:channel" json:"channel"`
	CampaignID int64        `gorm:"column:campaign_id" json:"campaign_id"`
	Country    string       `gorm:"column:country" json:"country"`
	TimeOff    sql.NullTime `gorm:"column:time_off" json:"time_off"`
}
