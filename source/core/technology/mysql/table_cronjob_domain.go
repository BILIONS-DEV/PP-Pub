package mysql

type TableCronjobDomain struct {
	Id      int64  `gorm:"column:id" json:"id"`
	Domain  string `gorm:"column:domain" json:"domain"`
	Type    string `gorm:"column:type" json:"type"`
	FeedUrl string `gorm:"column:feed_url" json:"feed_url"`
	Status  string `gorm:"column:status" json:"status"`
}

func (TableCronjobDomain) TableName() string {
	return Tables.CronjobDomain
}
