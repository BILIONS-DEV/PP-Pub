package model

import "github.com/asaskevich/govalidator"

func (ReportFacebookModel) TableName() string {
	return "report_facebook"
}

type ReportFacebookModel struct {
	TimeUTC     string  `gorm:"column:time_utc;primaryKey;type:varchar(100)"`
	CampaignID  string  `gorm:"column:campaign_id;primaryKey;type:varchar(100)"`
	Impressions int64   `gorm:"column:impressions;type:bigint(20)"`
	Clicks      int64   `gorm:"column:clicks;type:bigint(20)"`
	Spend       float64 `gorm:"column:spend;type:double"`
	RedirectID  int64   `gorm:"column:redirect_id;type:int(11)"`
}

func (t *ReportFacebookModel) IsFound() bool {
	if !govalidator.IsNull(t.TimeUTC) {
		return true
	}
	return false
}

func (t *ReportFacebookModel) Validate() (err error) {
	return
}
