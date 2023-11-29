package model

import "source/core/technology/mysql"

type LineItemBidderInfo struct{}

type LineItemBidderInfoRecord struct {
	mysql.TableLineItemBidderInfo
}

func (LineItemBidderInfoRecord) TableName() string {
	return mysql.Tables.LineItemBidderInfo
}

func (LineItemBidderInfo) CreateBidderInfo(record LineItemBidderInfoRecord) LineItemBidderInfoRecord {
	mysql.Client.Create(&record)
	return record
}

func (LineItemBidderInfo) DeleteBidderInfo(lineItemId int64) {
	mysql.Client.Where(LineItemBidderInfoRecord{mysql.TableLineItemBidderInfo{LineItemId: lineItemId}}).Delete(LineItemBidderInfoRecord{})
	return
}

func (LineItemBidderInfo) GetBidderInfoByLineItem(lineItemId int64) (records []LineItemBidderInfoRecord) {
	mysql.Client.Where("line_item_id = ?", lineItemId).Find(&records)
	return
}

func (LineItemBidderInfo) GetBidderInfoByBidderId(bidderId int64) (records []LineItemBidderInfoRecord) {
	mysql.Client.Where("bidder_id = ?", bidderId).Find(&records)
	return
}

func (LineItemBidderInfo) DeleteBidderByBidderId(bidderId int64) {
	mysql.Client.Where(LineItemBidderInfoRecord{mysql.TableLineItemBidderInfo{BidderId: bidderId}}).Delete(LineItemBidderInfoRecord{})
	return
}
