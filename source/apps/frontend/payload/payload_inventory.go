package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type InventorySubmit struct {
	Inventories string `json:"inventories" form:"inventories"`
}

type InventoryIndex struct {
	InventoryFilterPostData
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	OrderColumn int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string   `query:"order_dir" json:"order_dir" form:"order_dir"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Type        []string `query:"f_type" form:"f_type" json:"f_type"`
	AdsTxtSync  []string `query:"f_ads_sync" form:"f_ads_sync" json:"f_ads_sync"`
	WebLive     []string `query:"f_web_live" form:"f_web_live" json:"f_web_live"`
	User        []int64  `query:"f_user" form:"f_user" json:"f_user"`
}

type InventoryFilterPayload struct {
	datatable.Request
	PostData *InventoryFilterPostData `query:"postData"`
}

type InventoryFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	AdsTxtSync  interface{} `query:"f_ads_sync[]" json:"f_ads_sync" form:"f_ads_sync[]"`
	WebLive     interface{} `query:"f_web_live[]" json:"f_web_live" form:"f_web_live[]"`
	User        interface{} `query:"f_user[]" json:"f_user" form:"f_user[]"`
	// datatable input
	Page   int `query:"page" json:"page" form:"page"`
	Limit  int `query:"limit" json:"limit" form:"limit"`
	Start  int `query:"start" json:"start" form:"start"`
	Length int `query:"length" json:"length" form:"length"`
}

type PayloadInventoryChangeStatusConnection struct {
	InventoryId int64                                     `form:"inventory_id"`
	BidderId    int64                                     `form:"bidder_id"`
	Status      mysql.TYPEStatusInventoryConnectionDemand `form:"status"`
}

type GeneralInventory struct {
	Id            int64                  `json:"id"`
	InventoryId   int64                  `json:"inventory_id"`
	GamAccount    int64                  `json:"gam_account"`
	AdRefreshTime int                    `json:"ad_refresh_time"`
	PrebidTimeout int                    `json:"prebid_timeout"`
	LoadAdType    string                 `json:"load_ad_type"`
	PbRenderMode  mysql.TYPEPbRenderMode `json:"pb_render_mode"`
	AdRefresh     mysql.TypeOnOff        `form:"ad_refresh" json:"ad_refresh"`
	DirectSales   mysql.TypeOnOff        `form:"direct_sales" json:"direct_sales"`
	LoadAdRefresh string                 `query:"ad_refresh_type" json:"ad_refresh_type" form:"ad_refresh_type"`
	Gdpr          int                    `json:"gdpr"`
	GdprTimeout   int                    `json:"gdpr_timeout"`
	Ccpa          int                    `json:"ccpa"`
	CcpaTimeout   int                    `json:"ccpa_timeout"`
	Status        mysql.TYPEStatus       `json:"status"`
	UserId        int64                  `json:"user_id"`
	CustomBrand   int                    `json:"custom_brand"`
	SafeFrame     mysql.TypeOnOff        `json:"safe_frame"`
	Logo          string                 `json:"custom_brand_logo"`
	Title         string                 `json:"custom_brand_title"`
	Content       string                 `json:"custom_brand_content"`
	AuctionDelay  int                    `form:"column:auction_delay" json:"auction_delay"`
	SyncDelay     int                    `form:"column:sync_delay" json:"sync_delay"`
	// ModuleTesting int              `form:"column:module_testing" json:"module_testing"`
	// ModuleVolume  int              `form:"column:module_volume" json:"module_volume"`
	ModuleParams        []ModuleInfo    `json:"module_params" form:"module_params"`
	ListDel             []int64         `json:"list_del" form:"list_del"`
	GamAutoCreate       mysql.TypeOnOff `json:"gam_auto_create" form:"gam_auto_create"`
	FetchMarginPercent  int             `json:"fetch_margin_percent" json:"fetch_margin_percent"`
	RenderMarginPercent int             `json:"render_margin_percent" json:"render_margin_percent"`
	MobileScaling       int             `json:"mobile_scaling" json:"mobile_scaling"`
}

type Config struct {
	Id            int64  `json:"id"`
	InventoryId   int64  `json:"inventory_id"`
	AdRefreshTime int    `json:"ad_refresh_time"`
	PrebidTimeout int    `json:"prebid_timeout"`
	LoadAdType    string `json:"load_ad_type"`
	AdRefresh     int    `form:"ad_refresh" json:"ad_refresh"`
	LoadAdRefresh string `query:"ad_refresh_type" json:"ad_refresh_type" form:"ad_refresh_type"`
}

type Consent struct {
	Id          int64            `json:"id"`
	InventoryId int64            `json:"inventory_id"`
	Gdpr        int              `json:"gdpr"`
	GdprTimeout int              `json:"gdpr_timeout"`
	Ccpa        int              `json:"ccpa"`
	CcpaTimeout int              `json:"ccpa_timeout"`
	Status      mysql.TYPEStatus `json:"status"`
	UserId      int64            `json:"user_id"`
	CustomBrand int              `json:"custom_brand"`
	Logo        string           `json:"custom_brand_logo"`
	Title       string           `json:"custom_brand_title"`
	Content     string           `json:"custom_brand_content"`
}

type UserID struct {
	Id               int64 `json:"id"`
	InventoryId      int64 `json:"inventory_id"`
	AuctionDelay     int   `form:"column:auction_delay" json:"auction_delay"`
	Module           int   `form:"column:height" json:"height"`
	SharedTesting    int   `form:"column:shared_testing" json:"shared_testing"`
	SharedVolume     int   `form:"column:shared_volume" json:"shared_volume"`
	CriteoTesting    int   `form:"column:criteo_testing" json:"criteo_testing"`
	CriteoVolume     int   `form:"column:criteo_volume" json:"criteo_volume"`
	FlocTesting      int   `form:"column:floc_testing" json:"floc_testing"`
	FlocVolume       int   `form:"column:floc_volume" json:"floc_volume"`
	UniversalTesting int   `form:"column:universal_testing" json:"universal_testing"`
	UniversalVolume  int   `form:"column:universal_volume" json:"universal_volume"`
}

type SetupInventory struct {
	AdRefreshTime int     `json:"ad_refresh_time"`
	PrebidTimeout int     `json:"prebid_timeout"`
	LoadAdType    string  `json:"load_ad_type"`
	LoadAdRefresh string  `json:"load_ad_refresh"`
	AdsTxt        string  `json:"ads_txt"`
	Status        string  `json:"status"`
	AdTagName     string  `json:"ad_tag_name"`
	FloorPrice    float64 `json:"floor_price"`
	Type          string  `json:"type"`
	Size          string  `json:"size"`
	Gam           string  `json:"gam"`
	PassBack      string  `json:"pass_back"`
	IsPublish     string  `json:"is_publish"`
}

type ModuleParam struct {
	ModuleId    int64
	ModuleName  string
	ModuleIndex int
	Params      []ParamModuleUserId
	Storage     []StorageModuleUserId
	ParamValue  []ModuleInfo
}

type ModuleInfo struct {
	ModuleId    int64                 `json:"id"`
	ModuleName  string                `json:"name"`
	ModuleIndex int                   `json:"-"`
	Params      []ParamModuleUserId   `json:"params"`
	Storage     []StorageModuleUserId `json:"storage"`
	AbTesting   mysql.TYPEOnOff       `json:"ab_testing"`
	Volume      int                   `json:"volume"`
	// ParamsType  map[string]string     `json:"-"`
}

type ParamModuleUserId struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template"`
}

type StorageModuleUserId struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template"`
}

type ModuleUserIdAssign struct {
	Id         int64                 `gorm:"column:id" json:"id"`
	IdentityId int64                 `gorm:"column:identity_id" json:"identity_id"`
	ModuleId   int64                 `gorm:"column:module_id" json:"module_id"`
	Name       string                `gorm:"column:name" json:"name"`
	Params     []ParamModuleUserId   `gorm:"params" json:"params" form:"params"`
	Storage    []StorageModuleUserId `gorm:"storage" json:"storage" form:"storage"`
	AbTesting  mysql.TYPEOnOff       `gorm:"ab_testing" json:"ab_testing" form:"ab_testing"`
	Volume     int                   `gorm:"volume" json:"volume" form:"volume"`
}

type DefaultModule struct {
	Id      int64
	Name    string
	Params  string
	Storage string
}

type BuildScript struct {
	TagDesktop      int64           `json:"desktop_tag"`
	TagMobile       int64           `json:"mobile_tag"`
	PlaceHolder     mysql.TypeOnOff `json:"place_holder"`
	PlaceHolderText string          `json:"place_holder_text"`
	TextColor       string          `json:"text_color"`
	BorderColor     string          `json:"border_color"`
}

type ResponseApiAdsense struct {
	InsertedTo []string `json:"inserted_to"`
	InsertedID []string `json:"inserted_id"`
	Queries    []struct {
		ID   string `json:"id"`
		SQL  string `json:"sql"`
		Data struct {
			Domain   string `json:"domain"`
			DomainID string `json:"domain_id"`
			Adsense  string `json:"adsense"`
		} `json:"data"`
	} `json:"queries"`
	Status bool `json:"status"`
}
