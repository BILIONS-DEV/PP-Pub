package model

import "github.com/asaskevich/govalidator"

type ResponsePHPTaboola struct {
	Status  bool        `json:"status"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (ReportTaboolaModel) TableName() string {
	return "report_taboola"
}

type ReportTaboolaModel struct {
	Account            string  `gorm:"column:account;type:varchar(25);primaryKey;"`
	Date               string  `gorm:"column:date;type:varchar(50);primaryKey;"`
	CampaignID         string  `gorm:"column:campaign;type:varchar(100);primaryKey;"`
	CampaignName       string  `gorm:"column:campaign_name"`
	SiteID             string  `gorm:"column:site_id;type:varchar(50);primaryKey;"`
	Site               string  `gorm:"column:site;type:varchar(255);"`
	SiteName           string  `gorm:"column:site_name;type:varchar(255);"`
	Clicks             int64   `gorm:"column:clicks;type:bigint(20);"`
	Impressions        int64   `gorm:"column:impressions;type:bigint(20);"`
	VisibleImpressions int64   `gorm:"column:visible_impressions;type:bigint(20);"`
	Spent              float64 `gorm:"column:spent;type:double;"`
	ConversionsValue   float64 `gorm:"column:conversions_value;type:double;"`
	Roas               float64 `gorm:"column:roas;type:double;"`
	Ctr                float64 `gorm:"column:ctr;type:double;"`
	Vctr               float64 `gorm:"column:vctr;type:double;"`
	Cpm                float64 `gorm:"column:cpm;type:double;"`
	Vcpm               float64 `gorm:"column:vcpm;type:double;"`
	Cpc                float64 `gorm:"column:cpc;type:double;"`
	Cpa                float64 `gorm:"column:cpa;type:double;"`
	CpaActionsNum      int64   `gorm:"column:cpa_actions_num;type:bigint(20);"`
	CpaConversionRate  float64 `gorm:"column:cpa_conversion_rate;type:bigint(20);"`
	BlockingLevel      string  `gorm:"column:blocking_level;type:varchar(255);"`
	Currency           string  `gorm:"column:currency;type:varchar(255);"`
}

func (t *ReportTaboolaModel) IsFound() bool {
	if !govalidator.IsNull(t.CampaignID) && !govalidator.IsNull(t.SiteID) {
		return true
	}
	return false
}

func (t *ReportTaboolaModel) Validate() (err error) {
	return
}
