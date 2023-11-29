package model

import (
	// "gorm.io/gorm"
	"source/core/technology/mysql"
)

type RateSharing struct{}

type RateSharingRecord struct {
	mysql.TableRateSharing
}

func (RateSharingRecord) TableName() string {
	return mysql.Tables.RateSharing
}

func (t *RateSharing) RateDefault() (Rate int64) {
	Rate = mysql.RevenueShareDefault
	return
}

func (t *RateSharing) GetByUserId(UserId int64) (record RateSharingRecord, err error) {
	mysql.Client.Where("user_id = ? and domain_id = 0", UserId).Find(&record)

	if record.Id == 0 {
		var user UserRecord
		user = new(User).GetById(UserId)
		err = t.UpdateRateSharingForUser(user, t.RateDefault())
		mysql.Client.Where("user_id = ? and domain_id = 0", UserId).Find(&record)
	}
	return
}

func (t *RateSharing) GetByDomain(Domain InventoryRecord) (record RateSharingRecord, err error) {
	mysql.Client.Where("domain_id = ?", Domain.Id).Find(&record)
	//if record.Rate > 0 {
	//	return
	//}

	// nếu rate chưa đc set theo domain thì sẽ lấy rate theo user
	record, err = t.GetByUserId(Domain.UserId)
	return
}

func (t *RateSharing) UpdateRateSharingForUser(User UserRecord, Rate int64) (err error) {
	// remove rate old
	err = mysql.Client.Delete(&RateSharingRecord{}, "domain_id = ? and user_id = ?", 0, User.Id).Error
	if err != nil {
		return
	}
	// create rate
	record := RateSharingRecord{}
	//record.UserId = User.Id
	//record.DomainId = 0
	//record.Rate = Rate
	//record.Date = time.Now().Format("2006-01-02")
	err = mysql.Client.Create(&record).Error
	return
}
