package usecase

import (
	"source/internal/lang"
	"source/internal/repo"
	adschedule "source/internal/usecase/ad-schedule"
	adType "source/internal/usecase/ad-type"
	"source/internal/usecase/ads"
	adsense_tag_mapping "source/internal/usecase/adsense-tag-mapping"
	"source/internal/usecase/api"
	"source/internal/usecase/bidder"
	"source/internal/usecase/campaign"
	"source/internal/usecase/category"
	"source/internal/usecase/country"
	"source/internal/usecase/creative"
	"source/internal/usecase/cronjob"
	gamNetwork "source/internal/usecase/gam-network"
	"source/internal/usecase/history"
	"source/internal/usecase/history2"
	"source/internal/usecase/inventory"
	keyValue "source/internal/usecase/key-value"
	"source/internal/usecase/payment"
	"source/internal/usecase/quiz"
	report_adsense "source/internal/usecase/report-adsense"
	reportAff "source/internal/usecase/report-aff"
	reportAffSelection "source/internal/usecase/report-aff-selection"
	report_bodis "source/internal/usecase/report-bodis"
	report_codefuel "source/internal/usecase/report-codefuel"
	report_mgid "source/internal/usecase/report-mgid"
	reportOpenMail "source/internal/usecase/report-openmail"
	reportOutBrain "source/internal/usecase/report-outbrain"
	report_pocpoc "source/internal/usecase/report-pocpoc"
	report_taboola "source/internal/usecase/report-taboola"
	report_tonic "source/internal/usecase/report-tonic"
	"source/internal/usecase/user"
)

type Deps struct {
	Repos       *repo.Repositories
	Translation *lang.Translation
}

type UseCases struct {
	History2           history2.UsecaseHistory
	History            history.UsecaseHistory
	User               user.UsecaseUser
	Inventory          inventory.UsecaseInventory
	Quiz               quiz.UsecaseQuiz
	Category           category.UsecaseCategory
	AdSchedule         adschedule.UsecaseAdSchedule
	GamNetwork         gamNetwork.UsecaseGamNetwork
	AdType             adType.UsecaseAdType
	Bidder             bidder.UsecaseBidder
	KeyValue           keyValue.UsecaseKeyValue
	CronJob            cronjob.UsecaseCronjob
	ReportBodis        report_bodis.UsecaseReportBodis
	Campaign           campaign.UsecaseCampaign
	ReportOpenMail     reportOpenMail.UseCaseReportOpenMail
	ReportOutBrain     reportOutBrain.UseCaseReportOutBrain
	ReportAff          reportAff.UsecaseReportAff
	ReportAffSelection reportAffSelection.UsecaseReportAffSelection
	ReportTonic        report_tonic.UseCaseReportTonic
	ReportTaboola      report_taboola.UseCaseReportTaboola
	ReportAdsense      report_adsense.UseCaseReportAdsense
	API                api.UsecaseAPI
	Payment            payment.UsercasePayment
	Country            country.UsecaseCountry
	Creative           creative.UsecaseCreative
	ReportMgid         report_mgid.UseCaseReportMgid
	ReportPocPoc       report_pocpoc.UseCaseReportPocPoc
	Ads                ads.UsecaseAds
	ReportCodeFuel     report_codefuel.UseCaseReportCodeFuel
	AdsenseTagMapping  adsense_tag_mapping.UsecaseAdsenseTagMapping
}

func NewUseCases(deps *Deps) *UseCases {
	return &UseCases{
		History2:           history2.NewHistoryUsecase(deps.Repos, deps.Translation),
		History:            history.NewHistoryUsecase(deps.Repos, deps.Translation),
		AdSchedule:         adschedule.NewAdScheduleUC(deps.Repos),
		Quiz:               quiz.NewQuizUsecase(deps.Repos, deps.Translation),
		Inventory:          inventory.NewInventoryUC(deps.Repos),
		User:               user.NewUserFeUsecase(deps.Repos, deps.Translation),
		Category:           category.NewCategoryUsecase(deps.Repos, deps.Translation),
		GamNetwork:         gamNetwork.NewGamNetworkUC(deps.Repos),
		AdType:             adType.NewAdTypeUC(deps.Repos),
		Bidder:             bidder.NewBidderUC(deps.Repos),
		KeyValue:           keyValue.NewKeyValueUC(deps.Repos),
		CronJob:            cronjob.NewCronJobUC(deps.Repos),
		ReportBodis:        report_bodis.NewReportBodisUsecase(deps.Repos, deps.Translation),
		Campaign:           campaign.NewCampaignUsecase(deps.Repos, deps.Translation),
		ReportOpenMail:     reportOpenMail.NewReportOpenMailUC(deps.Repos),
		ReportOutBrain:     reportOutBrain.NewReportOutBrainUC(deps.Repos),
		ReportAff:          reportAff.NewReportAffUsecase(deps.Repos),
		ReportAffSelection: reportAffSelection.NewReportAffSelectionUsecase(deps.Repos),
		ReportTonic:        report_tonic.NewReportTonicUC(deps.Repos),
		ReportAdsense:      report_adsense.NewReportAdsenseUC(deps.Repos),
		ReportTaboola:      report_taboola.NewReportTaboolaUC(deps.Repos),
		API:                api.NewApiUsecase(deps.Repos),
		Payment:            payment.NewPaymentUsecase(deps.Repos, deps.Translation),
		Country:            country.NewCountryUC(deps.Repos),
		Creative:           creative.NewCreativeUsecase(deps.Repos),
		ReportMgid:         report_mgid.NewReportMgidUC(deps.Repos),
		ReportPocPoc:       report_pocpoc.NewReportPocPocUC(deps.Repos),
		Ads:                ads.NewAdsUsecase(deps.Repos),
		ReportCodeFuel:     report_codefuel.NewReportCodeFuelUC(deps.Repos),
		AdsenseTagMapping:  adsense_tag_mapping.NewAdsenseTagMappingUC(deps.Repos),
	}
}
