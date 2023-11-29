package mysql

type TableRlsBidderMediaType struct {
	Id               int64 `gorm:"column:id" json:"id"`
	BidderTemplateId int64 `gorm:"column:bidder_template_id" json:"bidder_template_id"`
	BidderId         int64 `gorm:"column:bidder_id" json:"bidder_id"`
	MediaTypeId      int64 `gorm:"column:media_type_id" json:"media_type_id"`
}

func (TableRlsBidderMediaType) TableName() string {
	return Tables.RlBidderMediaType
}
