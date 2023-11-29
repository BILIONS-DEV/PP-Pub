package report_adsense

import (
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/repo"
	report_adsense "source/internal/repo/report-adsense"
	"source/pkg/logger"
	"strconv"
	"strings"
	"time"
)

type UseCaseReportAdsense interface {
	HandlerReportAdsenseChannel(refreshToken string, accountAdsense string, startDate string, endDate string) (err error)
	HandlerReportAdsenseSubDomain(refreshToken string, accountAdsense string, startDate string, endDate string) (err error)
	Filter(input *InputFilter) (totalRecords int64, records []*model.ReportAffModel, recordTotal *model.ReportAffModel, err error)
	GetAllAdsense() (records []*model.TokenAdsense, err error)
}

func NewReportAdsenseUC(repos *repo.Repositories) *adsenseUC {
	return &adsenseUC{
		repos: repos,
	}
}

type adsenseUC struct {
	repos *repo.Repositories
}

func (t *adsenseUC) HandlerReportAdsenseChannel(refreshToken string, accountAdsense string, startDate string, endDate string) (err error) {
	// get accounts adsense
	accounts, err := t.repos.ReportAdsense.GetAllAccounts(refreshToken, accountAdsense)
	if err != nil {
		return
	}
	for _, account := range accounts {
		reports, err := t.repos.ReportAdsense.GetReports(report_adsense.InputGetReports{
			AccountAdsense: accountAdsense,
			RefreshToken:   refreshToken,
			Account:        account.Name,
			StartDate:      startDate,
			EndDate:        endDate,
		})
		if err != nil {
			continue
		}

		for _, report := range reports {
			//fmt.Printf("%+v \n", report)
			query := make(map[string]interface{})
			query["date"] = report.Date
			query["custom_channel_id"] = report.CustomChannelID
			query["custom_channel_name"] = report.CustomChannelName
			query["country_code"] = report.CountryCode
			reportAdsense, _ := t.repos.ReportAdsense.FindOneByQuery(query)
			if reportAdsense.IsFound() {
				record := report.ToModel()
				record.Account = accountAdsense
				record.ID = reportAdsense.ID
				// Channel ID sẽ có dạng partner-pub-3275351677467683:5842325056 => Phía cuối sau : chính là channelID
				if reportAdsense.CampaignID == 0 || report.Date == time.Now().UTC().Format("2006-01-02") {
					var campaignID int64
					splitChannelID := strings.Split(report.CustomChannelID, ":")
					if len(splitChannelID) > 0 {
						channelID := splitChannelID[1]
						adsenseChanel, _ := t.repos.AdsenseChannel.FindChannelCountry(channelID, report.CountryCode)
						campaignID = adsenseChanel.CampaignID
					}
					record.CampaignID = campaignID
				} else {
					record.CampaignID = reportAdsense.CampaignID
				}
				record.Type = model.TYPEReportAdsenseChannel
				_ = t.repos.ReportAdsense.Save(record)
			} else {
				record := report.ToModel()
				var campaignID int64
				// Channel ID sẽ có dạng partner-pub-3275351677467683:5842325056 => Phía cuối sau : chính là channelID
				splitChannelID := strings.Split(report.CustomChannelID, ":")
				if len(splitChannelID) > 0 {
					channelID := splitChannelID[1]
					adsenseChanel, _ := t.repos.AdsenseChannel.FindChannelCountry(channelID, report.CountryCode)
					campaignID = adsenseChanel.CampaignID
				}
				record.CampaignID = campaignID
				record.Account = accountAdsense
				record.Type = model.TYPEReportAdsenseChannel
				_ = t.repos.ReportAdsense.Save(record)
			}
		}
	}
	return
}

func (t *adsenseUC) HandlerReportAdsenseSubDomain(refreshToken string, accountAdsense string, startDate string, endDate string) (err error) {
	// get accounts adsense
	accounts, err := t.repos.ReportAdsense.GetAllAccounts(refreshToken, accountAdsense)
	if err != nil {
		return
	}
	var records []*model.ReportAdsenseModel
	for _, account := range accounts {
		reports, err := t.repos.ReportAdsense.GetReports2(report_adsense.InputGetReports{
			AccountAdsense: accountAdsense,
			RefreshToken:   refreshToken,
			Account:        account.Name,
			StartDate:      startDate,
			EndDate:        endDate,
		})
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		for _, report := range reports {
			splitHost := strings.Split(report.DomainCode, ".")
			if len(splitHost) < 3 {
				continue
			}
			// Xử lý subDomain
			var redirectID int64
			var sectionID string
			subDomain := splitHost[0]
			if strings.Contains(subDomain, "_") {
				splitSubDomain := strings.Split(subDomain, "_")
				if len(splitSubDomain) < 2 {
					continue
				}
				redirectID, _ = strconv.ParseInt(splitSubDomain[0], 10, 64)
				sectionID = splitSubDomain[1]
			} else {
				params, _ := t.repos.Campaign.FindCacheSubDomainParams(subDomain)
				redirectID = params.CampaignID
				sectionID = params.SectionID
			}
			//fmt.Printf("%+v \n", report)
			query := make(map[string]interface{})
			query["date"] = report.Date
			query["domain_code"] = report.DomainCode
			reportAdsense, _ := t.repos.ReportAdsense.FindOneByQuery(query)
			if reportAdsense.IsFound() {
				record := report.ToModel()
				record.Account = accountAdsense
				record.ID = reportAdsense.ID
				record.CampaignID = redirectID
				record.SectionID = sectionID
				record.Type = model.TYPEReportAdsenseSubDomain
				records = append(records, record)
			} else {
				record := report.ToModel()
				record.Account = accountAdsense
				record.CampaignID = redirectID
				record.SectionID = sectionID
				record.Type = model.TYPEReportAdsenseSubDomain
				records = append(records, record)
			}
		}
	}
	err = t.repos.ReportAdsense.SaveSlice(records)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

type InputFilter struct {
	datatable.Request
	UserID         int64
	QuerySearch    string
	StartDate      string
	EndDate        string
	Campaigns      interface{}
	TrafficSources interface{}
	PublisherID    interface{}
	GroupBy        interface{}
	SectionID      interface{}
}

func (t *adsenseUC) Filter(input *InputFilter) (totalRecords int64, records []*model.ReportAffModel, recordTotal *model.ReportAffModel, err error) {
	//=> khởi tạo input của Repo
	inputRepo := report_adsense.InputFilter{
		UserID:         input.UserID,
		Search:         input.QuerySearch,
		Offset:         input.Start,
		Limit:          input.Length,
		Order:          input.OrderString(),
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Campaigns:      input.Campaigns,
		TrafficSources: input.TrafficSources,
		SectionID:      input.SectionID,
		GroupBy:        input.GroupBy,
	}
	return t.repos.ReportAdsense.Filter(&inputRepo)
}

func (t *adsenseUC) GetAllAdsense() (records []*model.TokenAdsense, err error) {
	return t.repos.TokenAdsense.FindAll()
}
