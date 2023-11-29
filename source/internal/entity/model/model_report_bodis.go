package model

import (
	"gorm.io/gorm"
	"time"
)

func (ReportBodisModel) TableName() string {
	return "report_bodis"
}

type ReportBodisModel struct {
	ID                         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Subid                      string    `gorm:"column:subid" json:"subid"`
	Visits                     int       `gorm:"column:visits" json:"visits"`
	LandingPageVisits          int       `gorm:"column:landing_page_visits" json:"landing_page_visits"`
	Zeroclick_visits           int       `gorm:"column:zeroclick_visits" json:"zeroclick_visits"`
	Clicks                     int       `gorm:"column:clicks" json:"clicks"`
	CreditedRevenue            float64   `gorm:"column:credited_revenue" json:"credited_revenue"`
	LandingPageCreditedRevenue float64   `gorm:"column:landing_page_credited_revenue" json:"landing_page_credited_revenue"`
	ZeroclickCreditedRevenue   float64   `gorm:"column:zeroclick_credited_revenue" json:"zeroclick_credited_revenue"`
	Ctr                        float64   `gorm:"column:ctr" json:"ctr"`
	Epc                        float64   `gorm:"column:epc" json:"epc"`
	Rpm                        float64   `gorm:"column:rpm" json:"rpm"`
	LandingPageRpm             float64   `gorm:"column:landing_page_rpm" json:"landing_page_rpm"`
	ZeroclickRpm               float64   `gorm:"column:zeroclick_rpm" json:"zeroclick_rpm"`
	ClicksSpamRatio            float64   `gorm:"column:clicks_spam_ratio" json:"clicks_spam_ratio"`
	IsFinalized                float64   `gorm:"column:is_finalized" json:"is_finalized"`
	Date                       time.Time `gorm:"column:date" json:"date"`
}

type ResponseData struct {
	Subid                      interface{} `gorm:"column:subid" json:"subid"`
	Visits                     int         `gorm:"column:visits" json:"visits"`
	LandingPageVisits          int         `gorm:"column:landing_page_visits" json:"landing_page_visits"`
	Zeroclick_visits           int         `gorm:"column:zeroclick_visits" json:"zeroclick_visits"`
	Clicks                     int         `gorm:"column:clicks" json:"clicks"`
	CreditedRevenue            float64     `gorm:"column:credited_revenue" json:"credited_revenue"`
	LandingPageCreditedRevenue float64     `gorm:"column:landing_page_credited_revenue" json:"landing_page_credited_revenue"`
	ZeroclickCreditedRevenue   float64     `gorm:"column:zeroclick_credited_revenue" json:"zeroclick_credited_revenue"`
	Ctr                        float64     `gorm:"column:ctr" json:"ctr"`
	Epc                        float64     `gorm:"column:epc" json:"epc"`
	Rpm                        float64     `gorm:"column:rpm" json:"rpm"`
	LandingPageRpm             float64     `gorm:"column:landing_page_rpm" json:"landing_page_rpm"`
	ZeroclickRpm               float64     `gorm:"column:zeroclick_rpm" json:"zeroclick_rpm"`
	ClicksSpamRatio            float64     `gorm:"column:clicks_spam_ratio" json:"clicks_spam_ratio"`
	IsFinalized                float64     `gorm:"column:is_finalized" json:"is_finalized"`
}

func (ReportBodisTrafficModel) TableName() string {
	return "report_bodis_traffic"
}

type ReportBodisTrafficModel struct {
	ID               int64   `gorm:"column:id;primaryKey;autoIncrement"`
	Time             string  `gorm:"column:time"`
	VisitID          string  `gorm:"column:visit_id"`
	DomainName       string  `gorm:"column:domain_name"`
	IpAddress        string  `gorm:"column:ip_address"`
	Type             string  `gorm:"column:type"`
	CountryID        int64   `gorm:"column:country_id"`
	PageQuery        string  `gorm:"column:page_query"`
	TrafficSource    string  `gorm:"column:traffic_source"`
	Campaign         string  `gorm:"column:campaign"`
	SelectionID      string  `gorm:"column:selection_id"`
	Clicks           int64   `gorm:"column:clicks"`
	EstimatedRevenue float64 `gorm:"column:estimated_revenue"`
	gorm.Model
}

func (t *ReportBodisTrafficModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}

type TrackingAdsMessage struct {
	TrafficSource string `json:"trafficSource"`
	PubName       string `json:"pubName"`
	PubId         string `json:"pubId"`
	SelectionId   string `json:"selectionId"`
	SectionName   string `json:"sectionName"`
	Campaign      string `json:"campaign"`
	CampaignName  string `json:"campaignName"`
	AdId          string `json:"adId"`
	Partner       string `json:"partner"`
	LandingPage   string `json:"-"`

	Device      string `json:"device"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
	Time        string `json:"time"`

	Impressions int `json:"impressions"`
	Click       int `json:"click"`

	Amount       float64 `json:"amount"`
	GrossRevenue float64 `json:"grossRevenue"`
}

func (t *ReportBodisTrafficModel) ToMessageKafka() (request TrackingAdsMessage) {
	request = TrackingAdsMessage{
		TrafficSource: t.TrafficSource,
		SelectionId:   t.SelectionID,
		Campaign:      t.Campaign,
		Time:          t.Time,
		Click:         int(t.Clicks),
		GrossRevenue:  t.EstimatedRevenue,
	}
	return
}
