package model

import "github.com/asaskevich/govalidator"

func (ReportTikTokModel) TableName() string {
	return "report_tiktok"
}

type ReportTikTokModel struct {
	StartTimeHour string  `gorm:"column:stat_time_hour;primaryKey;type:varchar(100)"`
	AdvertiserID  string  `gorm:"column:advertiser_id;primaryKey;type:varchar(100)"`
	CampaignID    string  `gorm:"column:campaign_id;primaryKey;type:varchar(100)"`
	AdGroupID     string  `gorm:"column:adgroup_id;primaryKey;type:varchar(100)"`
	CampaignName  string  `gorm:"column:campaign_name;type:varchar(255)"`
	Impressions   int64   `gorm:"column:impressions;type:bigint(20)"`
	Clicks        int64   `gorm:"column:clicks;type:bigint(20)"`
	Spend         float64 `gorm:"column:spend;type:double"`
	TimeUTC       string  `gorm:"column:time_utc;type:varchar(100)"`
	RedirectID    int64   `gorm:"column:redirect_id;type:int(11)"`
}

func (t *ReportTikTokModel) IsFound() bool {
	if !govalidator.IsNull(t.StartTimeHour) {
		return true
	}
	return false
}

func (t *ReportTikTokModel) Validate() (err error) {
	return
}
