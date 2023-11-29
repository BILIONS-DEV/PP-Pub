package model

import "encoding/json"

type ReportOpenMailSubIdModel struct {
	ID               int64   `gorm:"column:id;primaryKey;autoIncrement"`
	UtcHour          string  `gorm:"column:utc_hour"`
	SectionId        string  `gorm:"column:section_id"`
	SectionName      string  `gorm:"column:section_name"`
	PublisherID      string  `gorm:"column:publisher_id"`
	PublisherName    string  `gorm:"column:publisher_name"`
	AdID             string  `gorm:"column:ad_id"`
	RedirectID       int64   `gorm:"column:redirect_id"`
	Campaign         string  `gorm:"column:campaign"`
	SubID            string  `gorm:"column:sub_id"`
	Searches         int64   `gorm:"column:searches"`
	Clicks           int64   `gorm:"column:clicks"`
	EstimatedRevenue float64 `gorm:"column:estimated_revenue"`
	LastUpdated      string  `gorm:"column:last_updated"`
}

type ReportOpenMailRequestDruid struct {
	UtcHour       string  `json:"utcHour"`
	TrafficSource string  `json:"trafficSource"`
	Campaign      string  `json:"campaign"`
	SelectionID   string  `json:"selectionId"`
	Searches      int64   `json:"searches"`
	Clicks        int64   `json:"clicks"`
	Revenue       float64 `json:"revenue"`
}

func (ReportOpenMailSubIdModel) TableName() string {
	return "report_open_mail_sub_id"
}

func (t *ReportOpenMailSubIdModel) Validate() (err error) {
	return
}

func (t *ReportOpenMailSubIdModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}

func (t *ReportOpenMailSubIdModel) ToJSON() string {
	jsonEncode, _ := json.Marshal(t)
	return string(jsonEncode)
}
