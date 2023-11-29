package mysql

type RuleBlockedPage struct {
	ID     int64  `gorm:"column:id" json:"id"`
	RuleID int64  `gorm:"column:rule_id" json:"rule_id"`
	Page   string `gorm:"column:page" json:"page"`
}

func (RuleBlockedPage) TableName() string {
	return "rule_blocked-page"
}