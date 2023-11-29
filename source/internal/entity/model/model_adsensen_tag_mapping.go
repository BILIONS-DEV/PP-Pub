package model

import "time"

func (AdsenseTagMappingModel) TableName() string {
	return "adsense_tag_mapping"
}

type AdsenseTagMappingModel struct {
	ID            int64                       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AdsenseAdUnit int64                       `gorm:"column:adsense_adunit" json:"adsense_adunit"`
	AdUnitType    TYPEAdUnitType              `gorm:"column:adunit_type" json:"adunit_type"`
	TagID         int64                       `gorm:"column:tag_id" json:"tag_id"`
	BidderID      int64                       `gorm:"column:bidder_id" json:"bidder_id"`
	Status        TYPEStatusAdsenseTagMapping `gorm:"column:status" json:"status"`
	Render        TYPERenderAdsenseTagMapping `gorm:"column:render" json:"render"`
	CreatedAt     time.Time                   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time                   `gorm:"column:updated_at" json:"updated_at"`
}

type TYPEAdUnitType string

const (
	TYPEAdUnitTypeFixed      TYPEAdUnitType = "fixed"
	TYPEAdUnitTypeResponsive TYPEAdUnitType = "responsive"
)

type TYPEStatusAdsenseTagMapping string

const (
	TYPEStatusAdsenseTagMappingOn  TYPEStatusAdsenseTagMapping = "on"
	TYPEStatusAdsenseTagMappingOff TYPEStatusAdsenseTagMapping = "off"
)

type TYPERenderAdsenseTagMapping string

const (
	TYPERenderAdsenseTagMappingDirect   TYPERenderAdsenseTagMapping = "direct"
	TYPERenderAdsenseTagMappingPassBack TYPERenderAdsenseTagMapping = "passback"
)

func (a *AdsenseTagMappingModel) IsFound() bool {
	if a.ID > 0 {
		return true
	}
	return false
}
