package model

func (ReportAffSelectionModel) TableName() string {
	return "report_aff_selection"
}

type ReportAffSelectionModel struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement" `
	Date          string `gorm:"column:date"`
	TrafficSource string `gorm:"column:traffic_source"`
	SectionId     string `gorm:"column:section_id"`
	SectionName   string `gorm:"column:section_name"`
	CampaignID    string `gorm:"column:campaign_id"`
	CampaignName  string `gorm:"column:campaign_name"`
	Partner       string `gorm:"column:partner"`
	PublisherID   string `gorm:"column:publisher_id"`

	Impressions   int64 `gorm:"column:impressions"`
	Click         int64 `gorm:"column:click"`
	SystemTraffic int64 `gorm:"column:system_traffic"`

	Revenue float64 `gorm:"column:revenue"`
}

func (t *ReportAffSelectionModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}
