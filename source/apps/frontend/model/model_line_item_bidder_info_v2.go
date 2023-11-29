package model

import "source/core/technology/mysql"

type LineItemBidderInfoV2 struct{}

type LineItemBidderInfoRecordV2 struct {
	mysql.TableLineItemBidderInfoV2
}

func (LineItemBidderInfoRecordV2) TableName() string {
	return mysql.Tables.LineItemBidderInfoV2
}

func (LineItemBidderInfoV2) CreateBidderInfo(record LineItemBidderInfoRecordV2) LineItemBidderInfoRecordV2 {
	mysql.Client.
		//Debug().
		Create(&record)
	return record
}

func (LineItemBidderInfoV2) DeleteBidderInfo(lineItemId int64) {
	mysql.Client.Where(LineItemBidderInfoRecordV2{mysql.TableLineItemBidderInfoV2{LineItemId: lineItemId}}).Delete(LineItemBidderInfoRecordV2{})
	return
}

func (LineItemBidderInfoV2) DeleteBidderInfoForLineItem(lineItemId, bidderId int64) {
	mysql.Client.Where(LineItemBidderInfoRecordV2{mysql.TableLineItemBidderInfoV2{LineItemId: lineItemId, BidderId: bidderId}}).Delete(LineItemBidderInfoRecordV2{})
	return
}

func (LineItemBidderInfoV2) GetBidderInfoByLineItem(lineItemId int64) (records []LineItemBidderInfoRecordV2) {
	mysql.Client.
		//Debug().
		Where("line_item_id = ?", lineItemId).Find(&records)
	return
}

func (LineItemBidderInfoV2) GetBidderInfoByBidderId(bidderId int64) (records []LineItemBidderInfoRecordV2) {
	mysql.Client.Where("bidder_id = ?", bidderId).Find(&records)
	return
}

func (LineItemBidderInfoV2) DeleteBidderByBidderId(bidderId int64) {
	mysql.Client.Where(LineItemBidderInfoRecordV2{mysql.TableLineItemBidderInfoV2{BidderId: bidderId}}).Delete(LineItemBidderInfoRecordV2{})
	return
}
