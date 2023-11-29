package mysql

import (
	"time"
)

type TableInventoryAd struct {
	Id          int64     `gorm:"column:id" json:"id"`
	InventoryId int64     `gorm:"column:inventory_id" json:"inventory_id"`
	UrlType     string    `gorm:"column:url_type" json:"url_type"`
	Url         string    `gorm:"column:url" json:"url"`
	Root        string    `gorm:"column:root" json:"root"`
	Status      string    `gorm:"column:status" json:"status"`
	SortAd      int       `gorm:"column:sort_ad" json:"sort_ad"`
	Ads         string    `gorm:"column:ads" json:"ads"` // Dùng InventoryAds để parse json
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (TableInventoryAd) TableName() string {
	return Tables.InventoryAd
}

type InventoryAds struct {
	ID      int64  `form:"id" json:"id"`
	UrlType string `form:"url_type" json:"url_type"`
	Status  string `form:"status" json:"status"`
	Url     string `form:"url" json:"url"`
	Root    string `form:"root" json:"root"`
	Number  int    `form:"number" json:"number"`
	Ads     []Ad   `form:"ads" json:"ads"`
}

type Ad struct {
	AdName                string   `form:"ad_name" json:"ad_name"`
	DeviceType            string   `form:"device_type" json:"device_type"`
	Device                []string `form:"device" json:"device"`
	CssSelector           string   `form:"css_selector" json:"css_selector"`
	AdPosition            string   `form:"ad_position" json:"ad_position"`
	AdSpace               int      `form:"ad_space" json:"ad_space"`
	LimitAd               int      `form:"limit_ad" json:"limit_ad"`
	RenderFirstImpression bool     `form:"render_first_impression" json:"render_first_impression"`
	PaddingLeft           int      `form:"padding_left" json:"padding_left"`
	PaddingTop            int      `form:"padding_top" json:"padding_top"`
	PaddingRight          int      `form:"padding_right" json:"padding_right"`
	PaddingBottom         int      `form:"padding_bottom" json:"padding_bottom"`
	Align                 string   `form:"align" json:"align"`
	Sticky4k              string   `form:"sticky_4k" json:"sticky_4k"`
	Advertisement         bool     `form:"advertisement" json:"advertisement"`
	AdContent             string   `form:"ad_content" json:"ad_content"`
	Tag                   int64    `form:"tag_id" json:"tag_id"`
	CustomScript          string   `form:"custom_script" json:"custom_script"`
	Scan                  string   `form:"scan" json:"scan"`
	DeviceCustomMin       []int    `form:"device_custom_min" json:"device_custom_min"`
}
