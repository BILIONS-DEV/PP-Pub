package model

import (
	"strconv"
	"time"
)

type ResponsePHPAdsense struct {
	Status  bool        `json:"status"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AdsenseAccount struct {
	CreateTime   time.Time       `json:"createTime"`
	DisplayName  string          `json:"displayName"`
	Name         string          `json:"name"`
	PendingTasks interface{}     `json:"pendingTasks"`
	Premium      interface{}     `json:"premium"`
	State        string          `json:"state"`
	TimeZone     AdsenseTimeZone `json:"timeZone"`
}

type AdsenseTimeZone struct {
	Id      string      `json:"id"`
	Version interface{} `json:"version"`
}

type AdsenseReport struct {
	Date              string `json:"date"`
	CustomChannelID   string `json:"custom_channel_id"`
	CustomChannelName string `json:"custom_channel_name"`
	CountryCode       string `json:"country_code"`
	DomainCode        string `json:"domain_code"`
	PageViews         string `json:"page_views"`
	Impressions       string `json:"impressions"`
	Clicks            string `json:"clicks"`
	AdRequestsCtr     string `json:"ad_requests_ctr"`
	CostPerClick      string `json:"cost_per_click"`
	AdRequestsRpm     string `json:"ad_requests_rpm"`
	EstimatedEarnings string `json:"estimated_earnings"`
}

func (t *AdsenseReport) ToModel() (record *ReportAdsenseModel) {
	pageView, _ := strconv.ParseInt(t.PageViews, 10, 64)
	impressions, _ := strconv.ParseInt(t.Impressions, 10, 64)
	clicks, _ := strconv.ParseInt(t.Clicks, 10, 64)
	adRequestCtr, _ := strconv.ParseFloat(t.AdRequestsCtr, 64)
	costPerClick, _ := strconv.ParseFloat(t.CostPerClick, 64)
	adRequestsRpm, _ := strconv.ParseFloat(t.AdRequestsRpm, 64)
	estimatedEarnings, _ := strconv.ParseFloat(t.EstimatedEarnings, 64)
	return &ReportAdsenseModel{
		Date:              t.Date,
		CustomChannelID:   t.CustomChannelID,
		CustomChannelName: t.CustomChannelName,
		CountryCode:       t.CountryCode,
		DomainCode:        t.DomainCode,
		PageViews:         pageView,
		Impressions:       impressions,
		Clicks:            clicks,
		AdRequestsCtr:     adRequestCtr,
		CostPerClick:      costPerClick,
		AdRequestsRpm:     adRequestsRpm,
		EstimatedEarnings: estimatedEarnings,
	}
}

type ReportAdsenseModel struct {
	ID                int64             `gorm:"column:id;primaryKey;autoIncrement"`
	Type              TYPEReportAdsense `gorm:"column:type"`
	CampaignID        int64             `gorm:"column:campaign_id"`
	SectionID         string            `gorm:"column:section_id"`
	Account           string            `gorm:"column:account"`
	Date              string            `gorm:"column:date"`
	CustomChannelID   string            `gorm:"column:custom_channel_id"`
	CustomChannelName string            `gorm:"column:custom_channel_name"`
	CountryCode       string            `gorm:"column:country_code"`
	DomainCode        string            `gorm:"column:domain_code"`
	PageViews         int64             `gorm:"column:page_views"`
	Impressions       int64             `gorm:"column:impressions"`
	Clicks            int64             `gorm:"column:clicks"`
	AdRequestsCtr     float64           `gorm:"column:ad_requests_ctr"`
	CostPerClick      float64           `gorm:"column:cost_per_click"`
	AdRequestsRpm     float64           `gorm:"column:ad_requests_rpm"`
	EstimatedEarnings float64           `gorm:"column:estimated_earnings"`
}

type TYPEReportAdsense string

const (
	TYPEReportAdsenseSubDomain TYPEReportAdsense = "sub_domain"
	TYPEReportAdsenseChannel   TYPEReportAdsense = "channel"
)

func (ReportAdsenseModel) TableName() string {
	return "report_adsense"
}
func (t *ReportAdsenseModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}
