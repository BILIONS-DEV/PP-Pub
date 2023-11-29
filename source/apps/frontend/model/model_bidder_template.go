package model

import "source/core/technology/mysql"

type BidderTemplate struct{}

type BidderTemplateRecord struct {
	mysql.TableBidderTemplate
}

func (BidderTemplateRecord) TableName() string {
	return mysql.Tables.BidderTemplate
}

func (t *BidderTemplate) GetBidderCodeById(id int64) string {
	rec := BidderTemplateRecord{}
	mysql.Client.Select("bidder_code").First(&rec, id)
	return rec.BidderCode
}

func (t *BidderTemplate) GetById(id int64) (record BidderTemplateRecord) {
	mysql.Client.First(&record, id)
	return
}

func (t *BidderTemplate) GetAll() (records []BidderTemplateRecord) {
	//Chỉ lấy toàn bộ list các bidder dùng cho pub với bidder_type = 1
	mysql.Client.Where("bidder_type = 1").Find(&records)
	return
}

func (t *BidderTemplate) CheckUniqueGoogleDisplayNameWithPrebidDemand(googleDisplayName string) (flag bool) {
	var rec BidderTemplateRecord
	mysql.Client.Select("id").Where("bidder_code = ? OR display_name = ?", googleDisplayName, googleDisplayName).Last(&rec)
	if rec.Id > 0 {
		flag = true
	} else {
		flag = false
	}
	return
}
