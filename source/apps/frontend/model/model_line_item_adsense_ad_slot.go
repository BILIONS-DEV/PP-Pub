package model

import (
	"source/core/technology/mysql"
)

type LineItemAdsenseAdSlot struct{}

type LineItemAdSenseAdSlotRecord struct {
	mysql.LineItemAdsenseAdSlot
}

func (t *LineItemAdsenseAdSlot) Push(lineItemId int64, Size string, adsenseAdSlotId string) (err error) {
	//var rec LineItemAdSenseAdSlotRecord
	//mysql.Client.Select("id").
	//	Where("line_item_id = ? AND size = ? AND adsense_ad_slot_id = ?", lineItemId, Size, adsenseAdSlotId).
	//	First(&rec)
	//if rec.Id > 0 {
	//	return errors.New("record already exists")
	//}
	err = mysql.Client.
		Create(&LineItemAdSenseAdSlotRecord{mysql.LineItemAdsenseAdSlot{
			LineItemId:      lineItemId,
			Size:            Size,
			AdsenseAdSlotId: adsenseAdSlotId,
		}}).Error
	return
}

func (t *LineItemAdsenseAdSlot) DeleteByLineItem(lineItemId int64) (err error) {
	err = mysql.Client.Unscoped().Where(LineItemAdSenseAdSlotRecord{mysql.LineItemAdsenseAdSlot{LineItemId: lineItemId}}).Delete(&LineItemAdSenseAdSlotRecord{}).Error
	return
}

func (t *LineItemAdsenseAdSlot) GetByLineItem(lineItemId int64) (records []LineItemAdSenseAdSlotRecord, err error) {
	err = mysql.Client.Where(LineItemAdSenseAdSlotRecord{mysql.LineItemAdsenseAdSlot{LineItemId: lineItemId}}).Find(&records).Error
	return
}
