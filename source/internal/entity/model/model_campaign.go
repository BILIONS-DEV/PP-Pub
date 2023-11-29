package model

import (
	"errors"
	"gorm.io/gorm"
	"source/pkg/utility"
	"strings"
	"time"
)

func (CampaignModel) TableName() string {
	return "campaign"
}

type CampaignModel struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name            string `gorm:"column:name" json:"name"`
	Status          string `gorm:"column:status" json:"status"`
	DeliveryStatus  string `gorm:"column:delivery_status" json:"delivery_status"`
	TrafficSource   string `gorm:"column:traffic_source" json:"traffic_source"`
	DemandSource    string `gorm:"column:demand_source" json:"demand_source"`
	UserID          int64  `gorm:"column:user_id" json:"user_id"`
	CampaignGroupID int64  `gorm:"column:campaign_group_id" json:"campaign_group_id"`
	CreativeID      int64  `gorm:"column:creative_id" json:"creative_id"`
	// PixelId        string                     `gorm:"column:pixel_id" json:"pixel_id"`
	Vertical        string                     `gorm:"column:vertical" json:"vertical"`
	LandingPages    string                     `gorm:"column:landing_pages" json:"landing_pages"`
	MainKeyword     string                     `gorm:"column:main_keyword" json:"main_keyword"`
	Channel         string                     `gorm:"column:channel" json:"channel"`
	GD              string                     `gorm:"column:gd" json:"gd"`
	UserFlow        string                     `gorm:"column:user_flow" json:"user_flow"`
	Bidding         string                     `gorm:"column:bidding" json:"bidding"`
	Cpc             float64                    `gorm:"column:cpc" json:"cpc"`
	ExpectedCpa     float64                    `gorm:"column:expected_cpa" json:"expected_cpa"`
	Budget          float64                    `gorm:"column:budget" json:"budget"`
	Location        []CampaignLocationModel    `gorm:"foreignKey:CampaignID;references:ID" json:"location"`
	Params          string                     `gorm:"column:params" json:"params"`
	UrlTrackImp     string                     `gorm:"column:url_track_imp" json:"url_track_imp"`
	URLTrackClick   string                     `gorm:"column:url_track_click" json:"url_track_click"`
	Keywords        []CampaignKeywordModel     `gorm:"foreignKey:CampaignID;references:ID;" json:"keywords"`
	LandingPageGeo  []CampaignLandingPageModel `gorm:"foreignKey:CampaignID;references:ID"`
	Device          []CampaignDeviceModel      `gorm:"foreignKey:CampaignID;references:ID" json:"device"`
	Country         []CampaignLocationModel    `gorm:"foreignKey:CampaignID;references:ID" json:"country"`
	Flag            string                     `gorm:"column:flag" json:"flag"`
	Account         string                     `gorm:"column:account" json:"account"`
	AccountAdsense  string                     `gorm:"column:account_adsense" json:"account_adsense"`
	Terms           string                     `gorm:"column:terms" json:"terms"`
	AdTitle         string                     `gorm:"column:ad_title" json:"ad_title"`
	TrafficSourceID string                     `gorm:"column:traffic_source_id" json:"traffic_source_id"`
	CreatedAt       time.Time                  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time                  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt             `gorm:"column:deleted_at" json:"deleted_at"`
	AccountModel    AccountModel               `gorm:"foreignKey:Account;references:Name" json:"account_model"`
}

type CampaignKeywordModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignID int64  `gorm:"column:campaign_id" json:"campaign_id"`
	Keyword    string `gorm:"column:keyword" json:"keyword"`
}

func (CampaignKeywordModel) TableName() string {
	return "campaign_keyword"
}

type CampaignLandingPageModel struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignID  int64  `gorm:"column:campaign_id" json:"campaign_id"`
	Country     string `gorm:"column:country" json:"country"`
	LandingPage string `gorm:"column:landing_page" json:"landing_page"`
}

func (CampaignLandingPageModel) TableName() string {
	return "campaign_landing_page"
}

type CampaignDeviceModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignID int64  `gorm:"column:campaign_id" json:"campaign_id"`
	Deivce     string `gorm:"column:device" json:"device"`
}

func (CampaignDeviceModel) TableName() string {
	return "campaign_device"
}

type CampaignLocationModel struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignID   int64  `gorm:"column:campaign_id" json:"campaign_id"`
	Location     string `gorm:"column:location" json:"location"`
	LocationType string `gorm:"column:location_type" json:"location_type"`
}

func (CampaignLocationModel) TableName() string {
	return "campaign_location"
}

type CampaignCountryModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignID int64  `gorm:"column:campaign_id" json:"campaign_id"`
	Country    string `gorm:"column:country" json:"country"`
}

func (CampaignCountryModel) TableName() string {
	return "campaign_country"
}

// type ResponseData struct {
// 	Subid                      interface{} `gorm:"column:subid" json:"subid"`
// 	Visits                     int         `gorm:"column:visits" json:"visits"`
// 	LandingPageVisits          int         `gorm:"column:landing_page_visits" json:"landing_page_visits"`
// 	Zeroclick_visits           int         `gorm:"column:zeroclick_visits" json:"zeroclick_visits"`
// 	DemandConversion                     int         `gorm:"column:clicks" json:"clicks"`
// 	CreditedRevenue            float64     `gorm:"column:credited_revenue" json:"credited_revenue"`
// 	LandingPageCreditedRevenue float64     `gorm:"column:landing_page_credited_revenue" json:"landing_page_credited_revenue"`
// 	ZeroclickCreditedRevenue   float64     `gorm:"column:zeroclick_credited_revenue" json:"zeroclick_credited_revenue"`
// 	Ctr                        float64     `gorm:"column:ctr" json:"ctr"`
// 	Epc                        float64     `gorm:"column:epc" json:"epc"`
// 	Rpm                        float64     `gorm:"column:rpm" json:"rpm"`
// 	LandingPageRpm             float64     `gorm:"column:landing_page_rpm" json:"landing_page_rpm"`
// 	ZeroclickRpm               float64     `gorm:"column:zeroclick_rpm" json:"zeroclick_rpm"`
// 	ClicksSpamRatio            float64     `gorm:"column:clicks_spam_ratio" json:"clicks_spam_ratio"`
// 	IsFinalized                float64     `gorm:"column:is_finalized" json:"is_finalized"`
// }

func (t *CampaignModel) Validate() (err error) {
	if strings.ToLower(t.TrafficSource) == "google" {
		if utility.ValidateString(t.AdTitle) == "" {
			return errors.New("AdTitle is required")
		}
	}
	if strings.ToLower(t.TrafficSource) == "taboola" && (strings.ToLower(t.DemandSource) == "adsense" || strings.ToLower(t.DemandSource) == "tonic") {
		if utility.ValidateString(t.Account) == "" {
			return errors.New("Account is required")
		}
	}
	if strings.ToLower(t.TrafficSource) == "outbrain" {
		if utility.ValidateString(t.Account) == "" {
			return errors.New("Account is required")
		}
	}
	if strings.ToLower(t.DemandSource) == "adsense" {
		if utility.ValidateString(t.AccountAdsense) == "" {
			return errors.New("Account Adsense is required")
		}
	}
	// check title có thuộc list block k?
	return
}

func (CampaignGroupModel) TableName() string {
	return "campaign_group"
}

type CampaignGroupModel struct {
	ID            int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"column:name" json:"name"`
	TrafficSource string         `gorm:"column:traffic_source" json:"traffic_source"`
	DemandSource  string         `gorm:"column:demand_source" json:"demand_source"`
	UserID        int64          `gorm:"column:user_id" json:"user_id"`
	CreativeID    int64          `gorm:"column:creative_id" json:"creative_id"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

type SubDomainParamsAE struct {
	CampaignID int64  `json:"campaign_id" as:"campaign_id"`
	SectionID  string `json:"section_id" as:"section_id"`
}

func (SubDomainParamsAE) SetName() string {
	return "aff_sub_domain_params"
}
