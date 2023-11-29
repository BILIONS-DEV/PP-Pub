package mysql

type TableBidderParams struct {
	Id               int64  `gorm:"column:id" json:"id"`
	UserId           int64  `gorm:"column:user_id" json:"user_id"`
	BidderId         int64  `gorm:"column:bidder_id" json:"bidder_id"`
	BidderTemplateId int64  `gorm:"column:bidder_template_id" json:"bidder_template_id"`
	Name             string `gorm:"column:name" json:"name"`
	Type             string `gorm:"column:type" json:"type"`
	Template         string `gorm:"column:template" json:"template"`
	IsRequired       bool   `gorm:"-"`
}

func (TableBidderParams) TableName() string {
	return Tables.BidderParams
}
