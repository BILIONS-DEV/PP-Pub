package mysql

type TableLineItemBidderInfoV2 struct {
	Id         int64                         `json:"id" form:"id"`
	LineItemId int64                         `json:"line_item_id"`
	BidderId   int64                         `json:"bidder_id"`
	ConfigType TYPEConfigType                `json:"config_type"`
	BidderType TYPEBidderType                `json:"bidder_type"`
	Name       string                        `json:"name"`
	Params     []TableLineItemBidderParamsV2 `gorm:"-"`
}

func (TableLineItemBidderInfoV2) TableName() string {
	return Tables.LineItemBidderInfoV2
}
