package mysql

import "github.com/lib/pq"

type TablePrebidModule struct {
	Name             string         `gorm:"column:name" json:"name"`
	BidAdapter       int            `gorm:"column:bid_adapter" json:"bid_adapter"`
	AnalyticsAdapter int            `gorm:"column:analytics_adapter" json:"analytics_adapter"`
	IdSystem         int            `gorm:"column:id_system" json:"id_system"`
	AdServerVideo    int            `gorm:"column:ad_server_video" json:"ad_server_video"`
	Floors           int            `gorm:"column:floors" json:"floors"`
	Files            pq.StringArray `gorm:"type:text[],column:files" json:"files"`
}

func (TablePrebidModule) TableName() string {
	return Tables.PrebidModule
}
