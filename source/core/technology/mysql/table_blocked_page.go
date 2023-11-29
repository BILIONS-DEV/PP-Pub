package mysql

func (TableRuleBlockedPage) TableName() string {
	return "rule_blocked-page"
}

type TableRuleBlockedPage struct {
	ID     int64  `gorm:"column:id" json:"id"`
	RuleID int64  `gorm:"column:rule_id" json:"rule_id"`
	Page   string `gorm:"column:page"`
}
