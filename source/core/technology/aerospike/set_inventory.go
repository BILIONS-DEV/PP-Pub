package aerospike

import (
	"source/core/technology/mysql"
	"time"
)

type SetInventory struct {
	Uuid                string                 `as:"uuid" json:"uuid"`
	Id                  int64                  `as:"id" json:"id"`
	UserID              int64                  `as:"user_id" json:"user_id"`
	Name                string                 `as:"name" json:"name"`
	Status              mysql.TYPEStatus       `as:"status" json:"status"`
	PrebidTimeOut       int                    `as:"prebid_timeout" json:"prebid_timeout"`
	LoadAdType          string                 `as:"load_ad_type" json:"load_ad_type"`
	MobileScaling       int                    `as:"mobile_scaling" json:"mobile_scaling"`
	FetchMarginPercent  int                    `as:"fetch_margin" json:"fetch_margin"`
	RenderMarginPercent int                    `as:"render_margin" json:"render_margin"`
	AdRefreshTime       int                    `as:"ad_refresh_time" json:"ad_refresh_time"`
	ReloadMode          string                 `as:"reload_mode" json:"reload_mode"`
	SafeFrame           mysql.TypeOnOff        `as:"safe_frame" json:"safe_frame"`
	JsMode              string                 `as:"js_mode" json:"js_mode"`
	PbRenderMode        mysql.TYPEPbRenderMode `as:"pb_render_mode" json:"pb_render_mode"`
	VpaidMode           string                 `as:"vpaid_mode" json:"vpaid_mode"`
	PrebidJs            string                 `as:"prebid_js" json:"prebid_js"`
	ListAdTag           []AdTag                `as:"-" json:"list_ad_tag"`
	KeyAdTags           []string               `as:"keyAdtags" json:"keyAdtags"`
	Bidders             []Bidder               `as:"bidders" json:"bidders"`
	BlockAdDomains      []string               `as:"blockAdDomains" json:"blockAdDomains"`
	UserConfig          mysql.TableConfig      `as:"user_config" json:"user_config"`
	// Consent field
	Gdpr           int                `as:"gdpr" json:"gdpr"`
	GdprTimeout    int                `as:"gdpr_timeout" json:"gdpr_timeout"`
	Ccpa           int                `as:"ccpa" json:"ccpa"`
	CcpaTimeout    int                `as:"ccpa_timeout" json:"ccpa_timeout"`
	CustomBrand    int                `as:"custom_brand" json:"custom_brand"`
	Logo           string             `as:"logo" json:"logo"`
	Title          string             `as:"title" json:"title"`
	Content        string             `as:"content" json:"content"`
	UserIdModule   UserIdModule       `as:"user_id_module" json:"user_id_module"`
	Uid            string             `as:"uid" json:"uid"`
	CreativeIds    []int64            `as:"creative_ids" json:"creative_ids"`
	BlockCreatives []string           `as:"blockCreatives" json:"blockCreatives"`
	GAM            map[int64]Currency `as:"gam" json:"gam"`
	AmazonCpm      map[string]float64 `as:"amazonCpm" json:"amazonCpm"`
	CaPub          string             `as:"caPub" json:"caPub"`
	EmailMd5       string             `as:"emailMd5" json:"emailMd5"`
	DirectSales    string             `as:"directSales" json:"directSales"`
	QuizUid        int64              `as:"quiz_uid" json:"quiz_uid"`
}

type Currency struct {
	Currency              string `as:"currency" json:"currency"`
	GranularityMultiplier int    `as:"granularityMultiplier" json:"granularityMultiplier"`
}

type UserIdModule struct {
	Configs struct {
		SyncDelay    int `as:"sync_delay"`
		AuctionDelay int `as:"auction_delay"`
	} `as:"configs" json:"configs"`
	Modules []UserIdModuleInfo `as:"modules" json:"modules"`
}

type UserIdModuleInfo struct {
	Name                 string `as:"name" json:"name"`
	PrebidModuleFileName string `as:"pbFileName" json:"prebid_module_filename"`
	Params               string `as:"params" json:"params"`
	Storage              string `as:"storage" json:"storage"`
	//AbTesting string `as:"ab_testing" json:"ab_testing"`
	//Volume    int    `as:"volume" json:"volume"`
}

type AdTag struct {
	Id                   int64                    `as:"id" json:"id"`
	Name                 string                   `as:"name" json:"name"`
	AdType               mysql.TYPEAdType         `as:"ad_type" json:"ad_type"`
	Renderer             mysql.TYPERenderer       `as:"renderer" json:"renderer"`
	GAM                  string                   `as:"gam" json:"gam"`
	GAMExt               string                   `as:"gam_ext" json:"gam_ext"`
	ShiftContent         bool                     `json:"shiftContent"`
	Size                 AdTagSize                `as:"size" json:"size"`
	BidOutStream         mysql.TypeOnOff          `as:"bid_out_stream" json:"bid_out_stream"`
	PassBack             string                   `as:"pass_back" json:"pass_back"`
	PassBackType         mysql.TYPEPassBackType   `as:"pass_back_type" json:"pass_back_type"`
	FeedUrl              string                   `as:"feed_url" json:"feed_url"`
	LineItemPubs         []LineItem               `as:"lineItemPub" json:"line_item_pubs"`
	LineItemAdmins       []LineItem               `as:"lineItemAdmin" json:"line_item_admins"`
	LineItemPubsV2       []LineItem               `as:"lineItemPubV2" json:"line_item_pubs_v2"`
	LineItemAdminsV2     []LineItem               `as:"lineItemAdminV2" json:"line_item_admins_v2"`
	LineItemPubMobiles   []LineItem               `as:"linePubMobile" json:"line_item_pub_mobiles"`
	LineItemAdminMobiles []LineItem               `as:"lineAdminMobile" json:"line_item_admin_mobiles"`
	LineForExt           LineForExt               `as:"lineForExt" json:"line_for_ext"`
	LineForBanner        LineForBanner            `as:"lineForBanner" json:"line_for_banner"`
	LineItemGoogle       LineItem                 `as:"lineGoogle"  json:"line_item_google"`
	LineItemGoogleExt    LineItem                 `as:"lineGoogleExt"  json:"line_item_google_ext"`
	Floors               []Floor                  `as:"floors" json:"floors"`
	FloorTests           []Floor                  `as:"floor_tests" json:"floor_tests"`
	Video                Video                    `as:"video" json:"video"`
	Position             mysql.TypePositionSticky `as:"position" json:"position"`
	CloseButton          int                      `as:"close_button" json:"close_button"`
	MobileConfig         AdTagMobileConfig        `as:"mobile_config" json:"mobile_config"`
	OrderMethod          string                   `as:"order_method" json:"order_method"` // Cho video
	BlockRPM             map[string][]string      `as:"block_rpm" json:"block_rpm"`
	FrequencyCaps        int                      `as:"frequencyCaps" json:"frequency_caps"`
	MainTitle            string                   `as:"main_title" json:"main_title"`
	BackgroundColor      string                   `as:"backgroundColor" json:"background_color"`
	TitleColor           string                   `as:"title_color" json:"title_color"`
	TotalAds             int64                    `as:"total_ads" json:"total_ads"`
	LayoutType           int64                    `as:"layout_type" json:"layout_type"`
	AdSize               mysql.TYPEAdSize         `as:"ad_size" json:"ad_size"`
	ResponsiveType       mysql.TYPEResponsiveType `as:"responsiveType" json:"responsive_type"`
	ContentType          mysql.TYPEContentType    `as:"contentType" json:"content_type"`
	EnableBanner         mysql.TYPEOnOff          `as:"enable_banner" json:"enable_banner"`
	EnableVideo          mysql.TYPEOnOff          `as:"enable_video" json:"enable_video"`
	AdsenseTag           AdsenseTag               `as:"adsense_tag" json:"adsense_tag"`
	AdRefreshTime        interface{}              `as:"adRefreshTime" json:"adRefreshTime"`
	EnableStickyDesktop  mysql.TYPEOnOff          `as:"eStickyDesktop" json:"eStickyDesktop"`
	EnableStickyMobile   mysql.TYPEOnOff          `as:"eStickyMobile" json:"eStickyMobile"`
	Title                string                   `as:"title" json:"title"`
	BtnApproved          string                   `as:"btnApproved" json:"btn_approved"`
	BtnCancel            string                   `as:"btnCancel" json:"btn_cancel"`
	Amount               string                   `as:"amount" json:"amount"`
	Type                 string                   `as:"type" json:"type"`
	MessageSuccess       string                   `as:"messageSuccess" json:"message_success"`
	BtnClose             string                   `as:"btnClose" json:"btn_close"`
	AdMode               string                   `as:"adMode" json:"adMode"`
}

type LineForExt struct {
	LineItemPubs   []LineItem `as:"lineItemPub" json:"line_item_pubs"`
	LineItemAdmins []LineItem `as:"lineItemAdmin" json:"line_item_admins"`
}

type LineForBanner struct {
	LineItemPubs   []LineItem `as:"lineItemPub" json:"line_item_pubs"`
	LineItemAdmins []LineItem `as:"lineItemAdmin" json:"line_item_admins"`
}

type AdTagMobileConfig struct {
	Size     string `as:"size" json:"size"`
	PassBack string `as:"passback,omitempty" json:"passback,omitempty"`
	Position string `as:"position,omitempty" json:"position,omitempty"`
	CloseBtn bool   `as:"closeBtn,omitempty" json:"closeBtn,omitempty"`
}

type Bidder struct {
	UserId           int64    `as:"user_id" json:"user_id"`
	Id               int64    `as:"id" json:"id"`
	BidderTemplateId int64    `as:"bidder_template_id" json:"bidder_template_id"`
	DisplayName      string   `as:"display_name" json:"display_name"`
	BidderCode       string   `as:"bidder_code" json:"bidder_code"`
	BidderAlias      string   `as:"bidder_alias" json:"bidder_alias"`
	MediaType        []string `as:"media_type" json:"media_type"`
	BidAdjustment    float64  `as:"bid_adjustment" json:"bid_adjustment"`
	PubId            string   `as:"pub_id" json:"pub_id"`
	Params           string   `as:"params" json:"params"`
}

type AdTagSize struct {
	DesktopPrimary   mysql.TableAdSize   `as:"desktop_primary"`
	DesktopAdditions []mysql.TableAdSize `as:"desktop_additions"`
	MobilePrimary    mysql.TableAdSize   `as:"mobile_primary"`
	MobileAdditions  []mysql.TableAdSize `as:"mobile_additions"`
}

type LineItem struct {
	Id             int64                    `as:"id" json:"id"`
	Type           int                      `as:"type" json:"type"`
	Priority       int                      `as:"priority" json:"priority"`
	StartDate      time.Time                `as:"start_date" json:"start_date"`
	EndDate        time.Time                `as:"end_date" json:"end_date"`
	Target         *Target                  `as:"target" json:"target"`
	BidderInfo     []BidderInfo             `as:"bidder_info" json:"bidder_info"`
	BidderGoogle   mysql.TableBidder        `as:"bidder_google" json:"bidder_google"`
	LinkedGam      int64                    `as:"linked_gam" json:"linked_gam"`
	ConnectionType mysql.TYPEConnectionType `as:"connection_type" json:"connection_type"`
	UserId         int64                    `as:"user_id" json:"user_id"`
}

type BidderInfo struct {
	Id         int64                `as:"id" json:"id"`
	Name       string               `as:"name" json:"name"`
	ConfigType mysql.TYPEConfigType `as:"configType" json:"configType"`
	Type       mysql.TYPEBidderType `as:"type" json:"type"`
	Bidder     mysql.TableBidder    `as:"bidder" json:"bidder"`
	Params     string               `as:"params" json:"params"`
}

type Floor struct {
	Id            int64                `as:"id" json:"id"`
	FloorType     mysql.TYPEFloor      `as:"floor_type" json:"floor_type"`
	FloorValue    float64              `as:"floor_value" json:"floor_value"`
	FloorDynamic  map[string]float64   `as:"floor_dynamic" json:"floor_dynamic"`
	FloorTest     map[string]FloorTest `as:"floor_test" json:"floor_test"`
	TimeTest      int64                `as:"time_test" json:"time_test"`
	AbTestId      int64                `as:"ab_test_id" json:"ab_test_id"`
	Priority      int                  `as:"priority" json:"priority"`
	PricingRuleID int64                `as:"pricingRuleID" json:"pricingRuleID"`
	Target        *Target              `as:"target" json:"target"`
}

type FloorTest struct {
	Floor []float64 `as:"floor" json:"floor"`
	Id    int64     `as:"id" json:"id"`
}

type Target struct {
	Geos    []string
	Devices []string
	Sizes   []string `as:"sizes,omitempty" json:"sizes,omitempty"`
	Hour    []int
}

type Video struct {
	Template mysql.TableTemplate `as:"template" json:"template"`
	PlayList PlayList            `as:"playlist" json:"playlist"`
	Contents []Content           `as:"contents" json:"contents"`
	Info     []Info              `as:"info" json:"info"`
}

type Content struct {
	Link        string   `json:"link"`
	Thumb       string   `json:"thumb"`
	Title       string   `json:"title"`
	VideoURL    VideoURL `json:"video_url"`
	ContentDesc string   `json:"content_desc"`
	AdsTime     AdsTime  `json:"ads_time"`
}

type AdsTime struct {
	Preroll  *AdBreak  `json:"preroll"`
	Midroll  []AdBreak `json:"midroll"`
	Postroll *AdBreak  `json:"postroll"`
}

type AdBreak struct {
	IsShow  bool  `json:"isShow"`
	AdsNums int   `json:"adsNums"`
	Time    int64 `json:"time"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type VideoURL struct {
	M3U8 string `json:"m3u8"`
	Mp4  string `json:"mp4"`
	Ogg  string `json:"ogg"`
}

type PlayList struct {
	Id          int64     `gorm:"column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Contents    []Content `as:"contents" json:"contents"`
}

type AdsenseTag struct {
	AdSlotID   string `as:"adslot_id" json:"adslot_id"`
	PubID      string `as:"pub_id" json:"pub_id"`
	AdUnitType string `as:"adunit_type" json:"adunit_type"`
	Render     string `as:"render" json:"render"`
}
