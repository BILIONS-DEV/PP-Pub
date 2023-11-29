package repo

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/repo/account"
	adschedule "source/internal/repo/ad-schedule"
	adType "source/internal/repo/ad-type"
	"source/internal/repo/ads"
	adsense_channel "source/internal/repo/adsense-channel"
	adsense_tag_mapping "source/internal/repo/adsense-tag-mapping"
	aff_block_section "source/internal/repo/aff-block-section"
	"source/internal/repo/bidder"
	"source/internal/repo/billing"
	"source/internal/repo/campaign"
	"source/internal/repo/campaign_traffic_source_id"
	"source/internal/repo/category"
	"source/internal/repo/country"
	"source/internal/repo/creative"
	"source/internal/repo/cronjob"
	"source/internal/repo/gam"
	gamNetwork "source/internal/repo/gam-network"
	googleAdsApi "source/internal/repo/google-ads-api"
	"source/internal/repo/history"
	"source/internal/repo/inventory"
	inventory_ad_tag "source/internal/repo/inventory-ad-tag"
	keyValue "source/internal/repo/key-value"
	keyValueGAM "source/internal/repo/key-value-gam"
	"source/internal/repo/payment"
	"source/internal/repo/quiz"
	report_adsense "source/internal/repo/report-adsense"
	reportAff "source/internal/repo/report-aff"
	report_aff_pixel "source/internal/repo/report-aff-pixel"
	reportAffSelection "source/internal/repo/report-aff-selection"
	reportbodis "source/internal/repo/report-bodis"
	reportBodisTraffic "source/internal/repo/report-bodis-traffic"
	report_codefuel "source/internal/repo/report-codefuel"
	report_domain_active "source/internal/repo/report-domain_active"
	report_facebook "source/internal/repo/report-facebook"
	report_mgid "source/internal/repo/report-mgid"
	reportOpenMailSubID "source/internal/repo/report-openmail-sub-id"
	reportOutBrain "source/internal/repo/report-outbrain"
	report_pocpoc "source/internal/repo/report-pocpoc"
	report_taboola "source/internal/repo/report-taboola"
	report_tiktok "source/internal/repo/report-tiktok"
	report_tonic "source/internal/repo/report-tonic"
	revenue_share "source/internal/repo/revenue-share"
	"source/internal/repo/section"
	token_adsense "source/internal/repo/token-adsense"
	"source/internal/repo/user"
)

type Deps struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

type Repositories struct {
	History                 history.RepoHistory
	User                    user.RepoUser
	Inventory               inventory.RepoInventory
	InventoryAdTag          inventory_ad_tag.RepoInventoryAdTag
	Quiz                    quiz.RepoQuiz
	Category                category.RepoCategory
	AdSchedule              adschedule.RepoAdSchedule
	Gam                     gam.RepoGam
	GamNetwork              gamNetwork.RepoGamNetwork
	AdType                  adType.RepoAdType
	Bidder                  bidder.RepoBidder
	KeyValue                keyValue.RepoKeyValue
	KeyValueGam             keyValueGAM.RepoKeyValueGam
	CronJob                 cronjob.RepoCronJob
	GoogleAdsAPI            googleAdsApi.RepoGoogleAdsAPI
	ReportBodis             reportbodis.RepoReportBodis
	ReportBodisTraffic      reportBodisTraffic.RepoReportBodisTraffic
	ReportOpenMailSubID     reportOpenMailSubID.RepoReportOpenMailSubID
	ReportOutBrain          reportOutBrain.RepoReportOutBrain
	ReportAff               reportAff.RepoReportAff
	ReportAffSelection      reportAffSelection.RepoReportAffSelection
	ReportAffPixel          report_aff_pixel.RepoReportAffPixel
	ReportTonic             report_tonic.RepoReportTonic
	ReportAdsense           report_adsense.RepoReportAdsense
	ReportTaboola           report_taboola.RepoReportTaboola
	AffBlockSection         aff_block_section.RepoAffBlockSection
	Campaign                campaign.RepoCampaign
	Country                 country.RepoCountry
	Payment                 payment.RepoPayment
	Billing                 billing.RepoUserBilling
	RevenueShare            revenue_share.RepoRevenueShare
	Creative                creative.RepoCreative
	Account                 account.RepoAccount
	AdsenseChannel          adsense_channel.RepoAdsenseChannel
	Section                 section.RepoSection
	ReportMgid              report_mgid.RepoReportMgid
	ReportPocPoc            report_pocpoc.RepoReportPocPoc
	Ads                     ads.RepoAds
	TokenAdsense            token_adsense.RepoTokenAdsense
	ReportCodeFuel          report_codefuel.RepoReportCodeFuel
	ReportTikTok            report_tiktok.RepoReportTikTok
	ReportFacebook          report_facebook.RepoReportFacebook
	CampaignTrafficSourceID campaign_traffic_source_id.RepoCampaignTrafficSourceID
	ReportDomainActive      report_domain_active.RepoReportDomainActive
	AdsenseTagMapping       adsense_tag_mapping.RepoAdsenseTagMapping
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		AdSchedule:              adschedule.NewAdScheduleRepo(deps.Db, deps.Cache),
		Quiz:                    quiz.NewQuizRepo(deps.Db, deps.Cache),
		Inventory:               inventory.NewInventoryRepo(deps.Db, deps.Cache),
		InventoryAdTag:          inventory_ad_tag.NewInventoryAdTagRepo(deps.Db, deps.Cache),
		User:                    user.NewUserRepo(deps.Db, deps.Cache),
		History:                 history.NewHistoryRepo(deps.Db, deps.Cache),
		Category:                category.NewCategoryRepo(deps.Db, deps.Cache),
		Gam:                     gam.NewGamRepo(deps.Db, deps.Cache),
		GamNetwork:              gamNetwork.NewGamNetworkRepo(deps.Db, deps.Cache),
		AdType:                  adType.NewAdTypeRepo(deps.Db),
		Bidder:                  bidder.NewBidderRepo(deps.Db),
		KeyValue:                keyValue.NewKeyValueRepo(deps.Db),
		KeyValueGam:             keyValueGAM.NewKeyValueRepo(deps.Db),
		CronJob:                 cronjob.NewCronJobRepo(deps.Db),
		GoogleAdsAPI:            googleAdsApi.NewGoogleAdsAPIRepo(deps.Db),
		ReportBodis:             reportbodis.NewReportBodisRepo(deps.Db),
		ReportBodisTraffic:      reportBodisTraffic.NewReportBodisTrafficRepo(deps.Db, deps.Kafka),
		ReportOpenMailSubID:     reportOpenMailSubID.NewReportOpenMailSubIDRepo(deps.Db, deps.Kafka),
		ReportOutBrain:          reportOutBrain.NewReportOutBrainRepo(deps.Db, deps.Kafka),
		ReportAff:               reportAff.NewReportAffRepo(deps.Db, deps.Cache, deps.Kafka),
		ReportAffSelection:      reportAffSelection.NewReportAffSelectionRepo(deps.Db, deps.Cache, deps.Kafka),
		ReportAffPixel:          report_aff_pixel.NewReportAffPixelRepo(deps.Db, deps.Cache, deps.Kafka),
		ReportTonic:             report_tonic.NewReportTonic(deps.Db),
		ReportAdsense:           report_adsense.NewReportAdsenseRepo(deps.Db),
		ReportTaboola:           report_taboola.NewReportTaboolaRepo(deps.Db),
		AffBlockSection:         aff_block_section.NewAffBlockSectionRepo(deps.Db),
		Campaign:                campaign.NewCampaignRepo(deps.Db, deps.Cache),
		Billing:                 billing.NewUserBillingRepo(deps.Db),
		RevenueShare:            revenue_share.NewRevenueShareRepo(deps.Db),
		Payment:                 payment.NewPaymentRepo(deps.Db),
		Country:                 country.NewCountryRepo(deps.Db),
		Creative:                creative.NewCreativeRepo(deps.Db),
		Account:                 account.NewAccountRepo(deps.Db),
		AdsenseChannel:          adsense_channel.NewAdsenseChannelRP(deps.Db),
		Section:                 section.NewSectionRP(deps.Db),
		Ads:                     ads.NewAdsRepo(),
		ReportMgid:              report_mgid.NewReportMgidRepo(deps.Db),
		ReportPocPoc:            report_pocpoc.NewReportPocPocRepo(deps.Db),
		TokenAdsense:            token_adsense.NewTokenAdsenseRP(deps.Db),
		ReportCodeFuel:          report_codefuel.NewReportCodeFuel(deps.Db),
		ReportTikTok:            report_tiktok.NewReportTikTokRP(deps.Db),
		ReportFacebook:          report_facebook.NewReportFacebookRP(deps.Db),
		CampaignTrafficSourceID: campaign_traffic_source_id.NewCampaignTrafficSourceIDRP(deps.Db),
		ReportDomainActive:      report_domain_active.NewReportDomainActiveRP(deps.Db),
		AdsenseTagMapping:       adsense_tag_mapping.NewAdsenseTagMappingRepo(deps.Db),
	}
}
