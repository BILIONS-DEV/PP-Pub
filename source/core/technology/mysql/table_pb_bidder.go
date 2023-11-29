package mysql

type TablePbBidder struct {
	Id         int64  `gorm:"column:id" json:"id"`
	Name       string `gorm:"column:name" json:"name"`
	BidderCode string `gorm:"column:bidder_code" json:"bidder_code"`
	Status     string `gorm:"column:status" json:"status"`
}

func (TablePbBidder) TableName() string {
	return Tables.PbBidder
}
