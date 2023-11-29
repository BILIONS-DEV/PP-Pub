package lang

type Inventory struct {
	CreateAdTag TYPETranslation `json:"create_ad_tag"`
	AdTagEmpty  TYPETranslation `json:"ad_tag_empty"`

	SubmitDomain SubmitDomain `json:"submit_domain"`

	Setup struct {
		TabConfig  TabConfig  `json:"tab_config"`
		TabConsent TabConsent `json:"tab_consent"`
		TabUserID  TabUserID  `json:"tab_user_id"`
	} `json:"setup"`
}

type TabUserID struct {
	SyncDelay        TYPETranslation `json:"sync_delay"`
	SyncDelayDesc    TYPETranslation `json:"sync_delay_desc"`
	AuctionDelay     TYPETranslation `json:"auction_delay"`
	AuctionDelayDesc TYPETranslation `json:"auction_delay_desc"`
	Module           TYPETranslation `json:"module"`
	ModuleDesc       TYPETranslation `json:"module_desc"`
	Button           TYPETranslation `json:"button"`
}

type TabConfig struct {
	GAMAccount              TYPETranslation `json:"gam_account"`
	GAMAccountPlaceholder   TYPETranslation `json:"gam_account_placeholder"`
	GAMAdUnitAutoCreate     TYPETranslation `json:"gam_ad_unit_auto_create"`
	GAMAdUnitAutoCreateDesc TYPETranslation `json:"gam_ad_unit_auto_create_desc"`
	SafeFrame               TYPETranslation `json:"safe_frame"`
	PrebidTimeout           TYPETranslation `json:"prebid_timeout"`
	PrebidTimeoutDesc       TYPETranslation `json:"prebid_timeout_desc"`
	LoadAdType              TYPETranslation `json:"load_ad_type"`
	LoadAdTypeDesc          TYPETranslation `json:"load_ad_type_desc"`
	LoadAdTypePlaceholder   TYPETranslation `json:"load_ad_type_placeholder"`
	PbRenderMode            TYPETranslation `json:"pb_render_mode"`
	PbRenderModeDesc        TYPETranslation `json:"pb_render_mode_desc"`
	PbRenderModePlaceholder TYPETranslation `json:"pb_render_mode_placeholder"`
	AdRefresh               TYPETranslation `json:"ad_refresh"`
	AdRefreshDesc           TYPETranslation `json:"ad_refresh_desc"`
	AdRefreshTime           TYPETranslation `json:"ad_refresh_time"`
	AdRefreshTimeDesc       TYPETranslation `json:"ad_refresh_time_desc"`
	AdRefreshType           TYPETranslation `json:"ad_refresh_type"`
	AdRefreshTypeDesc       TYPETranslation `json:"ad_refresh_type_desc"`
	Button                  TYPETranslation `json:"button"`
	FetchMarginPercent      TYPETranslation `json:"fetch_margin_percent"`
	FetchMarginPercentDesc  TYPETranslation `json:"fetch_margin_percent_desc"`
	RenderMarginPercent     TYPETranslation `json:"render_margin_percent"`
	RenderMarginPercentDesc TYPETranslation `json:"render_margin_percent_desc"`
	MobileScaling           TYPETranslation `json:"mobile_scaling"`
	MobileScalingDesc       TYPETranslation `json:"mobile_scaling_desc"`
	DirectSales             TYPETranslation `json:"direct_sales"`
	DirectSalesDesc         TYPETranslation `json:"direct_sales_desc"`
}

type TabConsent struct {
	GDPR               TYPETranslation `json:"gdpr"`
	GDPRDesc           TYPETranslation `json:"gdpr_desc"`
	CCPA               TYPETranslation `json:"ccpa"`
	CCPADesc           TYPETranslation `json:"ccpa_desc"`
	Timeout            TYPETranslation `json:"timeout"`
	TimeoutDesc        TYPETranslation `json:"timeout_desc"`
	CustomBrand        TYPETranslation `json:"custom_brand"`
	CustomBrandDesc    TYPETranslation `json:"custom_brand_desc"`
	Logo               TYPETranslation `json:"logo"`
	LogoDesc           TYPETranslation `json:"logo_desc"`
	LogoPlaceholder    TYPETranslation `json:"logo_placeholder"`
	Title              TYPETranslation `json:"title"`
	TitleDesc          TYPETranslation `json:"title_desc"`
	TitlePlaceholder   TYPETranslation `json:"title_placeholder"`
	Content            TYPETranslation `json:"content"`
	ContentDesc        TYPETranslation `json:"content_desc"`
	ContentPlaceholder TYPETranslation `json:"content_placeholder"`
}

type SubmitDomain struct {
	ModalTitle TYPETranslation `json:"modal_title"`
	Title      TYPETranslation `json:"title"`
	Desc       TYPETranslation `json:"desc"`
	Example    TYPETranslation `json:"example"`
	Button     TYPETranslation `json:"button"`
	Errors     struct {
		AlreadyExist TYPETranslation `json:"already_exist"`
	} `json:"errors"`
}

type InventoryError struct {
	Submit       TYPETranslation `json:"submit"`
	SetupConfig  TYPETranslation `json:"setup_config"`
	SetupConsent TYPETranslation `json:"setup_consent"`
	SetupUserId  TYPETranslation `json:"setup_userid"`
	Edit         TYPETranslation `json:"edit"`
	Delete       TYPETranslation `json:"delete"`
	List         TYPETranslation `json:"list"`
	NotFound     TYPETranslation `json:"not_found"`
	RootDomain   TYPETranslation `json:"root_domain"`
	SaveAdsTxt   TYPETranslation `json:"save_ads_txt"`
	ScanAdsTxt   TYPETranslation `json:"scan_ads_txt"`
}
