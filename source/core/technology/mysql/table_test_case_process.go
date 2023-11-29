package mysql

type TableTestCaseProcess struct {
	Id         int64   `gorm:"column:id" json:"id"`
	TestCaseId int64   `gorm:"column:test_case_id" json:"test_case_id"`
	FloorTest  string  `gorm:"column:floor_test" json:"floor_test"`
	Winner     float64 `gorm:"column:winner" json:"winner"`
	WinCPM     float64 `gorm:"column:win_cpm" json:"win_cpm"`
	StartAt    int64   `gorm:"column:start_at" json:"start_at"`
	Status     string  `gorm:"column:status" json:"status"`
}

func (TableTestCaseProcess) TableName() string {
	return Tables.TestCaseProcess
}
