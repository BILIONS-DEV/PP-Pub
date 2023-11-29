package model

import "source/core/technology/mysql"

type RlBidderMediaType struct{}

type RlsBidderMediaTypeRecord struct {
	mysql.TableRlsBidderMediaType
}

func (RlsBidderMediaTypeRecord) TableName() string {
	return mysql.Tables.RlBidderMediaType
}

func (RlBidderMediaType) Create(rl RlsBidderMediaTypeRecord) (err error) {
	err = mysql.Client.Create(&rl).Error
	return err
}

func (RlBidderMediaType) GetByBidderTemplateId(bidderTemplateId int64) (record []RlsBidderMediaTypeRecord) {
	mysql.Client.Where(RlsBidderMediaTypeRecord{mysql.TableRlsBidderMediaType{BidderTemplateId: bidderTemplateId}}).Find(&record)
	return
}

func (RlBidderMediaType) GetByBidderId(bidderId int64) (record []RlsBidderMediaTypeRecord) {
	mysql.Client.Where(RlsBidderMediaTypeRecord{mysql.TableRlsBidderMediaType{BidderId: bidderId}}).Find(&record)
	return
}

func (RlBidderMediaType) GetByMediaTypeId(mediaTypeId int64) (record []RlsBidderMediaTypeRecord) {
	mysql.Client.Where("bidder_id != 0 and media_type_id = ?", mediaTypeId).Find(&record)
	return
}

func (RlBidderMediaType) DeleteByBidderId(bidderId int64) (record []RlsBidderMediaTypeRecord) {
	mysql.Client.Where(RlsBidderMediaTypeRecord{mysql.TableRlsBidderMediaType{BidderId: bidderId}}).Delete(RlsBidderMediaTypeRecord{})
	return
}
