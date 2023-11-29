package model

import (
	"time"
)

func (AdsenseChannelModel) TableName() string {
	return "adsense_channel"
}

type AdsenseChannelModel struct {
	ID      int64  `gorm:"column:id;primary_key" json:"id"`
	Channel string `gorm:"column:channel" json:"channel"`
	Account string `gorm:"column:account" json:"account"`
	// CampaignID int64     `gorm:"column:campaign_id" json:"campaign_id"`
	TimeOff time.Time `gorm:"column:time_off" json:"time_off"`
}
