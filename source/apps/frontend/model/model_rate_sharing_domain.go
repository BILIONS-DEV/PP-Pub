package model

import (
	"source/core/technology/mysql"
	"strconv"
	"time"
)

type RateSharingInventory struct{}

type RateSharingInventoryRecord struct {
	mysql.TableRateSharingInventory
}

func (RateSharingInventoryRecord) TableName() string {
	return mysql.Tables.RateSharingInventory
}

func (t *RateSharingInventory) RateDefault() (Rate int64) {
	Rate = mysql.RevenueShareDefault
	return
}

func (t *RateSharingInventory) GetRevShareByDomain(inventoryId int64) (rate int64, err error) {
	var record RateSharingInventoryRecord
	mysql.Client.Where("inventory_id = ? and rate_type = 'revenue'", inventoryId).Order("created_at DESC").Find(&record)
	if record.Id > 0 {
		rate = record.Rate
	} else {
		rate = 100
	}
	return
}

func (t *RateSharingInventory) GetRateFillAdjustedByDomain(inventoryId int64) (rate int64, err error) {
	var record RateSharingInventoryRecord
	mysql.Client.Where("inventory_id = ? and rate_type = 'fill-rate'", inventoryId).Order("created_at DESC").Find(&record)
	if record.Id > 0 {
		rate = record.Rate
	} else {
		rate = 100
	}
	return
}

func (t *RateSharingInventory) GetImpShareByDomain(inventoryId int64) (rate int64, err error) {
	var record RateSharingInventoryRecord
	mysql.Client.Where("inventory_id = ? and rate_type = 'impressions'", inventoryId).Order("created_at DESC").Find(&record)
	if record.Id > 0 {
		rate = record.Rate
	} else {
		rate = 100
	}
	return
}

func (t *RateSharingInventory) GetRateSharing(rateType string, inventoryId int64) (record RateSharingInventoryRecord, err error) {
	if rateType == "fill-rate" || rateType == "impressions" {
		timeCreate := time.Now()
		utc, _ := time.LoadLocation("Iceland")
		applicableDate, _ := strconv.ParseInt(timeCreate.In(utc).Format("20060102"), 10, 64)
		mysql.Client.Where("inventory_id = ? and rate_type = ? and applicable_date = ? ", inventoryId, rateType, applicableDate).Order("created_at DESC").Find(&record)
	} else {
		mysql.Client.Where("inventory_id = ? and rate_type = ?", inventoryId, rateType).Order("created_at DESC").Find(&record)
	}
	return
}

func (t *RateSharingInventory) CreateRateSharing(rateType string, rate int64, inventoryId int64) (err error) {
	var record RateSharingInventoryRecord
	record.InventoryId = inventoryId
	record.RateType = rateType
	record.Rate = rate
	timeCreate := time.Now()
	utc, _ := time.LoadLocation("Iceland")
	record.CreatedAt = timeCreate
	applicableDate := timeCreate.In(utc).Format("20060102")
	record.ApplicableDate, _ = strconv.ParseInt(applicableDate, 10, 64)
	err = mysql.Client.Create(&record).Error
	return
}

func (t *RateSharingInventory) UpdateRateSharing(record RateSharingInventoryRecord) (err error) {
	err = mysql.Client.Save(&record).Error
	return
}

func (t *RateSharingInventory) UpdateRateSharingForDomain(inventoryId int64, rate int64, rateType string) (err error) {
	if rateType == "revenue" {
		err = t.CreateRateSharing(rateType, rate, inventoryId)
	} else {
		fillSharing, _ := t.GetRateSharing(rateType, inventoryId)
		if fillSharing.Id == 0 {
			err = t.CreateRateSharing(rateType, rate, inventoryId)
		} else {
			fillSharing.Rate = rate
			err = t.UpdateRateSharing(fillSharing)
		}
	}
	return
}
