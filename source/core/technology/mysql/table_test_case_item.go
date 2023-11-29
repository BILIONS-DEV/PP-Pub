package mysql

import (
	"time"
)

type TableTestCaseItem struct {
	Id            int64     `gorm:"column:id" json:"id"`
	TestCaseId    int64     `gorm:"column:test_case_id" json:"test_case_id"`
	TestProcessId int64     `gorm:"column:test_process_id" json:"test_process_id"`
	Hour          int       `gorm:"column:_hour" json:"_hour"`
	TagId         int64     `gorm:"column:tag_id" json:"tag_id"`
	AdUnitCode    string    `gorm:"column:ad_unit_code" json:"ad_unit_code"`
	CountryCode   string    `gorm:"column:country_code" json:"country_code"`
	Device        string    `gorm:"column:device" json:"device"`
	CreatedTime   time.Time `gorm:"column:created_time" json:"created_time"`
}

func (TableTestCaseItem) TableName() string {
	return Tables.TestCaseItem
}
