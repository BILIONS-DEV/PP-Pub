package model

import "source/core/technology/mysql"

type AdTagSizeAdditional struct{}

type AdTagSizeAdditionalRecord struct {
	mysql.TableRlsAdTagSizeAdditional
}

func (AdTagSizeAdditionalRecord) TableName() string {
	return mysql.Tables.RlAdTagSizeAdditional
}

func (this *AdTagSizeAdditional) Create(adTagId, device, adSizeId int64) {
	mysql.Client.
		//Debug().
		Create(&AdTagSizeAdditionalRecord{mysql.TableRlsAdTagSizeAdditional{
			AdTagId:  adTagId,
			Device:   device,
			AdSizeId: adSizeId,
		}})
}

func (this *AdTagSizeAdditional) DeleteAllUnscopedByTagId(adTagId int64) {
	mysql.Client.Unscoped().Where(AdTagSizeAdditionalRecord{mysql.TableRlsAdTagSizeAdditional{AdTagId: adTagId}}).Delete(&AdTagSizeAdditionalRecord{})
}

func (this *AdTagSizeAdditional) DeleteAllByTagId(adTagId int64) {
	mysql.Client.Where(AdTagSizeAdditionalRecord{mysql.TableRlsAdTagSizeAdditional{AdTagId: adTagId}}).Delete(&AdTagSizeAdditionalRecord{})
}

func (this *AdTagSizeAdditional) GetAllByTagId(adTagId int64) (records []AdTagSizeAdditionalRecord) {
	mysql.Client.Where(AdTagSizeAdditionalRecord{mysql.TableRlsAdTagSizeAdditional{AdTagId: adTagId}}).Find(&records)
	return
}
