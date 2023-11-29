package model

func (ReportCodeFuelModel) TableName() string {
	return "report_codefuel"
}

type ReportCodeFuelModel struct {
	Time                   string  `gorm:"column:time;type:varchar(50);primaryKey;"`
	CampaignID             int64   `gorm:"column:campaign_id;type:int(11);primaryKey;"`
	TotalMonetizedSearches int64   `gorm:"column:total_monetized_searches;type:bigint(20)"`
	AdClicks               int64   `gorm:"column:ad_clicks;type:bigint(20)"`
	Amount                 float64 `gorm:"column:amount;type:double"`
}

func (t *ReportCodeFuelModel) IsFound() bool {
	if t.CampaignID != 0 {
		return true
	}
	return false
}

func (t *ReportCodeFuelModel) Validate() (err error) {
	return
}
