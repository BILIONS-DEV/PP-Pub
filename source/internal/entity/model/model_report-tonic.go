package model

func (ReportTonicModel) TableName() string {
	return "report_tonic"
}

type ReportTonicModel struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Date          string `gorm:"column:date" json:"date"`
	CampaignId    string `gorm:"column:campaign_id" json:"campaign_id"`
	CampaignName  string `gorm:"column:campaign_name" json:"campaign_name"`
	SectionId     string `gorm:"column:section_id" json:"section_id"`
	SectionName   string `gorm:"column:section_name" json:"section_name"`
	PublisherID   string `gorm:"column:publisher_id" json:"publisher_id"`
	PublisherName string `gorm:"column:publisher_name" json:"publisher_name"`
	AdID          string `gorm:"column:ad_id" json:"ad_id"`
	RedirectID    int64  `gorm:"column:redirect_id" json:"redirect_id"`
	Clicks        string `gorm:"column:clicks" json:"clicks"`
	RevenueUsd    string `gorm:"column:revenue_usd" json:"revenueUsd"`
	Subid1        string `gorm:"column:subid1" json:"subid1"`
	Subid2        string `gorm:"column:subid2" json:"subid2"`
	Subid3        string `gorm:"column:subid3" json:"subid3"`
	Subid4        string `gorm:"column:subid4" json:"subid4"`
	Keyword       string `gorm:"column:keyword" json:"keyword"`
	Timestamp     string `gorm:"column:timestamp" json:"timestamp"`
	Adtype        string `gorm:"column:adtype" json:"adtype"`
	Advertiser    string `gorm:"column:advertiser" json:"advertiser"`
	Template      string `gorm:"column:template" json:"template"`
}

func (t *ReportTonicModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}

func (t *ReportTonicModel) Validate() (err error) {
	return
}
