package mysql

type TableBlockRPM struct {
	Id          int64   `gorm:"column:id;primary_key" json:"id"`
	UserId      int64   `gorm:"column:user_id" json:"user_id"`
	BidderCode  string  `gorm:"column:bidderCode" json:"bidderCode"`
	AdUnitCode  string  `gorm:"column:adUnitCode" json:"adUnitCode"`
	TagId       int64   `gorm:"column:tagId" json:"tagId"`
	Device      string  `gorm:"column:device" json:"device"`
	CountryCode string  `gorm:"column:countryCode" json:"countryCode"`
	Status      string  `gorm:"column:status" json:"status"`
	Test        int     `gorm:"column:test" json:"test"`
	Rpm         float64 `gorm:"column:rpm" json:"rpm"`
	BlockDate   int     `gorm:"column:blocked_date" json:"blocked_date"`
}

func (TableBlockRPM) TableName() string {
	return Tables.BlockRPM
}
