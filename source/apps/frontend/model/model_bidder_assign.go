package model

import (
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/ajax"
)

type BidderAssign struct{}

type BidderAssignRecord struct {
	mysql.TableBidderAssign
}

func (BidderAssignRecord) TableName() string {
	return mysql.Tables.BidderAssign
}

func (t *BidderAssign) Create(inputs payload.BidderCreate, user UserRecord, bidder BidderRecord) (record BidderAssignRecord, errs []ajax.Error) {
	// Validate inputs
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	// Check exist
	mysql.Client.Where("user_id = ? AND bidder_id = ?", user.Id, bidder.Id).Last(&record)
	if record.IsFound() {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "bidder is already in use",
		})
		return
	}
	// Insert to database
	record.makeInfoCreate(inputs)
	record.UserId = user.Id
	record.BidderId = bidder.Id
	err := mysql.Client.Create(&record).Error
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	return
}

func (rec *BidderAssignRecord) makeInfoCreate(inputs payload.BidderCreate) {
	rec.BidAdjustment = inputs.BidAdjustment
	rec.RefreshAd = inputs.RefreshAd
	rec.AdFormat = inputs.AdFormat
	rec.Status = inputs.Status.Standardized()
	rec.GeoIpOption = inputs.GeoIpOption.Standardized()
	rec.GeoIpList = inputs.GeoIpList
	rec.AdsTxt = inputs.AdsTxt
}

func (t *BidderAssign) ValidateCreate(inputs payload.BidderCreate) (errs []ajax.Error) {
	if inputs.BidAdjustment < 1 {
		errs = append(errs, ajax.Error{
			Id:      "bid_adjustment",
			Message: "Bid Adjustment is required",
		})
	}
	if len(inputs.AdFormat) < 1 {
		errs = append(errs, ajax.Error{
			Id:      "ad_format",
			Message: "Ad Format is required",
		})
	}
	return
}
