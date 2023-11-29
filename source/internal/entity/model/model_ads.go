package model

type AdsModel struct {
	ID          int64    `gorm:"column:id" json:"id"`
	UserID      int64    `gorm:"column:user_id" json:"user_id"`
	CampaignID  int64    `gorm:"column:campaign_id" json:"campaign_id"`
	AdGroupID   int64    `gorm:"column:ad_group_id" json:"ad_group_id"`
	AdType      string   `gorm:"column:type" json:"type"`
	Size        string   `gorm:"column:size" json:"size"`
	SiteName    string   `gorm:"default:;column:site_name" json:"site_name"`
	Headline    string   `gorm:"column:headline" json:"headline"`
	Description string         `gorm:"column:description" json:"description"`
	Media       string   `gorm:"column:media" json:"media"`
	MediaURL    string   `json:"media_url"`
	MediaType   string   `gorm:"column:media_type" json:"media_type"`
	StatusAd    bool     `gorm:"default:pending;column:status_ad" json:"status_ad"`
	Thumbnail   string   `gorm:"default:;column:thumbnail" json:"thumbnail"`
	ClickURL    string   `gorm:"column:click_url" json:"click_url"`
	CTA         string   `gorm:"default:;column:call_to_action" json:"call_to_action"`
	Action      string   `gorm:"default:on;column:action" json:"action"`
	AdErrorID   string   `gorm:"default:;column:ad_error_id" json:"ad_error_id"`
	Audience    Audience `json:"audience,omitempty"`
	Impressions int64    `json:"impressions"`
	Revenue     float64  `json:"revenue"`
}

type Audience struct {
	ID        int64          `gorm:"column:id" json:"id"`
	UserID    int64          `gorm:"column:user_id" json:"user_id"`
	Name      string         `gorm:"column:name" json:"name"`
	Inventory AudienceTarget `gorm:"column:inventory" json:"inventory"`
	Category  AudienceTarget `gorm:"column:category" json:"category"`
	Device    AudienceTarget `gorm:"column:device" json:"device"`
	Locations AudienceTarget `gorm:"column:locations" json:"locations"`
	Languages AudienceTarget `gorm:"column:languages" json:"languages"`
}

type AudienceTarget struct {
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}
