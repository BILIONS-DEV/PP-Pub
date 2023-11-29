package model

import (
	"source/core/technology/mysql"
)

type BidderParams struct{}

type BidderParamsRecord struct {
	mysql.TableBidderParams
}

func (BidderParamsRecord) TableName() string {
	return mysql.Tables.BidderParams
}

func (t *BidderParams) GetByBidderTemplateId(bidderTemplateId int64) (records []BidderParamsRecord) {
	mysql.Client.Where("bidder_template_id = ?", bidderTemplateId).Find(&records)
	return
}

func (t *BidderParams) GetByBidderId(bidderId int64) (records []BidderParamsRecord) {
	mysql.Client.Where("bidder_id = ?", bidderId).Order("bidder_template_id desc").Find(&records)
	for k, rec := range records {
		rec.IsRequired = isRequired(rec)
		records[k] = rec
	}
	return
}

func isRequired(rec BidderParamsRecord) (required bool) {
	bidder := new(Bidder).GetByIdNoCheckUser(rec.BidderId)
	return new(PbBidder).IsParamRequiredByBidder(bidder.BidderCode, rec.Name)
}
