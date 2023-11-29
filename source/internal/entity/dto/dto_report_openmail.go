package dto

import "source/internal/entity/model"

type ReportOpenMailSubIdEstimatedHourly struct {
	SubIds []ReportOpenMailSubId `json:"subids" mapstructure:"subids"`
}

type ReportOpenMailSubId struct {
	UtcHour          string  `json:"utc_hour" mapstructure:"utc_hour"`
	Campaign         string  `json:"campaign" mapstructure:"campaign"`
	SubID            string  `json:"sub_id" mapstructure:"sub_id"`
	Searches         int64   `json:"searches" mapstructure:"searches"`
	Clicks           int64   `json:"clicks" mapstructure:"clicks"`
	EstimatedRevenue float64 `json:"estimated_revenue" mapstructure:"estimated_revenue"`
	LastUpdated      string  `json:"last_updated" mapstructure:"last_updated"`
}

func (t *ReportOpenMailSubId) ToModel() model.ReportOpenMailSubIdModel {
	record := model.ReportOpenMailSubIdModel{
		UtcHour:          t.UtcHour,
		Campaign:         t.Campaign,
		SubID:            t.SubID,
		Searches:         t.Searches,
		Clicks:           t.Clicks,
		EstimatedRevenue: t.EstimatedRevenue,
		LastUpdated:      t.LastUpdated,
	}
	return record
}
