package model

import "strconv"

func (*ReportAffModel) TableName() string {
	return "report_aff"
}

type ReportAffForBlockSectionModel struct {
	ReportAffModel
	StatusBlock            TYPEState
	ExpectedCPA            float64
	RealtimeCPA            float64
	RuleBlock              string
	LevelBlock             string
	TrafficSourceAccountID string
}

type ReportAffModel struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement" `
	UserID        int64  `gorm:"column:user_id"`
	Date          string `gorm:"column:date"`
	TrafficSource string `gorm:"column:traffic_source"`
	SectionId     string `gorm:"column:section_id"`
	SectionName   string `gorm:"column:section_name"`
	CampaignID    string `gorm:"column:campaign_id"`
	CampaignName  string `gorm:"column:campaign_name"`
	RedirectID    int64  `gorm:"column:redirect_id"`
	StyleID       string `gorm:"column:style_id"`
	LayoutID      int64  `gorm:"column:layout_id"`
	LayoutVersion string `gorm:"column:layout_version"`
	Device        string `gorm:"column:device"`
	Geo           string `gorm:"column:geo"`
	Partner       string `gorm:"column:partner"`

	ImpQuantumdex int64 `gorm:"column:imp_quantumdex"`
	Impressions   int64 `gorm:"column:impressions"`
	Click         int64 `gorm:"column:click"`
	Click4        int64 `gorm:"column:click4"`
	SystemTraffic int64 `gorm:"column:system_traffic"`
	PreBotTraffic int64 `gorm:"column:pre_bot_traffic"`
	BotTraffic    int64 `gorm:"column:bot_traffic"`

	Cost                float64 `gorm:"column:cost"`
	CostCPC             float64 `gorm:"column:cost_cpc"`
	Revenue             float64 `gorm:"column:revenue"`
	EstimatedRevenue    float64 `gorm:"column:estimated_revenue"`
	PreEstimatedRevenue float64 `gorm:"column:pre_estimated_revenue"`

	SupplyClick       int64 `gorm:"column:supply_click"`
	SupplyConversions int64 `gorm:"column:supply_conversions"`
	ClickAdsense      int64 `gorm:"column:click_adsense"`

	CTR    float64 `gorm:"<-:false"` //=> SELECT SUM(click)/SUM(impression)*100  as ctr
	CPC    float64 `gorm:"<-:false"` //=> SELECT SUM(revenue)/SUM(click) as cpc
	CPA    float64 `gorm:"<-:false"` //=> SELECT SUM(cost)/SUM(click) as cpa
	Profit float64 `gorm:"<-:false"` //=> SELECT SUM(revenue)-SUM(cost) as profit
	Roi    float64 `gorm:"<-:false"` //=> SELECT (((SUM(revenue)-SUM(cost))/SUM(cost))*100) AS roi
}

func (t *ReportAffModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}

func (*ReportAffPixelModel) TableName() string {
	return "report_aff_pixel"
}

type ReportAffPixelModel struct {
	ID                 int64   `gorm:"column:id;primaryKey;autoIncrement" `
	UID                string  `gorm:"column:uid"`
	UserID             int64   `gorm:"column:user_id"`
	TimeConversion     string  `gorm:"column:time_conversion"`
	TrafficSource      string  `gorm:"column:traffic_source"`
	DemandSource       string  `gorm:"column:demand_source"`
	CampaignID         string  `gorm:"column:campaign_id"`
	SectionId          string  `gorm:"column:section_id"`
	SectionName        string  `gorm:"column:section_name"`
	PublisherID        string  `gorm:"column:publisher_id"`
	PublisherName      string  `gorm:"column:publisher_name"`
	AdID               string  `gorm:"column:ad_id"`
	RedirectID         string  `gorm:"column:redirect_id"`
	ImpQuantumdex      int64   `gorm:"column:imp_quantumdex"`
	SystemTraffic      int64   `gorm:"column:system_traffic"`
	PreBotTraffic      int64   `gorm:"column:pre_bot_traffic"`
	BotTraffic         int64   `gorm:"column:bot_traffic"`
	StyleID            string  `gorm:"column:style_id"`
	Impressions        int64   `gorm:"column:impressions"`
	Click              int64   `gorm:"column:click"`
	Click2             int64   `gorm:"column:click2"`
	Click3             int64   `gorm:"column:click3"`
	Click4             int64   `gorm:"column:click4"`
	PreEstimateRevenue float64 `gorm:"column:pre_estimate_revenue"`
	EstimateRevenue    float64 `gorm:"column:estimate_revenue"`
	Cost               float64 `gorm:"column:cost"`
	Account            string  `gorm:"column:account"`
	CPC                string  `gorm:"column:cpc"`
	Referrer           string  `gorm:"column:referrer"`
	LayoutID           int64   `gorm:"column:layout_id"`
	LayoutVersion      string  `gorm:"column:layout_version"`
	Device             string  `gorm:"column:device"`
	Geo                string  `gorm:"column:geo"`
	Key                string  `gorm:"-"`
}

func (t *ReportAffPixelModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}

func (PixelAff) SetName() string {
	return "pixel_aff"
}

type PixelAff struct {
	Key           string `as:"-"`
	UserID        int64  `as:"user_id"`
	TimeUTC       string `as:"time_utc"`
	TrafficSource string `as:"traffic_source"`
	DemandSource  string `as:"demand_source"`
	CampaignID    string `as:"campaign_id"`
	SectionId     string `as:"section_id"`
	SectionName   string `as:"section_name"`
	PublisherID   string `as:"publisher_id"`
	PublisherName string `as:"publisher_name"`
	AdID          string `as:"ad_id"`
	RedirectID    string `as:"redirect_id"`
	StyleID       string `as:"style_id"`
	CPC           string `as:"cpc"`
	Account       string `as:"account"`

	ImpQuantumdex      int64   `as:"imp_quantumdex"`
	SystemTraffic      int64   `as:"system_traffic"`
	PreBotTraffic      int64   `as:"pre_bot_traffic"`
	BotTraffic         int64   `as:"bot_traffic"`
	Impressions        int64   `as:"impressions"`
	Click              int64   `as:"click"`
	Click2             int64   `as:"click2"` // Dùng cho google
	Click3             int64   `as:"click3"` // Dùng cho google
	Click4             int64   `as:"click4"` // Dùng cho google
	PreEstimateRevenue float64 `as:"pre_revenue"`
	EstimateRevenue    float64 `as:"revenue"`

	UID            string `as:"uid"`
	TimeConversion string `as:"time_conversion"`
	Referrer       string `as:"referrer"`

	Uuid []string `as:"-"` // Xử lý lưu list uuid key aerospike
}

func (PixelAffCheck) SetName() string {
	return "pixel_aff_check"
}

type PixelAffCheck struct {
	Key            string `as:"-" json:"key"`
	UserID         int64  `as:"user_id" json:"userID"`
	TimeConversion string `as:"time_conversion" json:"timeConversion"`
	TrafficSource  string `as:"traffic_source" json:"trafficSource"`
	DemandSource   string `as:"demand_source" json:"demandSource"`
	CampaignID     string `as:"campaign_id" json:"campaignID"`
	SectionId      string `as:"section_id" json:"sectionId"`
	SectionName    string `as:"section_name" json:"sectionName"`
	PublisherID    string `as:"publisher_id" json:"publisherID"`
	PublisherName  string `as:"publisher_name" json:"publisherName"`
	AdID           string `as:"ad_id" json:"adID"`
	RedirectID     string `as:"redirect_id" json:"redirectID"`
	CPC            string `as:"cpc" json:"CPC"`
	Account        string `as:"account" json:"account"`

	ImpQuantumdex      int64   `as:"imp_quantumdex" json:"impQuantumdex"`
	SystemTraffic      int64   `as:"system_traffic" json:"systemTraffic"`
	PreBotTraffic      int64   `as:"pre_bot_traffic" json:"preBotTraffic"`
	BotTraffic         int64   `as:"bot_traffic" json:"botTraffic"`
	Impressions        int64   `as:"impressions" json:"impressions"`
	Click              int64   `as:"click" json:"click"`
	Click2             int64   `as:"click2" json:"click2"`
	Click3             int64   `as:"click3" json:"click3"`
	Click4             int64   `as:"click4" json:"click4"`
	PreEstimateRevenue float64 `as:"pre_revenue" json:"preEstimateRevenue"`
	EstimateRevenue    float64 `as:"revenue" json:"estimateRevenue"`
	Referrer           string  `as:"referrer" json:"referrer"`

	ClickID string `as:"click_id" json:"clickID"`
}

type PixelAffForAll struct {
	Key            string `json:"key"`
	UserID         int64  `json:"user_id"`
	UID            string `json:"uid"`
	TimeConversion string `json:"time_conversion"`
	TimeUTC        string `json:"time_utc"`
	TrafficSource  string `json:"traffic_source"`
	DemandSource   string `json:"demand_source"`
	CampaignID     string `json:"campaign_id"`
	SectionId      string `json:"section_id"`
	SectionName    string `json:"section_name"`
	PublisherID    string `json:"publisher_id"`
	PublisherName  string `json:"publisher_name"`
	AdID           string `json:"ad_id"`
	RedirectID     string `json:"redirect_id"`
	StyleID        string `json:"style_id"`
	CPC            string `json:"cpc"`
	Account        string `json:"account"`

	ImpQuantumdex      int64   `json:"imp_quantumdex"`
	SystemTraffic      int64   `json:"system_traffic"`
	PreBotTraffic      int64   `json:"pre_bot_traffic"`
	BotTraffic         int64   `json:"bot_traffic"`
	Impressions        int64   `json:"impressions"`
	Click              int64   `json:"click"`
	Click2             int64   `json:"click2"`
	Click3             int64   `json:"click3"`
	Click4             int64   `json:"click4"`
	PreEstimateRevenue float64 `json:"pre_revenue"`
	EstimateRevenue    float64 `json:"revenue"`
	Referrer           string  `json:"referrer"`
	LayoutID           string  `as:"layoutID" json:"layoutID"`
	LayoutVersion      string  `as:"layoutVersion" json:"layoutVersion"`
	Device             string  `as:"device" json:"device"`
	Geo                string  `as:"geo" json:"geo"`

	Uuid []string `json:"-"` // Xử lý lưu list uuid key aerospike
}

func (t *PixelAffForAll) MakeModelReportPixel() ReportAffPixelModel {
	layoutID, _ := strconv.ParseInt(t.LayoutID, 10, 64)
	return ReportAffPixelModel{
		UID:                t.UID,
		UserID:             t.UserID,
		TimeConversion:     t.TimeConversion,
		TrafficSource:      t.TrafficSource,
		DemandSource:       t.DemandSource,
		CampaignID:         t.CampaignID,
		SectionId:          t.SectionId,
		SectionName:        t.SectionName,
		PublisherID:        t.PublisherID,
		PublisherName:      t.PublisherName,
		AdID:               t.AdID,
		RedirectID:         t.RedirectID,
		StyleID:            t.StyleID,
		ImpQuantumdex:      t.ImpQuantumdex,
		SystemTraffic:      t.SystemTraffic,
		PreBotTraffic:      t.PreBotTraffic,
		BotTraffic:         t.BotTraffic,
		Impressions:        t.Impressions,
		Click:              t.Click,
		Click2:             t.Click2,
		Click3:             t.Click3,
		Click4:             t.Click4,
		PreEstimateRevenue: t.PreEstimateRevenue,
		EstimateRevenue:    t.EstimateRevenue,
		Account:            t.Account,
		CPC:                t.CPC,
		Key:                t.Key,
		Referrer:           t.Referrer,
		LayoutID:           layoutID,
		LayoutVersion:      t.LayoutVersion,
		Device:             t.Device,
		Geo:                t.Geo,
	}
}
