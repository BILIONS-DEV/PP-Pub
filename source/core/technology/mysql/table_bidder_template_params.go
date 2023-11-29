package mysql

type TableBidderTemplateParams struct {
	Id               int64  `gorm:"column:id" json:"id"`
	UserId           int64  `gorm:"column:user_id" json:"user_id"`
	BidderTemplateId int64  `gorm:"column:bidder_template_id" json:"bidder_template_id"`
	Name             string `gorm:"column:name" json:"name"`
	Type             string `gorm:"column:type" json:"type"`
	Template         string `gorm:"column:template" json:"template"`
}

func (TableBidderTemplateParams) TableName() string {
	return Tables.BidderTemplateParams
}
