package model

import "gorm.io/gorm"

type AdTypeModel struct {
	Id          int64           `gorm:"column:id" json:"id"`
	IdName      string          `gorm:"column:id_name" json:"id_name"`
	Name        string          `gorm:"column:name" json:"name"`
	DisplayName string          `gorm:"column:display_name" json:"display_name"`
	Type        TypeBannerVideo `gorm:"column:type" json:"type"`
	DeletedAt   gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`
}

func (AdTypeModel) TableName() string {
	return "ad_type"
}

type TypeBannerVideo int

const (
	TypeBanner TypeBannerVideo = iota + 1
	TypeVideo
)

func (t TypeBannerVideo) String() string {
	switch t {
	case TypeBanner:
		return "Banner"
	case TypeVideo:
		return "Video"
	default:
		return ""
	}
}

type AdType int

const (
	AdTypeDisplay AdType = iota + 1
	AdTypeInstream
	AdTypeOutstream
	AdTypePinZone
	AdTypeStickyBanner
	AdTypeInterstitial
	AdTypeRelatedZone
)

func (t AdType) String() string {
	switch t {
	case AdTypeDisplay:
		return "Display"
	case AdTypeInstream:
		return "Instream"
	case AdTypeOutstream:
		return "Outstream"
	case AdTypePinZone:
		return "Pin Zone"
	case AdTypeStickyBanner:
		return "Sticky Banner"
	case AdTypeInterstitial:
		return "Interstitial"
	case AdTypeRelatedZone:
		return "Related Zone"
	default:
		return ""
	}
}
