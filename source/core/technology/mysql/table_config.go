package mysql

type TableConfig struct {
	Id                    int64  `gorm:"column:id" json:"id"`
	UserId                int64  `gorm:"column:user_id" json:"user_id"`
	PrebidTimeOut         int    `gorm:"column:prebid_time_out" json:"prebid_time_out"`
	AdRefreshTime         int    `gorm:"column:ad_refresh_time" json:"ad_refresh_time"`
	Currency              string `gorm:"column:currency" json:"currency"`
	GranularityMultiplier int  `gorm:"column:granularity_multiplier" json:"granularity_multiplier"`
}
