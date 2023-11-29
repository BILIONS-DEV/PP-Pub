package model

import "source/core/technology/mysql"

type BidderTemplateParams struct{}

type BidderTemplateParamsRecord struct {
	mysql.TableBidderTemplateParams
}

func (BidderTemplateParamsRecord) TableName() string {
	return mysql.Tables.BidderTemplateParams
}

func (t *BidderTemplateParams) GetByBidderTemplateId(id int64) (records []BidderTemplateParamsRecord) {
	mysql.Client.Where("bidder_template_id = ?", id).Find(&records)
	return
}
