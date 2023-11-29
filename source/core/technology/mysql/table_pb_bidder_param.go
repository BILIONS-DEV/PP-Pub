package mysql

type TablePbBidderParam struct {
	Id           int64  `gorm:"column:id" json:"id"`
	PbBidderId   int64  `gorm:"column:pb_bidder_id" json:"pb_bidder_id"`
	PbBidder     string `gorm:"column:pb_bidder" json:"pb_bidder"`
	BidderCode   string `gorm:"column:bidder_code" json:"bidder_code"`
	Name         string `gorm:"column:name" json:"name"`
	PrebidType   int    `gorm:"column:prebid_type" json:"prebid_type"`
	MediaType    string `gorm:"column:media_type" json:"media_type"`
	Scope        string `gorm:"column:scope" json:"scope"`
	Description  string `gorm:"column:description" json:"description"`
	Type         string `gorm:"column:type" json:"type"`
	Example      string `gorm:"column:example" json:"example"`
	DefaultValue string `gorm:"column:default_value" json:"default_value"`
	Url          string `gorm:"column:url" json:"url"`
	NumberTable  int    `gorm:"column:number_table" json:"number_table"`
	H3           string `gorm:"column:h_3" json:"h_3"`
	Status       string `gorm:"column:status" json:"status"`
}

func (TablePbBidderParam) TableName() string {
	return Tables.PbBidderParam
}
