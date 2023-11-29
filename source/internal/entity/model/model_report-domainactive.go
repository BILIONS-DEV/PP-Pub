package model

func (ReportDomainActiveModel) TableName() string {
	return "report_domainactive"
}

type ReportDomainActiveModel struct {
	Tg1          string `gorm:"column:tg1;primaryKey;" json:"tg1"`
	CampaignName string `gorm:"column:campaign__name" json:"campaign__name"`
	//CampaignType           string  `gorm:"column:campaign__type" json:"campaign__type"`
	CampaignId int64 `gorm:"column:campaign_id;primaryKey;" json:"campaign_id"`
	//Hour                   int64   `gorm:"column:hour" json:"hour"`
	//LanderSearches         int64   `gorm:"column:lander_searches" json:"lander_searches"`
	//LanderVisitors         int64   `gorm:"column:lander_visitors" json:"lander_visitors"`
	PublisherRevenueAmount float64 `gorm:"column:publisher_revenue_amount" json:"publisher_revenue_amount"`
	RevenueClicks          int64   `gorm:"column:revenue_clicks" json:"revenue_clicks"`
	TotalVisitors          int64   `gorm:"column:total_visitors" json:"total_visitors"`
	TrackedVisitors        int64   `gorm:"column:tracked_visitors" json:"tracked_visitors"`
	Date                   string  `gorm:"column:date;primaryKey;" json:"date"`
	Gclid                  string  `gorm:"column:gclid" json:"gclid"`
	TrafficSource          string  `gorm:"column:traffic_source" json:"traffic_source"`
	RedirectID             int64   `gorm:"column:redirect_id" json:"redirect_id"`
	SectionID              string  `gorm:"column:section_id" json:"section_id"`
	TimeUTC                string  `gorm:"column:time_utc" json:"time_utc"`
}

//func (t *ReportDomainActiveModel) IsFound() bool {
//
//	return false
//}

func (t *ReportDomainActiveModel) IsFound() bool {
	if t.CampaignId == 0 {
		return true
	}
	return false
}

func (t *ReportDomainActiveModel) Validate() (err error) {
	return
}
