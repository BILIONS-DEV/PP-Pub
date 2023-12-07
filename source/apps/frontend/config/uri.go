package config

const (
	URIHome       = "/"
	URIDashboards = "/dashboards"
	URIITest      = "/test"

	URIUser           = "/user"
	URIRegister       = "/register"
	URILogin          = "/login"
	URILogout         = "/user/logout"
	URIBillingSetting = "/user/billing"
	URIAccountSetting = "/user/account"
	URIChangePassWord = "/user/password"
	URIForgotPassWord = "/user/forgot-password"
	URIResetPassWord  = "/user/new-password"
	URIChangePassword = "/user/changePassword"
	URIAtl            = "/altlog"
	URIAtlQuick       = "/altlogq"
	// => UserBackend
	URIUserBe = "/be/user"

	URIInventory                       = "/websites"
	URIInventorySubmit                 = "/websites/submit"
	URIInventorySetup                  = "/websites/setup"
	URIInventoryConsent                = "/websites/setupConsent"
	URIInventoryUserId                 = "/websites/setupUserId"
	URIInventoryDelete                 = "/websites/del"
	URIInventoryAdTag                  = "/websites/adtag"
	URIInventoryLoadParam              = "/websites/loadParam"
	URIInventoryCollapse               = "/websites/collapse"
	URIInventoryCopyAdTag              = "/websites/copyAdTag"
	URIInventoryBuildScript            = "/websites/buildScript"
	URIInventoryConnection             = "/websites/connection"
	URIInventoryChangeStatusConnection = "/websites/change-status-connection"

	URIInventorySetupV2 = "/supply-v2/setup"

	URILineItem              = "/line-item"
	URILineItemLoadParam     = "/line-item/loadParam"
	URILineItemAdd           = "/line-item/add"
	URILineItemEdit          = "/line-item/edit"
	URILineItemDelete        = "/line-item/del"
	URILineItemView          = "/line-item/view"
	URILineItemCreate        = "/line-item/create"
	URILineItemChoose        = "/line-item/choose"
	URISearchDomain          = "/line-item/searchDomain"
	URISearchAdFormat        = "/line-item/searchAdFormat"
	URISearchAdSize          = "/line-item/searchAdSize"
	URISearchAdTag           = "/line-item/searchAdTag"
	URISearchDevice          = "/line-item/searchDevice"
	URISearchCountry         = "/line-item/searchCountry"
	URILineItemCollapse      = "/line-item/collapse"
	URILineItemCheckParam    = "/line-item/checkParam"
	URILineItemListLinkedGam = "/line-item/listLinkedGam"
	URIAddParamBidder        = "/line-item/addParamBidder"

	//Line V2
	URILineItemV2              = "/line-item-v2"
	URILineItemV2LoadParam     = "/line-item-v2/loadParam"
	URILineItemV2Add           = "/line-item-v2/add"
	URILineItemV2Edit          = "/line-item-v2/edit"
	URILineItemV2Delete        = "/line-item-v2/del"
	URILineItemV2View          = "/line-item-v2/view"
	URILineItemV2Create        = "/line-item-v2/create"
	URILineItemV2Choose        = "/line-item-v2/choose"
	URISearchDomainV2          = "/line-item-v2/searchDomain"
	URISearchAdFormatV2        = "/line-item-v2/searchAdFormat"
	URISearchAdSizeV2          = "/line-item-v2/searchAdSize"
	URISearchAdTagV2           = "/line-item-v2/searchAdTag"
	URISearchDeviceV2          = "/line-item-v2/searchDevice"
	URISearchCountryV2         = "/line-item-v2/searchCountry"
	URILineItemV2Collapse      = "/line-item-v2/collapse"
	URILineItemV2CheckParam    = "/line-item-v2/checkParam"
	URILineItemV2ListLinkedGam = "/line-item-v2/listLinkedGam"
	URIAddParamBidderV2        = "/line-item-v2/addParamBidder"

	URITarget              = "/target"
	URITargetAdd           = "/target/add"
	URITargetEdit          = "/target/edit"
	URITargetLoadInventory = "/target/loadInventory"
	URITargetLoadSelected  = "/target/loadSelected"
	URITargetFilterAdTAg   = "/target/filterAdTag"

	URIAdTag                  = "/adtag"
	URIAdTagAdd               = "/adtag/add"
	URIAdTagEdit              = "/adtag/edit"
	URIAdTagDel               = "/adtag/del"
	URIAdTagGetSizeAdditional = "/adtag/getSizeAdditional"
	URIAdTagCollapse          = "/adtag/collapse"

	URIAdTagAddV2  = "/adtag-v2/add"
	URIAdTagEditV2 = "/adtag-v2/edit"

	URIPlayerTemplate          = "/player/template"
	URIPlayerAddTemplate       = "/player/template/add"
	URIPlayerEditTemplate      = "/player/template/edit"
	URIPlayerDelTemplate       = "/player/template/del"
	URIPlayerViewTemplate      = "/player/template/view"
	URIPlayerDuplicateTemplate = "/player/template/duplicate"
	URIPlayerCollapse          = "/player/template/collapse"
	URIPlayerPreview           = "/player/template/preview"

	URIPlayerTemplateV2     = "/player-v2/template"
	URIPlayerAddTemplateV2  = "/player-v2/template/add"
	URIPlayerEditTemplateV2 = "/player-v2/template/edit"
	URIPlayerViewTemplateV2 = "/player-v2/template/view"
	URIPlayerPreviewV2      = "/player-v2/template/preview"

	URIPlaylist         = "/playlist"
	URIPlaylistAdd      = "/playlist/add"
	URIPlaylistEdit     = "/playlist/edit"
	URIPlaylistView     = "/playlist/view"
	URIPlaylistDel      = "/playlist/del"
	URIPlaylistCollapse = "/playlist/collapse"

	URIReport          = "/report"
	URIReportTop       = "/report/top"
	URIReportDimension = "/report/dimension"
	URIReportResponse  = "/report/response"
	URIReportSaved     = "/report/saved"

	URIPayment = "/payment"
	// URIPaymentPreview = "/payment/preview"
	// URIPaymentExport  = "/payment/export-pdf"

	URIFloor         = "/floor"
	URIFloorAdd      = "/floor/add"
	URIFloorEdit     = "/floor/edit"
	URIFloorDel      = "/floor/del"
	URIFloorCollapse = "/floor/collapse"

	URIGam               = "/gam"
	URIGamAdd            = "/gam/add"
	URIGamEdit           = "/gam/edit"
	URIGamPushLine       = "/gam/pushLine"
	URIGamConnect        = "/gam/connect"
	URIGamCallback       = "/gam/callback"
	URIGamSelectNetwork  = "/gam/select-network"
	URIGamGetNetworks    = "/gam/get-networks"
	URIGamCheckApiAccess = "/gam/checkApiAccess"

	URIContent          = "/content"
	URIContentAdd       = "/content/add"
	URIContentEdit      = "/content/edit"
	URIContentDel       = "/content/del"
	URIContentAddVideo  = "/content/add/video"
	URIContentEditVideo = "/content/edit/video"
	URIContentQuiz      = "/quiz"
	URIContentAddQuiz   = "/content/quiz/add"
	URIContentEditQuiz  = "/content/quiz/edit"
	URIContentDelQuiz   = "/content/quiz/del"

	// System
	URISystemBidder            = "/bidder"
	URISystemBidderAdd         = "/bidder/add"
	URISystemBidderAddTemplate = "/bidder/addTemplate"
	URISystemBidderEdit        = "/bidder/edit"
	URISystemBidderDel         = "/bidder/del"
	URISystemBidderView        = "/bidder/view"
	URISystemBidderUploadXlsx  = "/bidder/uploadXlsx"
	URISystemAddParam          = "/bidder/addParam"

	URIConfig     = "/config"
	URIConfigSave = "/config/save"
	URIConfigAdd  = "/config/add"
	URIConfigEdit = "/config/edit"

	URILinkVideo = "/video"

	URIAdsTxt       = "/ads_txt"
	URIAdsTxtDetail = "/ads_txt/detail"
	URIAdsTxtScan   = "/ads_txt/scan"
	URIAdsTxtLoad   = "/ads_txt/load"

	URIBlocking               = "/blocking"
	URIBlockingAdd            = "/blocking/add"
	URIBlockingEdit           = "/blocking/edit"
	URIBlockingLoadInventory  = "/blocking/load_inventory"
	URIBlockingDelete         = "/blocking/del"
	URIBlockingValidateDomain = "/blocking/validateDomain"

	URISUPPORT                    = "/support"
	URISUPPORTPRODUCTDESCRIOPTION = "/support/product-description"
	URISUPPORTTICKETSNEW          = "/support/tickets/new"

	URILineItemSystem       = "/line_item_system"
	URILineItemSystemAdd    = "/line_item_system/add"
	URILineItemSystemEdit   = "/line_item_system/edit"
	URILineItemSystemDelete = "/line_item_system/del"

	URINOTIFICATION        = "/notification"
	URINOTIFICATIONREADALL = "/notification/ReadAll"
	URINOTIFICATIONREAD    = "/notification/read"

	// Identity
	URIIdentity             = "/identity"
	URIIdentityAdd          = "/identity/add"
	URIIdentityEdit         = "/identity/edit"
	URIIdentityDel          = "/identity/del"
	URIIdentityChangeStatus = "/identity/change-status"
	URIIdentityCollapse     = "/identity/collapse"
	// Channels
	URIChannels     = "/channels"
	URIChannelsAdd  = "/channels/add"
	URIChannelsEdit = "/channels/edit"
	URIChannelsDel  = "/channels/del"

	// Ab Testing
	URIAbTesting         = "/ab-testing"
	URIAbTestingAdd      = "/ab-testing/add"
	URIAbTestingEdit     = "/ab-testing/edit"
	URIAbTestingDel      = "/ab-testing/del"
	URIAbTestingCollapse = "/ab-testing/collapse"

	// History
	URIHistory      = "/activity"
	URIHistoryLoad  = "/history/load-histories"
	URIObjectByPage = "/history/object-by-page"

	// Rule
	URIRule = "/rule"

	// Rule Blocked Page
	URIBlockedPageAdd       = "/blocked-page/add"
	URIBlockedPageEdit      = "/blocked-page/edit"
	URIBlockedPageDel       = "/blocked-page/del"
	URIBlockedPageImportCSV = "/blocked-page/importCSV"

	// Advertising Schedules
	URIAdvertisingSchedules       = "/ad-schedules"
	URIAdvertisingSchedulesAdd    = "/ad-schedules/add"
	URIAdvertisingSchedulesEdit   = "/ad-schedules/edit"
	URIAdvertisingSchedulesDelete = "/ad-schedules/delete/:id"
	URIAdvertisingSchedulesDetail = "/ad-schedules/:id"

	// Sales
	URISalesReport = "/report/sales"
	URICampaign    = "/campaigns"

	// AdBlock
	URIAdBlockAnalytics      = "/adblock/analytics"
	URIAdBlockAlertGenerator = "/adblock/generator"
)

type SidebarSetupUri struct {
	Bidder []string
	Config []string

	Report    []string
	Top       []string
	Dimension []string
	Video     []string
	Response  []string
	Saved     []string

	Inventory []string
	Demand    []string
	Rule      []string
	Floor     []string
	Gam       []string
	AdsTxt    []string
	AbTesting []string

	Template []string
	Channels []string
	Playlist []string
	Content  []string
	Quiz     []string
	// GROUP
	SystemGroup  []string
	SetupGroup   []string
	VideoGroup   []string
	ReportGroup  []string
	SupportGroup []string
	SalesGroup   []string

	// Support
	Support []string

	Blocking []string

	Identity []string

	Payment []string

	History []string

	SalesReport []string
	Campaign    []string

	AdBlock []string
}

var SidebarSetup SidebarSetupUri

func init() {
	SidebarSetup.Bidder = append(SidebarSetup.Bidder,
		URISystemBidder, URISystemBidderAdd, URISystemBidderEdit, URISystemBidderAddTemplate, URISystemBidderDel, URISystemBidderView,
	)
	SidebarSetup.Config = append(SidebarSetup.Config,
		URIConfig, URIConfigAdd, URIConfigEdit, URIConfigSave,
	)
	SidebarSetup.Inventory = append(SidebarSetup.Inventory,
		URIInventory, URIInventoryAdTag, URIInventorySetup, URIInventoryDelete, URIInventorySubmit,
		URIAdTag, URIAdTagAdd, URIAdTagEdit, // => AdTag
	)
	SidebarSetup.Demand = append(SidebarSetup.Demand,
		URILineItem, URILineItemAdd, URILineItemEdit, URILineItemCreate, URILineItemView, // => LineItem
	)
	SidebarSetup.Rule = append(SidebarSetup.Rule,
		URIRule, URIFloorAdd, URIFloorEdit, URIBlockingAdd, URIBlockingEdit, URIBlockedPageAdd, URIBlockedPageEdit, // => Bidder
	)
	SidebarSetup.Floor = append(SidebarSetup.Floor,
		URIFloor, URIFloorAdd, URIFloorEdit, // => Floor
	)
	SidebarSetup.Gam = append(SidebarSetup.Gam,
		URIGam, URIGamAdd, URIGamEdit, // => GAM
	)
	SidebarSetup.AdsTxt = append(SidebarSetup.AdsTxt,
		URIAdsTxt, URIAdsTxtDetail, // => AdsTxt
	)
	SidebarSetup.AbTesting = append(SidebarSetup.AbTesting,
		URIAbTesting, URIAbTestingAdd, URIAbTestingEdit, URIAbTestingDel, // => A/B Testing
	)
	SidebarSetup.Template = append(SidebarSetup.Template,
		URIPlayerTemplate, URIPlayerAddTemplate, URIPlayerEditTemplate, URIPlayerViewTemplate, URIPlayerDuplicateTemplate, // => Template
	)
	SidebarSetup.Channels = append(SidebarSetup.Channels,
		URIChannels, URIChannelsAdd, URIChannelsEdit, URIChannelsDel, // => Channels
	)
	SidebarSetup.Playlist = append(SidebarSetup.Playlist,
		URIPlaylist, URIPlaylistAdd, URIPlaylistEdit, URIPlaylistView, // => Playlist
	)
	SidebarSetup.Content = append(SidebarSetup.Content,
		URIContent, URIContentAdd, URIContentEdit, URIContentAddVideo, URIContentEditVideo, // => Content
	)
	SidebarSetup.Quiz = append(SidebarSetup.Quiz,
		URIContentQuiz, URIContentAddQuiz, URIContentEditQuiz, // => Quiz
	)
	SidebarSetup.Blocking = append(SidebarSetup.Blocking,
		URIBlocking, URIBlockingAdd, URIBlockingEdit, // => Blocking
	)
	//SidebarSetup.Report = append(SidebarSetup.Blocking,
	//	URIReport, URIReportTop, URIReportDimension, URIReportResponse, URIReportSaved, // => Report
	//)
	SidebarSetup.Support = append(SidebarSetup.Support,
		URISUPPORT, // => Support
	)
	SidebarSetup.Identity = append(SidebarSetup.Identity,
		URIIdentity, URIIdentityAdd, URIIdentityEdit, // => Support
	)
	SidebarSetup.Payment = append(SidebarSetup.Payment,
		URIPayment, // => Payment
	)
	SidebarSetup.History = append(SidebarSetup.History,
		URIHistory, // => Payment
	)
	SidebarSetup.SalesReport = append(SidebarSetup.SalesReport,
		URISalesReport, // => Direct Sales
	)
	SidebarSetup.Campaign = append(SidebarSetup.Campaign,
		URICampaign, // => Direct Sales
	)
	SidebarSetup.AdBlock = append(SidebarSetup.AdBlock,
		URIAdBlockAnalytics, URIAdBlockAlertGenerator, // => AdBlock
	) // SystemGroup
	SidebarSetup.ReportGroup = append(SidebarSetup.ReportGroup,
		URIReport, URIReportDimension, URIReportSaved, // => Report
	) // ReportGroup
	SidebarSetup.Dimension = append(SidebarSetup.Dimension,
		URIReportDimension, // => Dimension
	) // Dimension
	SidebarSetup.Saved = append(SidebarSetup.Saved,
		URIReportSaved, // => Saved
	) // Saved
	SidebarSetup.Report = append(SidebarSetup.Report,
		URIReport, // => Report
	) // Saved
	SidebarSetup.SystemGroup = append(SidebarSetup.SystemGroup, SidebarSetup.Bidder...)
	SidebarSetup.SystemGroup = append(SidebarSetup.SystemGroup, SidebarSetup.Config...)
	// SetupGroup
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.Inventory...)
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.Demand...)
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.Rule...)
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.Floor...)
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.Gam...)
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.AdsTxt...)
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.AbTesting...)
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.Blocking...)
	SidebarSetup.SetupGroup = append(SidebarSetup.SetupGroup, SidebarSetup.Identity...)
	// VideoGroup
	SidebarSetup.VideoGroup = append(SidebarSetup.VideoGroup, SidebarSetup.Template...)
	SidebarSetup.VideoGroup = append(SidebarSetup.VideoGroup, SidebarSetup.Playlist...)
	SidebarSetup.VideoGroup = append(SidebarSetup.VideoGroup, SidebarSetup.Content...)
	SidebarSetup.VideoGroup = append(SidebarSetup.VideoGroup, SidebarSetup.Quiz...)
	SidebarSetup.VideoGroup = append(SidebarSetup.VideoGroup, SidebarSetup.Channels...)
	// SupportGroup
	SidebarSetup.SupportGroup = append(SidebarSetup.SupportGroup, SidebarSetup.Support...)
	// SalesGroup
	SidebarSetup.SalesGroup = append(SidebarSetup.SalesGroup, SidebarSetup.SalesReport...)
	SidebarSetup.SalesGroup = append(SidebarSetup.SalesGroup, SidebarSetup.Campaign...)
}
