package model

import "source/core/technology/mysql"

type LineItemAccount struct{}

type LineItemAccountRecord struct {
	mysql.TableLineItemAccount
}

func (LineItemAccountRecord) TableName() string {
	return mysql.Tables.LineItemAccount
}

func (t *LineItemAccount) Create(lineItemId int64, bidderId int64) (err error) {
	err = mysql.Client.Create(&LineItemAccountRecord{mysql.TableLineItemAccount{
		LineItemId: lineItemId,
		BidderId:   bidderId,
	}}).Error
	return
}

func (t *LineItemAccount) GetAllLineItem(lineItemId int64) (records []LineItemAccountRecord, err error) {
	err = mysql.Client.Where("line_item_id = ?", lineItemId).Find(&records).Error
	return
}

func (t *LineItemAccount) GetByLineItem(lineItemId int64) (records LineItemAccountRecord, err error) {
	err = mysql.Client.Where("line_item_id = ?", lineItemId).Last(&records).Error
	return
}

func (LineItemAccount) DeleteAccount(lineItemId int64) {
	mysql.Client.Where(LineItemAccountRecord{mysql.TableLineItemAccount{LineItemId: lineItemId}}).Delete(LineItemAccountRecord{})
	return
}

func (LineItemAccount) DeleteAccountByBidder(bidderId int64) {
	mysql.Client.Where(LineItemAccountRecord{mysql.TableLineItemAccount{BidderId: bidderId}}).Delete(LineItemAccountRecord{})
	return
}
