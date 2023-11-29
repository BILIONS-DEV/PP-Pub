package model

func (ReportPocPocModel) TableName() string {
	return "report_pocpoc"
}

type ReportPocPocModel struct {
	Time         string  `gorm:"column:time;type:varchar(50);primaryKey;"`
	CampaignID   string  `gorm:"column:campaign_id;type:varchar(100);primaryKey;"`
	CampaignName string  `gorm:"column:campaign_name;type:varchar(500)"`
	SectionID    string  `gorm:"column:section_id;type:varchar(100);primaryKey;"`
	SectionName  string  `gorm:"column:section_name;type:varchar(500)"`
	Clicks       int64   `gorm:"column:clicks;type:bigint(20)"`
	Spent        float64 `gorm:"column:spent;type:double"`
	CPC          float64 `gorm:"column:cpc;type:double"`
}

func (t *ReportPocPocModel) IsFound() bool {
	if t.CampaignID != "" {
		return true
	}
	return false
}

func (t *ReportPocPocModel) Validate() (err error) {
	return
}
