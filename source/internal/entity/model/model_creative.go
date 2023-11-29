package model

import (
	"errors"
)

type CreativeModel struct {
	ID     int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID int64 `gorm:"column:user_id" json:"user_id"`
	// CampaignID int64                `gorm:"column:campaign_id" json:"campaign_id"`
	SiteName string               `gorm:"column:site_name" json:"site_name"`
	Titles   []CreativeTitleModel `gorm:"foreignKey:CreativeID;references:ID"`
	Images   []CreativeImageModel `gorm:"foreignKey:CreativeID;references:ID"`
}

func (CreativeModel) TableName() string {
	return "campaign_creative"
}

type CreativeTitleModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreativeID int64  `gorm:"column:creative_id" json:"creative_id"`
	Title      string `gorm:"column:title" json:"title"`
}

func (CreativeTitleModel) TableName() string {
	return "creative_title"
}

type CreativeImageModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreativeID int64  `gorm:"column:creative_id" json:"creative_id"`
	Image      string `gorm:"column:image" json:"image"`
}

func (CreativeImageModel) TableName() string {
	return "creative_image"
}

type CreativeSubmitModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignID int64  `gorm:"column:campaign_id" json:"cmapaign_id"`
	CreativeID int64  `gorm:"column:creative_id" json:"creative_id"`
	SiteName   string `gorm:"column:site_name" json:"site_name"`
	Title      string `gorm:"column:title" json:"title"`
	Image      string `gorm:"column:image" json:"image"`
	New        string `gorm:"column:new" json:"new"`
	Flag       int    `gorm:"column:flag" json:"flag"`
}

func (CreativeSubmitModel) TableName() string {
	return "creative_submit"
}

type CreativeSubmitFilter struct {
	CampaignID int `json:"campaign_id" form:"campaign_id"`
	CreativeID int `json:"creative_id" form:"creative_id"`
	Page       int `query:"page" json:"page" form:"page"`
	Limit      int `query:"limit" json:"limit" form:"limit"`
	Start      int `query:"start" json:"start" form:"start"`
	Length     int `query:"length" json:"length" form:"length"`
}

// type ResponseData struct {
// 	Subid                      interface{} `gorm:"column:subid" json:"subid"`
// 	Visits                     int         `gorm:"column:visits" json:"visits"`
// 	LandingPageVisits          int         `gorm:"column:landing_page_visits" json:"landing_page_visits"`
// 	Zeroclick_visits           int         `gorm:"column:zeroclick_visits" json:"zeroclick_visits"`
// 	Clicks                     int         `gorm:"column:clicks" json:"clicks"`
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

func (t *CreativeModel) Validate() (err error) {
	// check title có thuộc list block k?
	var ListBlock []string
	ListBlock = append(ListBlock, "Affordable", "Angriest", "Best", "Biggest", "Brightest", "Broadest", "Busiest", "Calmest", "Cheapest", "Cleanest", "Closest", "Compare", "Corona", "covid", "Deadliest", "Deal", "Deals", "Easiest", "Fastest", "Fewest", "Finest", "For Free", "Free", "Greatest", "Healthiest", "Heaviest", "Here", "Here's", "Highest", "inexpensive", "largest", "Latest", "Local", "Lowest", "Medicaid", "Medicare", "Monkeypox", "Most Affordable", "Most difficult", "Most interesting", "Nastiest", "Newest", "Nicest", "Nicest", "Poorest", "ppp", "Quickest", "Rarest", "Richest", "Roomiest", "Roughest", "Safest", "Save", "See Prices", "see symptoms", "See treatments", "Shortest", "Signs", "Simplest", "Slimmest", "Slowest", "Smallest", "Smartest", "Strangest", "Symptoms", "Tallest", "Thinnest", "Tiniest", "Top", "Toughest", "Weakest", "Worst", "Youngest")
	if len(t.Titles) > 0 {
		for _, value := range t.Titles {
			if contains(ListBlock, value.Title) {
				err = errors.New("title " + value.Title + " belongs to block list")
				return
			}
		}
	}

	return
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
