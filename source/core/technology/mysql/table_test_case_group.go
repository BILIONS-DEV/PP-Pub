package mysql

import (
	"database/sql"
	"time"
)

type TableTestCaseGroup struct {
	Id               int64        `gorm:"column:id" json:"id"`
	MinCPM           float64      `gorm:"column:min_cpm" json:"min_cpm"`
	MaxCPM           float64      `gorm:"column:max_cpm" json:"max_cpm"`
	AvgCPM           float64      `gorm:"column:avg_cpm" json:"avg_cpm"`
	MinFloor         float64      `gorm:"column:min_floor" json:"min_floor"`
	MaxFloor         float64      `gorm:"column:max_floor" json:"max_floor"`
	FloorList        string       `gorm:"column:floor_list" json:"floor_list"`
	Winner           float64      `gorm:"column:winner" json:"winner"`
	GroupTime        string       `gorm:"column:group_time" json:"group_time"`
	TotalRequests    int64        `gorm:"column:total_requests" json:"total_requests"`
	TotalImpressions int64        `gorm:"column:total_impressions" json:"total_impressions"`
	Status           string       `gorm:"column:status" json:"status"`
	CreateTime       time.Time    `gorm:"column:create_time" json:"create_time"`
	UpdateTime       sql.NullTime `gorm:"column:update_time" json:"update_time"`
}

func (TableTestCaseGroup) TableName() string {
	return Tables.TestCaseGroup
}
