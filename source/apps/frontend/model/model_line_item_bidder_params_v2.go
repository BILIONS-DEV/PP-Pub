package model

import "source/core/technology/mysql"

type LineItemBidderParamsV2 struct{}

type LineItemBidderParamsRecordV2 struct {
	mysql.TableLineItemBidderParamsV2
}

func (LineItemBidderParamsRecordV2) TableName() string {
	return mysql.Tables.LineItemBidderParamsV2
}

func (t *LineItemBidderParamsV2) GetParams(lineItemBidderId int64) (records []mysql.TableLineItemBidderParamsV2) {
	mysql.Client.Model(&LineItemBidderParamsRecordV2{}).Where("line_item_bidder_id = ?", lineItemBidderId).Find(&records)
	return
}

func (t *LineItemBidderParamsV2) DeleteByLineItemId(lineItemId int64) {
	mysql.Client.Where("line_item_id = ?", lineItemId).Delete(&LineItemBidderParamsRecordV2{})
}

func (t *LineItemBidderParamsV2) DeleteByLineItemAndBidder(lineItemId, bidderId int64) {
	mysql.Client.Where("line_item_id = ? and bidder_id = ?", lineItemId, bidderId).Delete(&LineItemBidderParamsRecordV2{})
}

func (t *LineItemBidderParamsV2) DeleteByBidderId(bidderId int64) {
	mysql.Client.Where(LineItemBidderParamsRecordV2{mysql.TableLineItemBidderParamsV2{BidderId: bidderId}}).Delete(LineItemBidderParamsRecordV2{})
	return
}
