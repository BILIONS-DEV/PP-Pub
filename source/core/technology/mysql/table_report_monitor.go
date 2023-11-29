package mysql

import (
	"time"
)

type TableReportMonitor struct {
	ID                 int64              `gorm:"column:id" json:"id"`
	ObservedDimensions ObservedDimensions `gorm:"column:observed_dimensions" json:"observed_dimensions"`
	RuleMetrics        RuleMetrics        `gorm:"column:rule_metrics" json:"rule_metrics"`
	NotifyUsers        NotifyUsers        `gorm:"column:notify_users" json:"notify_users"`
	NumberCompareDate  int                `gorm:"column:number_compare_date" json:"number_compare_date"`
	CompareHour        int                `gorm:"column:compare_hour" json:"compare_hour"`
	Status             string             `gorm:"column:status" json:"status"`
	LastRun            time.Time          `gorm:"column:last_run" json:"last_run"`
}

func (TableReportMonitor) TableName() string {
	return Tables.ReportMonitor
}
