package mysql

import (
	"gorm.io/gorm"
)

type TableInventoryConfig struct {
	Id                  int64                 `gorm:"column:id" json:"id"`
	InventoryId         int64                 `gorm:"column:inventory_id" json:"inventory_id"`
	GamAccount          int64                 `gorm:"column:gam_account" json:"gam_account"`
	AdRefresh           TypeOnOff             `gorm:"column:ad_refresh" json:"ad_refresh"`
	DirectSales         TypeOnOff             `gorm:"column:direct_sales" json:"direct_sales"`
	AdRefreshTime       int                   `gorm:"column:ad_refresh_time" json:"ad_refresh_time"`
	PrebidTimeOut       int                   `gorm:"column:prebid_time_out" json:"prebid_time_out"`
	LoadAdType          string                `gorm:"column:load_ad_type" json:"load_ad_type"`
	LoadAdRefresh       string                `gorm:"column:load_ad_refresh" json:"load_ad_refresh"`
	Gdpr                int                   `gorm:"column:gdpr" json:"gdpr"`
	GdprTimeout         int                   `gorm:"column:gdpr_timeout" json:"gdpr_timeout"`
	Ccpa                int                   `gorm:"column:ccpa" json:"ccpa"`
	CcpaTimeout         int                   `gorm:"column:ccpa_timeout" json:"ccpa_timeout"`
	CustomBrand         int                   `gorm:"column:custom_brand" json:"custom_brand"`
	SafeFrame           TypeOnOff             `gorm:"column:safe_frame" json:"safe_frame"`
	Logo                string                `gorm:"column:logo" json:"logo"`
	Title               string                `gorm:"column:title" json:"title"`
	Content             string                `gorm:"column:content" json:"content"`
	AuctionDelay        int                   `gorm:"column:auction_delay" json:"auction_delay"`
	SyncDelay           int                   `gorm:"column:sync_delay" json:"sync_delay"`
	GamAutoCreate       TypeOnOff             `gorm:"column:gam_auto_create" json:"gam_auto_create"`
	PbRenderMode        TYPEPbRenderMode      `gorm:"column:pb_render_mode" json:"pb_render_mode"`
	FetchMarginPercent  int                   `gorm:"column:fetch_margin_percent" json:"fetch_margin_percent"`
	RenderMarginPercent int                   `gorm:"column:render_margin_percent" json:"render_margin_percent"`
	MobileScaling       int                   `gorm:"column:mobile_scaling" json:"mobile_scaling"`
	CustomGoogleShare   *string               `gorm:"default:null;column:custom_google_share" json:"custom_google_share"`
	GoogleRequestMode   TYPEGoogleRequestMode `gorm:"default:always;column:google_request_mode" json:"google_request_mode"`
	DeletedAt           gorm.DeletedAt        `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableInventoryConfig) TableName() string {
	return Tables.InventoryConfig
}

type TYPEPbRenderMode int

const (
	TYPEPbRenderModeInIframe = iota + 1
	TYPEPbRenderModeInRootDocument
)

func (t TYPEPbRenderMode) String() string {
	switch t {
	case TYPEPbRenderModeInIframe:
		return "iframe"
	case TYPEPbRenderModeInRootDocument:
		return "root"
	default:
		return ""
	}
}

type TYPEGoogleRequestMode string

const (
	TYPEGoogleRequestModeAlways TYPEGoogleRequestMode = "always"
	TYPEGoogleRequestModeOnly   TYPEGoogleRequestMode = "only"
	TYPEGoogleRequestModeInView TYPEGoogleRequestMode = "inview"
)
