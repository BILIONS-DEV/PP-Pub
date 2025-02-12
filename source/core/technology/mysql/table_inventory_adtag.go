package mysql

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type TableInventoryAdTag struct {
	Id                  int64              `gorm:"column:id" json:"id"`
	Name                string             `gorm:"column:name" json:"name"`
	UserId              int64              `gorm:"column:user_id" json:"user_id"`
	InventoryId         int64              `gorm:"column:inventory_id" json:"inventory_id"`
	Gam                 string             `gorm:"column:gam" json:"gam"`
	GamAuto             string             `gorm:"column:gam_auto" json:"gam_auto"`
	Type                TYPEAdType         `gorm:"column:type" json:"type"`
	PassBack            string             `gorm:"column:pass_back" json:"pass_back"`
	PassBackType        TYPEPassBackType   `gorm:"column:passback_type" json:"passback_type"`
	InlineTag           int64              `gorm:"column:inline_tag" json:"inline_tag"`
	Renderer            TYPERenderer       `gorm:"column:renderer" json:"renderer"`
	Status              TypeStatusAdTag    `gorm:"column:status" json:"status"`
	AdSize              TYPEAdSize         `gorm:"column:ad_size" json:"ad_size"`
	ResponsiveType      TYPEResponsiveType `gorm:"column:responsive_type;default:null" json:"responsive_type"`
	PrimaryAdSize       int64              `gorm:"column:primary_ad_size;default:94" json:"primary_ad_size"`
	PrimaryAdSizeMobile int64              `gorm:"column:primary_ad_size_mobile" json:"primary_ad_size_mobile"`
	SizeOnMobile        int64              `gorm:"column:size_on_mobile" json:"size_on_mobile"`
	PassBackMobile      string             `gorm:"column:pass_back_mobile" json:"pass_back_mobile"`
	// AdditionalAdSize  string             `gorm:"column:additional_ad_size" json:"additional_ad_size"`
	TemplateId     int64              `gorm:"column:template_id" json:"template_id"`
	ContentSource  TypeContentSource  `gorm:"column:content_source" json:"content_source"`
	PlaylistId     int64              `gorm:"column:playlist_id" json:"playlist_id"`
	FeedUrl        string             `gorm:"column:feed_url" json:"feed_url"`
	RelatedContent TYPERelatedContent `gorm:"column:related_content" json:"related_content"`
	ContentType    TYPEContentType    `gorm:"column:content_type" json:"content_type"`
	Uuid           string             `gorm:"column:uuid" json:"uuid"`
	BidOutStream   TypeOnOff          `gorm:"column:bid_out_stream" json:"bid_out_stream"`
	GamSticky      string             `gorm:"column:gam_sticky" json:"gam_sticky"`
	// SizeSticky        int64              `gorm:"size_sticky" json:"size_sticky"`
	PositionSticky          TypePositionSticky       `gorm:"position_sticky" json:"position_sticky"`
	PositionStickyMobile    TypePositionSticky       `gorm:"position_sticky_mobile" json:"position_sticky_mobile"`
	CloseButtonSticky       int                      `gorm:"close_button_sticky" json:"close_button_sticky"`
	CloseButtonStickyMobile TYPEOnOff                `gorm:"close_button_sticky_mobile" json:"close_button_sticky_mobile"`
	FrequencyCaps           int                      `gorm:"frequency_caps" json:"frequency_caps"`
	Output                  string                   `gorm:"output" json:"output"`
	MainTitle               string                   `gorm:"main_title;default:null" json:"main_title"`
	BackgroundColor         string                   `gorm:"background_color;default:null" json:"background_color"`
	TitleColor              string                   `gorm:"title_color;default:null" json:"title_color"`
	ShiftContent            TYPEOnOff                `gorm:"shift_content;default:2" json:"shift_content"`
	EnableStickyDesktop     TYPEOnOff                `gorm:"column:enable_sticky_desktop" json:"enable_sticky_desktop"`
	EnableStickyMobile      TYPEOnOff                `gorm:"column:enable_sticky_mobile" json:"enable_sticky_mobile"`
	BannerAD                TYPEOnOff                `gorm:"column:banner_ad" json:"banner_ad"`
	VideoAD                 TYPEOnOff                `gorm:"column:video_ad" json:"video_ad"`
	SyncPocPoc              string                   `gorm:"default:pending;column:sync_pocpoc" json:"sync_pocpoc"`
	AdRefresh               TYPEAdRefresh            `gorm:"default:domain_configuration;column:ad_refresh" json:"ad_refresh"`
	AdRefreshTime           int64                    `gorm:"default:0;column:ad_refresh_time" json:"ad_refresh_time"`
	Sticky4k                string                   `gorm:"default:off;column:sticky4k" json:"sticky4k"`
	Version                 int64                    `gorm:"default:0;column:version" json:"version"`
	CreatedAt               time.Time                `gorm:"column:created_at" json:"created_at"`
	UpdatedAt               time.Time                `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt               gorm.DeletedAt           `gorm:"column:deleted_at" json:"deleted_at"`
	Template                TableTemplate            `gorm:"-"`
	Playlist                TablePlaylist            `gorm:"-"`
	AdditionalAdSize        []TableAdSize            `gorm:"-"`
	AdditionalAdSizeMobile  []TableAdSize            `gorm:"-"`
	RewardedAdConfig        TableTagRewardedAdConfig `gorm:"-"`
}

func (TableInventoryAdTag) TableName() string {
	return Tables.InventoryAdTag
}

func (rec *TableInventoryAdTag) GetFullData() {
	if !rec.IsFound() {
		return
	}
	rec.GetTemplate()
	rec.GetPlaylist()
	rec.GetSizeAdditional()
	//rec.GetAdsenseTagModel()
	rec.GetRewardedAdConfig()
}

func (rec *TableInventoryAdTag) GetTemplate() {
	if !rec.IsFound() || rec.TemplateId == 0 {
		return
	}
	var template TableTemplate
	Client.Find(&template, rec.TemplateId)
	rec.Template = template
}

func (rec *TableInventoryAdTag) GetPlaylist() {
	if !rec.IsFound() || rec.PlaylistId == 0 {
		return
	}
	var playlist TablePlaylist
	Client.Find(&playlist, rec.PlaylistId)
	rec.Playlist = playlist
}

func (rec *TableInventoryAdTag) GetSizeAdditional() {
	if !rec.IsFound() {
		return
	}
	var rlsAdTagSizeAdditionals []TableRlsAdTagSizeAdditional
	Client.Where(TableRlsAdTagSizeAdditional{AdTagId: rec.Id, Device: 1}).Find(&rlsAdTagSizeAdditionals)
	var additionalAdSize []TableAdSize
	for _, rlsAdTagSizeAdditional := range rlsAdTagSizeAdditionals {
		var sizeAdditional TableAdSize
		Client.Find(&sizeAdditional, rlsAdTagSizeAdditional.AdSizeId)
		additionalAdSize = append(additionalAdSize, sizeAdditional)
	}
	rec.AdditionalAdSize = additionalAdSize

	var rlsAdTagSizeAdditionalsMobile []TableRlsAdTagSizeAdditional
	Client.Where(TableRlsAdTagSizeAdditional{AdTagId: rec.Id, Device: 2}).Find(&rlsAdTagSizeAdditionalsMobile)
	var additionalAdSizeMobile []TableAdSize
	for _, rlsAdTagSizeAdditional := range rlsAdTagSizeAdditionalsMobile {
		var sizeAdditional TableAdSize
		Client.Find(&sizeAdditional, rlsAdTagSizeAdditional.AdSizeId)
		additionalAdSizeMobile = append(additionalAdSizeMobile, sizeAdditional)
	}
	rec.AdditionalAdSizeMobile = additionalAdSizeMobile
}

func (rec *TableInventoryAdTag) GetRewardedAdConfig() {
	if !rec.IsFound() || rec.Type != TYPERewardedAd {
		return
	}
	var cf TableTagRewardedAdConfig
	Client.Find(&cf, rec.Id)
	rec.RewardedAdConfig = cf
}

func (rec *TableInventoryAdTag) IsFound() bool {
	if rec.Id > 0 {
		return true
	}
	return false
}

func (rec *TableInventoryAdTag) GetAdTag(isSticky4k string) string {
	var userRecord TableUser
	Client.Find(&userRecord, rec.UserId)

	tag := "pw_" + strconv.FormatInt(rec.Id, 10)
	//var class, bTag, bAPI, divTag, tagDomain string
	var class, bTag, bAPI, divTag string
	if userRecord.ParentSub == "yes" {
		class = "futureads"
		bTag = "wapTag"
		bAPI = "wAPITag"
		divTag = "div"
		//tagDomain = "cdn.bilsyndication.com"
	} else {
		if userRecord.SystemSync == 1 {
			class = "adsbyvli"
			bTag = "vitag"
			bAPI = "viAPItag"
			divTag = "div"
			//tagDomain = "cdn.vlitag.com"
		} else {
			class = "adsbypubpower"
			bTag = "powerTag"
			bAPI = "powerAPITag"
			divTag = "pubtag"
			//tagDomain = "nc.pubpowerplatform.io"
		}
	}

	switch rec.Type {
	case TYPEDisplay:
		var style string
		if rec.AdSize == TYPEAdSizeFixed {
			var size TableAdSize
			Client.Find(&size, rec.PrimaryAdSize)
			style = "style='width:" + strconv.Itoa(size.Width) + "px;height:" + strconv.Itoa(size.Height) + "px;display:inline-block;'"
		}
		var sticky4k string
		if isSticky4k == "on" {
			sticky4k = "data-ad-mode='sticky'"
		} else {
			sticky4k = ""
		}
		return fmt.Sprintf("<%s class='%s' %s data-ad-slot='%s' %s></%s> <script type='text/javascript'>(%s.Init = window.%s.Init || []).push(function () { %s.display('%s') })</script>",
			divTag, class, style, tag, sticky4k, divTag, bTag, bTag, bAPI, tag,
		)
	case TYPEInStream:
		var strTag string
		if rec.Renderer != TYPERendererPubPower && rec.Renderer != TYPERendererOverlayAd && rec.Status == TypeStatusAdTagLive {
			//strTag = "https://nc.pubpowerplatform.io/vpaid/w/" + copyTag.Inventory.Uuid + "?tagid=" + strconv.FormatInt(tag.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize
			//strTag = = "https://" + tagDomain + "/vpaid/w/" + uuid + "?tagid=" + strconv.FormatInt(rec.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize
		} else {
			if rec.Renderer == TYPERendererOverlayAd {
				strTag = "<button onclick='loadAd();'>Load AD</button> <" + divTag + " class='" + class + "' data-ad-slot='" + tag + "'></" + divTag + "> <script> var " + bTag + " = " + bTag + " || {}; " + bTag + ".instreamConfig = " + bTag + ".instreamConfig || {}; " + bTag + ".instreamConfig = {loadingText: 'Loading advertisement..',started: function(){}, complete: function() {}, error: function() {}, hidden: function() {} }; function loadAd() { (" + bTag + ".Init = window." + bTag + ".Init || []).push(function () { " + bAPI + ".initPowerInstream('" + tag + "'); }); } </script>"
			} else {
				strTag = "<" + divTag + " class='" + class + "' data-ad-slot='" + tag + "'></" + divTag + "><script type='text/javascript'> (" + bTag + ".Init = window." + bTag + ".Init || []).push(function () { " + bAPI + ".initPowerInstream('" + tag + "'); }); </script>"
			}
		}
		return strTag
	case TYPEOutStream:
		var strTag string
		if rec.Renderer != TYPERendererPubPower && rec.Status == TypeStatusAdTagLive {
			//strTag = "https://nc.pubpowerplatform.io/vpaid/w/" + copyTag.Inventory.Uuid + "?tagid=" + strconv.FormatInt(tag.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize
			//strTag = "https://" + tagDomain + "/vpaid/w/" + uuid + "?tagid=" + strconv.FormatInt(rec.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize
		} else {
			strTag = "<" + divTag + " class='" + class + "' data-ad-slot='" + tag + "'></" + divTag + "><script type='text/javascript'> (" + bTag + ".Init = window." + bTag + ".Init || []).push(function () { " + bAPI + ".initPowerOutstream('pw_{{.AdTag.TableInventoryAdTag.Id}}'); }); </script>"
		}
		return strTag
	case TYPETopArticles:
		return ""
	case TYPEStickyBanner:
		return "<script type='text/javascript'> (" + bTag + ".Init = window." + bTag + ".Init || []).push(function () { " + bAPI + ".initStickyBanner('" + tag + "'); }); </script>"
	case TYPEInterstitial:
		return ""
	case TYPEPlayZone:
		return "<" + divTag + " class='" + class + "' data-ad-slot='" + tag + "'></" + divTag + "> <script type='text/javascript'> (" + bTag + ".Init = window." + bTag + ".Init || []).push(function () { " + bAPI + ".initPlayZone('" + tag + "') }) </script>"
	case TYPEPNative:
		return "<" + divTag + " class='" + class + "' data-ad-slot='" + tag + "'></" + divTag + "> <script type='text/javascript'> (" + bTag + ".Init = window." + bTag + ".Init || []).push(function () { " + bAPI + ".initNativeAds('" + tag + "') }) </script>"
	case TYPEVideo:
		var strTag string
		if rec.Renderer != TYPERendererPubPower && rec.Renderer != TYPERendererOverlayAd && rec.Status == TypeStatusAdTagLive {
			//strTag = "https://" + tagDomain + "/vpaid/w/" + uuid + "?tagid=" + strconv.FormatInt(rec.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize
		} else {
			if rec.Renderer == TYPERendererOverlayAd {
				strTag = "<button onclick='loadAd();'>Load AD</button> <" + divTag + " class='" + class + "' data-ad-slot='" + tag + "'></" + divTag + "> <script> var " + bTag + " = " + bTag + " || {}; " + bTag + ".instreamConfig = " + bTag + ".instreamConfig || {}; " + bTag + ".instreamConfig = {loadingText: 'Loading advertisement..',started: function(){}, complete: function() {}, error: function() {}, hidden: function() {} }; function loadAd() { (" + bTag + ".Init = window." + bTag + ".Init || []).push(function () { " + bAPI + ".initPowerVideoAds('" + tag + "'); }); } </script>"
			} else {
				strTag = "<" + divTag + " class='" + class + "' data-ad-slot='" + tag + "'></" + divTag + "><script type='text/javascript'> (" + bTag + ".Init = window." + bTag + ".Init || []).push(function () { " + bAPI + ".initPowerVideoAds('" + tag + "'); }); </script>"
			}
		}
		return strTag
	}
	return ""
}

func (t *TableInventoryAdTag) SaveConfigRewardedAd(rewardedAdConfig TableTagRewardedAdConfig) (err error) {
	rewardedAdConfig.TagID = t.Id
	err = Client.Debug().Save(&rewardedAdConfig).Error
	return
}

type TypeStatusAdTag int

const (
	TypeStatusAdTagLive TypeStatusAdTag = iota + 1
	TypeStatusAdTagNotLive
	TypeStatusAdTagArchived
)

func (s TypeStatusAdTag) String() string {
	switch s {
	case TypeStatusAdTagLive:
		return "Live"
	case TypeStatusAdTagNotLive:
		return "Paused"
	case TypeStatusAdTagArchived:
		return "Archive"
	default:
		return ""
	}
}

func (s TypeStatusAdTag) Int() int {
	switch s {
	case TypeStatusAdTagLive:
		return 1
	case TypeStatusAdTagNotLive:
		return 2
	case TypeStatusAdTagArchived:
		return 3
	default:
		return 0
	}
}

type TYPEAdType int

const (
	TYPEDisplay TYPEAdType = iota + 1
	TYPEInStream
	TYPEOutStream
	TYPETopArticles
	TYPEStickyBanner
	TYPEInterstitial
	TYPEPlayZone
	TYPEPNative
	TYPEVideo
	TYPERewardedAd
)

func (t TYPEAdType) String() string {
	switch t {
	case TYPEDisplay:
		return "Display"
	case TYPEInStream:
		return "Instream"
	case TYPEOutStream:
		return "Outstream"
	case TYPETopArticles:
		return "Pin Zone"
	case TYPEStickyBanner:
		return "Sticky Banner"
	case TYPEInterstitial:
		return "Interstitial"
	case TYPEPlayZone:
		return "Play Zone"
	case TYPEPNative:
		return "Native"
	case TYPEVideo:
		return "Video"
	case TYPERewardedAd:
		return "Rewarded Ad"
	default:
		return ""
	}
}

func (t TYPEAdType) StringV2() string {
	switch t {
	case TYPEDisplay:
		return "Display"
	case TYPEInStream, TYPEOutStream, TYPEVideo:
		return "Video"
	case TYPETopArticles:
		return "Pin Zone"
	case TYPEStickyBanner:
		return "Sticky Banner"
	case TYPEInterstitial:
		return "Interstitial"
	case TYPEPlayZone:
		return "Play Zone"
	case TYPEPNative:
		return "Native"
	case TYPERewardedAd:
		return "Rewarded Ad"
	default:
		return ""
	}
}

func (t TYPEAdType) Int() int64 {
	switch t {
	case TYPEDisplay:
		return 1
	case TYPEInStream:
		return 2
	case TYPEOutStream:
		return 3
	case TYPETopArticles:
		return 4
	case TYPEStickyBanner:
		return 5
	case TYPEInterstitial:
		return 6
	case TYPEPlayZone:
		return 7
	case TYPEPNative:
		return 8
	case TYPEVideo:
		return 9
	case TYPERewardedAd:
		return 10
	default:
		return 0
	}
}

func (t TYPEAdType) IsBanner() bool {
	switch t {
	case TYPEDisplay:
		return true
	case TYPEStickyBanner:
		return true
	case TYPEInterstitial:
		return true
	case TYPEPlayZone:
		return true
	default:
		return false
	}
}

func (t TYPEAdType) IsVideo() bool {
	switch t {
	case TYPEInStream:
		return true
	case TYPEOutStream:
		return true
	case TYPETopArticles:
		return true
	case TYPEVideo:
		return true
	default:
		return false
	}
}

type TypeContentSource int

const (
	TypeContentSourcePlaylist TypeContentSource = iota + 1
	TypeContentSourceFeed
	TypeContentSourceDirect
	TypeContentSourceAuto
)

func (t TypeContentSource) String() string {
	switch t {
	case TypeContentSourcePlaylist:
		return "Playlist"
	case TypeContentSourceFeed:
		return "Feed"
	case TypeContentSourceDirect:
		return "Direct"
	case TypeContentSourceAuto:
		return "Auto"
	default:
		return ""
	}
}

type TypePositionSticky int

func (t *TypePositionSticky) Base() int {
	return int(*t)
}

const (
	TypePositionStickyBottomCenter TypePositionSticky = iota + 1
	TypePositionStickyBottomLeft
	TypePositionStickyBottomRight
	TypePositionStickyTop
	TypePositionStickyBottom
	TypePositionStickyTopCenter
)

func (t TypePositionSticky) String() string {
	switch t {
	case TypePositionStickyBottomCenter:
		return "Bottom Center"
	case TypePositionStickyBottomLeft:
		return "Bottom Left"
	case TypePositionStickyBottomRight:
		return "Bottom Right"
	case TypePositionStickyTop:
		return "Top"
	case TypePositionStickyBottom:
		return "Bottom"
	case TypePositionStickyTopCenter:
		return "Top Center"
	default:
		return ""
	}
}

func (t TypePositionSticky) StringJson() string {
	switch t {
	case TypePositionStickyBottomCenter:
		return "bottom_center"
	case TypePositionStickyBottomLeft:
		return "bottom_left"
	case TypePositionStickyBottomRight:
		return "bottom_right"
	case TypePositionStickyTop:
		return "top"
	case TypePositionStickyBottom:
		return "bottom"
	case TypePositionStickyTopCenter:
		return "top_center"
	default:
		return ""
	}
}

func (t TypePositionSticky) Int() int {
	switch t {
	case TypePositionStickyBottomCenter:
		return 1
	case TypePositionStickyBottomLeft:
		return 2
	case TypePositionStickyBottomRight:
		return 3
	case TypePositionStickyTop:
		return 4
	case TypePositionStickyBottom:
		return 5
	case TypePositionStickyTopCenter:
		return 5
	default:
		return 0
	}
}

type TYPEPassBackType int

const (
	TYPEPassBackTypeInline TYPEPassBackType = iota + 1
	TYPEPassBackTypeCustom
)

func (t TYPEPassBackType) String() string {
	switch t {
	case TYPEPassBackTypeInline:
		return "inline"
	case TYPEPassBackTypeCustom:
		return "custom"
	default:
		return ""
	}
}

func (t TYPEPassBackType) Int() int {
	switch t {
	case TYPEPassBackTypeInline:
		return 1
	case TYPEPassBackTypeCustom:
		return 2
	default:
		return 0
	}
}

type TYPERenderer int

const (
	TYPERendererPubPower TYPERenderer = iota + 1
	TYPERendererOther
	TYPERendererJWPlayer
	TYPERendererVideoJS
	TYPERendererFlowPlayer
	TYPERendererOverlayAd
)

func (t TYPERenderer) String() string {
	switch t {
	case TYPERendererPubPower:
		return "PubPower Player"
	case TYPERendererOther:
		return "Other"
	case TYPERendererJWPlayer:
		return "JWPlayer"
	case TYPERendererVideoJS:
		return "VideoJS"
	case TYPERendererFlowPlayer:
		return "FlowPlayer"
	case TYPERendererOverlayAd:
		return "Overlay Ad"
	default:
		return ""
	}
}

func (t TYPERenderer) Int() int {
	switch t {
	case TYPERendererPubPower:
		return 1
	case TYPERendererOther:
		return 2
	case TYPERendererJWPlayer:
		return 3
	case TYPERendererVideoJS:
		return 4
	case TYPERendererFlowPlayer:
		return 5
	case TYPERendererOverlayAd:
		return 6
	default:
		return 0
	}
}

type TYPERelatedContent int

const (
	TYPERelatedContentNewest TYPERelatedContent = iota + 1
	TYPERelatedContentMostViewed
)

func (t TYPERelatedContent) String() string {
	switch t {
	case TYPERelatedContentNewest:
		return "Newest"
	case TYPERelatedContentMostViewed:
		return "Most viewed"
	default:
		return ""
	}
}

func (t TYPERelatedContent) Code() string {
	switch t {
	case TYPERelatedContentNewest:
		return "newest"
	case TYPERelatedContentMostViewed:
		return "mostviewed"
	default:
		return ""
	}
}

type TYPEContentType int

const (
	TYPEContentTypeRelated TYPEContentType = iota + 1
	TYPEContentTypeQuiz
)

func (t TYPEContentType) Int() int64 {
	switch t {
	case TYPEContentTypeRelated:
		return 1
	case TYPEContentTypeQuiz:
		return 2
	default:
		return 0
	}
}

type TYPEAdSize int

const (
	TYPEAdSizeFixed TYPEAdSize = iota + 1
	TYPEAdSizeResponsive
)

func (t TYPEAdSize) String() string {
	switch t {
	case TYPEAdSizeFixed:
		return "Fixed"
	case TYPEAdSizeResponsive:
		return "Responsive"
	default:
		return ""
	}
}

func (t TYPEAdSize) Int() int64 {
	switch t {
	case TYPEAdSizeFixed:
		return 1
	case TYPEAdSizeResponsive:
		return 2
	default:
		return 0
	}
}

type TYPEResponsiveType int

const (
	TYPEResponsiveTypeHorizontal TYPEResponsiveType = iota + 1
	TYPEResponsiveTypeSquare
	TYPEResponsiveTypeVertical
)

func (t TYPEResponsiveType) String() string {
	switch t {
	case TYPEResponsiveTypeHorizontal:
		return "Horizontal"
	case TYPEResponsiveTypeSquare:
		return "Square"
	case TYPEResponsiveTypeVertical:
		return "Vertical"
	default:
		return ""
	}
}

func (t TYPEResponsiveType) Int() int64 {
	switch t {
	case TYPEResponsiveTypeHorizontal:
		return 1
	case TYPEResponsiveTypeSquare:
		return 2
	case TYPEResponsiveTypeVertical:
		return 3
	default:
		return 0
	}
}

func (t TYPEResponsiveType) GetAllSize() (adSizes []TableAdSize) {
	switch t {
	case TYPEResponsiveTypeHorizontal:
		listSize := []string{"970x250", "970x90", "728x90", "468x60", "320x100", "320x50", "300x50"}
		Client.Where("name in ?", listSize).Find(&adSizes)
		return
	case TYPEResponsiveTypeSquare:
		listSize := []string{"336x280", "300x250"}
		Client.Where("name in ?", listSize).Find(&adSizes)
		return
	case TYPEResponsiveTypeVertical:
		listSize := []string{"300x600", "160x600", "120x600"}
		Client.Where("name in ?", listSize).Find(&adSizes)
		return
	default:
		return
	}
}

func (t TYPEResponsiveType) GetListSize() (listSize []string) {
	switch t {
	case TYPEResponsiveTypeHorizontal:
		listSize = []string{"970x250", "970x90", "728x90", "468x60", "320x100", "320x50", "300x50"}
		return
	case TYPEResponsiveTypeSquare:
		listSize = []string{"336x280", "300x250"}
		return
	case TYPEResponsiveTypeVertical:
		listSize = []string{"300x600", "160x600", "120x600"}
		return
	default:
		return
	}
}

type TYPEAdRefresh string

const (
	TYPEAdRefreshDomainConfiguration        TYPEAdRefresh = "domain_configuration"
	TYPEAdRefreshIndividualTagConfiguration TYPEAdRefresh = "individual_tag_configuration"
)
