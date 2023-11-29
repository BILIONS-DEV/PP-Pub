package model

import "source/core/technology/mysql"

type BlockingRestrictions struct{}

type BlockingRestrictionsRecord struct {
	mysql.TableBlockingRestrictions
}

func (BlockingRestrictionsRecord) TableName() string {
	return mysql.Tables.BlockingRestrictions
}

func (t *BlockingRestrictions) CreateAdvertiseDomain(blockingId int64, advertiseDomain string) (err error) {
	err = mysql.Client.Create(&BlockingRestrictionsRecord{mysql.TableBlockingRestrictions{
		BlockingId:      blockingId,
		AdvertiseDomain: advertiseDomain,
	}}).Error
	return
}

func (t *BlockingRestrictions) CreateCreativeId(blockingId int64, creativeId string) (err error) {
	err = mysql.Client.Create(&BlockingRestrictionsRecord{mysql.TableBlockingRestrictions{
		BlockingId: blockingId,
		CreativeId: creativeId,
	}}).Error
	return
}

func (t *BlockingRestrictions) GetAdvertiseDomainByBlocking(blockingId int64) (records []BlockingRestrictionsRecord, err error) {
	err = mysql.Client.Where("blocking_id = ? and advertiser_domain != ''", blockingId).Find(&records).Error
	return
}

func (t *BlockingRestrictions) GetCreativeIdByBlocking(blockingId int64) (records []BlockingRestrictionsRecord, err error) {
	err = mysql.Client.Where("blocking_id = ? and creative_id != ''", blockingId).Find(&records).Error
	return
}

func (t *BlockingRestrictions) DeleteByBlockingId(blockingId int64) (err error) {
	err = mysql.Client.Unscoped().Where("blocking_id = ?", blockingId).Delete(&BlockingRestrictionsRecord{}).Error
	return
}
