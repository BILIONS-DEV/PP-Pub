package mysql

type TableLineItemBidderInfo struct {
	Id         int64                       `json:"id" form:"id"`
	LineItemId int64                       `json:"line_item_id"`
	BidderId   int64                       `json:"bidder_id"`
	ConfigType TYPEConfigType              `json:"config_type"`
	BidderType TYPEBidderType              `json:"bidder_type"`
	Name       string                      `json:"name"`
	Params     []TableLineItemBidderParams `gorm:"-"`
}

func (TableLineItemBidderInfo) TableName() string {
	return Tables.LineItemBidderInfo
}

type TYPEConfigType string

const (
	TYPEConfigTypeBanner         TYPEConfigType = "banner"
	TYPEConfigTypeVideoInstream  TYPEConfigType = "video_instream"
	TYPEConfigTypeVideoOutStream TYPEConfigType = "video_outstream"
	TYPEConfigTypeNative         TYPEConfigType = "native"
)
