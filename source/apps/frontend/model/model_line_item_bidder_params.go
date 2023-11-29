package model

import "source/core/technology/mysql"

type LineItemBidderParams struct{}

type LineItemBidderParamsRecord struct {
	mysql.TableLineItemBidderParams
}

func (LineItemBidderParamsRecord) TableName() string {
	return mysql.Tables.LineItemBidderParams
}

func (t *LineItemBidderParams) GetParams(lineItemBidderId int64) (records []mysql.TableLineItemBidderParams) {
	mysql.Client.Model(&LineItemBidderParamsRecord{}).Where("line_item_bidder_id = ?", lineItemBidderId).Find(&records)
	return
}

func (t *LineItemBidderParams) DeleteByLineItemId(lineItemId int64) {
	mysql.Client.Where("line_item_id = ?", lineItemId).Delete(&LineItemBidderParamsRecord{})
}

func (t *LineItemBidderParams) DeleteByBidderId(bidderId int64) {
	mysql.Client.Where(LineItemBidderParamsRecord{mysql.TableLineItemBidderParams{BidderId: bidderId}}).Delete(LineItemBidderParamsRecord{})
	return
}
