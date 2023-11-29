package model

import (
	"source/core/technology/mysql"

	"github.com/asaskevich/govalidator"
)

type PbBidder struct{}

type PbBidderRecord struct {
	mysql.TablePbBidder
}
type PbBidderParamRecord struct {
	mysql.TablePbBidderParam
}

func (t *PbBidder) GetPbBidderByBidderCode(bidderCode string) (record PbBidderRecord) {
	mysql.Client.Where("bidder_code like ?", bidderCode).Last(&record)
	return
}

func (t *PbBidder) GetPbBidderParamCheckRequired(bidderId int64, param string) (record PbBidderParamRecord) {
	mysql.Client.Where("pb_bidder_id = ? and name like ?", bidderId, param).Last(&record)
	return
}

func (t *PbBidder) GetMapCheckParamFloor() (mapCheckParam map[string]int64) {
	mapCheckParam = make(map[string]int64)
	var records []PbBidderParamRecord
	mysql.Client.Where("description LIKE '%floor%'").Find(&records)
	for _, v := range records {
		mapCheckParam[v.Name] = v.Id
	}
	return
}

func (t *PbBidder) IsParamRequiredByBidder(bidderCode, paramName string) (flag bool) {
	bidder := new(PbBidder).GetPbBidderByBidderCode(bidderCode)
	var bidderParam PbBidderParamRecord
	mysql.Client.Select("id").Where("pb_bidder = ? AND name = ? and scope = 'required'", bidder.Name, paramName).Find(&bidderParam)
	if bidderParam.Id > 0 {
		return true
	}
	return
}

func (t *PbBidder) GetLinkByBidderCode(bidderCode string) (url string) {
	bidder := new(PbBidder).GetPbBidderByBidderCode(bidderCode)
	var bidderParam PbBidderParamRecord
	mysql.Client.Select("url").Where("pb_bidder like ?", bidder.Name).Last(&bidderParam)
	if !govalidator.IsNull(bidderParam.Url) {
		return bidderParam.Url
	}
	return "https://docs.prebid.org/dev-docs/bidders.html"
}

func (t *PbBidder) GetPbBidderParamsByBidder(bidderId int64) (records []PbBidderParamRecord) {
	mysql.Client.Where("pb_bidder_id = ?", bidderId).Find(&records)
	return
}
