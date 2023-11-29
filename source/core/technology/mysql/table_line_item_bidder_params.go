package mysql

type TableLineItemBidderParams struct {
	Id               int64  `json:"id"`
	LineItemId       int64  `json:"line_item_id"`
	BidderId         int64  `json:"bidder_id"`
	LineItemBidderId int64  `json:"line_item_bidder_id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Value            string `json:"value"`
}

func (TableLineItemBidderParams) TableName() string {
	return Tables.LineItemBidderParams
}