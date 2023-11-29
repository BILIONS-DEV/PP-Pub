package model

func (ReportMgidModel) TableName() string {
	return "report_mgid"
}

type ReportMgidModel struct {
	Time       string  `gorm:"column:time;type:varchar(50);primaryKey;"`
	CampaignID string  `gorm:"column:campaign_id;type:varchar(100);primaryKey;"`
	SectionID  string  `gorm:"column:section_id;type:varchar(100);primaryKey;"`
	Clicks     int64   `gorm:"column:clicks;type:bigint(20)"`
	Spent      float64 `gorm:"column:spent;type:double"`
	CPC        float64 `gorm:"column:cpc;type:double"`
}

func (t *ReportMgidModel) IsFound() bool {
	if t.CampaignID != "" {
		return true
	}
	return false
}

func (t *ReportMgidModel) Validate() (err error) {
	return
}
