package model

import (
	"gorm.io/gorm"
)

func (AdScheduleModel) TableName() string {
	return "ad_schedule"
}

type AdScheduleModel struct {
	ID             int64                   `gorm:"column:id;primaryKey"`
	UserID         int64                   `gorm:"column:user_id;index"`
	Name           string                  `gorm:"column:name;size:500"`
	ClientType     AdScheduleClientTYPE    `gorm:"column:client_type;size:50"`
	AdBreakType    AdScheduleAdBreakTYPE   `gorm:"column:ad_break_type;size:50"`
	VpaidMode      AdScheduleVpaidModeTYPE `gorm:"column:vpaid_mode;size:50"`
	AdBreakConfigs []AdScheduleConfigModel `gorm:"foreignKey:AdScheduleID;references:ID"`
	gorm.Model
}

func (a *AdScheduleModel) Validate() (err error) {
	return
}

func (a *AdScheduleModel) IsFound() bool {
	if a.ID > 0 {
		return true
	}
	return false
}

func (AdScheduleConfigModel) TableName() string {
	return "ad_schedule_config"
}

type AdScheduleConfigModel struct {
	ID              int64                           `gorm:"column:id;primaryKey"`
	AdScheduleID    int64                           `gorm:"column:ad_schedule_id"`
	ConfigType      AdScheduleAdBreakConfigTYPE     `gorm:"column:config_type;size:50"`
	SkipSecond      int                             `gorm:"column:skip_second;type:int(11)"`
	OverlayAd       bool                            `gorm:"column:overlay_ad"`
	BreakTiming     int                             `gorm:"column:break_timing;type:int(11)"`
	BreakTimingType AdScheduleBreakTimingTYPE       `gorm:"column:break_timing_type;size:50"`
	AdTagUrls       []AdScheduleConfigAdTagUrlModel `gorm:"foreignKey:ConfigID;references:ID"`
}

func (AdScheduleConfigAdTagUrlModel) TableName() string {
	return "ad_schedule_config_ad-tag-url"
}

type AdScheduleConfigAdTagUrlModel struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	ConfigID int64  `gorm:"column:config_id"`
	AdTagUrl string `gorm:"column:ad_tag_url"`
}

type AdScheduleBreakTimingTYPE string

const (
	BreakTimingSecondIntoVideo AdScheduleBreakTimingTYPE = "seconds_into_video"
	BreakTimingTimeCode        AdScheduleBreakTimingTYPE = "time_code"
	BreakTimingPercentOfVideo  AdScheduleBreakTimingTYPE = "percent_of_video"
)

func (a AdScheduleBreakTimingTYPE) String() string {
	switch a {
	case BreakTimingSecondIntoVideo:
		return "Seconds into Video"
	case BreakTimingTimeCode:
		return "Time code"
	case BreakTimingPercentOfVideo:
		return "% of Video"
	default:
		return ""
	}
}

type AdScheduleAdBreakConfigTYPE string

const (
	ConfigTypePreroll  AdScheduleAdBreakConfigTYPE = "preroll"
	ConfigTypeMidroll  AdScheduleAdBreakConfigTYPE = "midroll"
	ConfigTypePostroll AdScheduleAdBreakConfigTYPE = "postroll"
	ConfigTypeVmap     AdScheduleAdBreakConfigTYPE = "vmap"
)

func (a AdScheduleAdBreakConfigTYPE) String() string {
	switch a {
	case ConfigTypePreroll:
		return "Pre-roll"
	case ConfigTypeMidroll:
		return "Mid-roll"
	case ConfigTypePostroll:
		return "Post-roll"
	case ConfigTypeVmap:
		return "vmap"
	default:
		return ""
	}
}

type AdScheduleClientTYPE string

const (
	ClientTypeVast      AdScheduleClientTYPE = "vast"
	ClientTypeGoogleIma AdScheduleClientTYPE = "google_ima"
)

func (a AdScheduleClientTYPE) String() string {
	switch a {
	case ClientTypeVast:
		return "Vast"
	case ClientTypeGoogleIma:
		return "Google IMA"
	default:
		return ""
	}
}

type AdScheduleAdBreakTYPE string

const (
	AdBreakManually AdScheduleAdBreakTYPE = "manually"
	AdBreakVmap     AdScheduleAdBreakTYPE = "vmap"
)

func (a AdScheduleAdBreakTYPE) String() string {
	switch a {
	case AdBreakManually:
		return "Manually"
	case AdBreakVmap:
		return "VMAP"
	default:
		return ""
	}
}

type AdScheduleVpaidModeTYPE string

const (
	VpaidModeInsecure AdScheduleVpaidModeTYPE = "insecure"
	VpaidModeSecure   AdScheduleVpaidModeTYPE = "secure"
	VpaidModeDisabled AdScheduleVpaidModeTYPE = "disabled"
	VpaidModeNull     AdScheduleVpaidModeTYPE = ""
)

func (a AdScheduleVpaidModeTYPE) String() string {
	switch a {
	case VpaidModeInsecure:
		return "Insecure"
	case VpaidModeSecure:
		return "Secure"
	case VpaidModeDisabled:
		return "Disabled"
	default:
		return ""
	}
}
