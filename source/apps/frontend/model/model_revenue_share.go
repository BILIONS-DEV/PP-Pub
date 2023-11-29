package model

import (
	// "gorm.io/gorm"
	"source/core/technology/mysql"
	"time"
)

type RevenueShare struct{}

type RevenueShareRecord struct {
	mysql.TableRevenueShare
}

func (RevenueShareRecord) TableName() string {
	return mysql.Tables.RevenueShare
}

func (t *RevenueShare) RateDefault() (Rate int64) {
	Rate = mysql.RevenueShareDefault
	return
}

func (t *RevenueShare) GetByUserId(UserId int64) (record RevenueShareRecord, err error) {
	mysql.Client.Where("user_id = ?", UserId).Order("date DESC").Find(&record)

	if record.Id == 0 {
		var user UserRecord
		user = new(User).GetById(UserId)
		if err != nil {
			return
		}

		err = t.CreateRevenueShareForUser(user, t.RateDefault())
		mysql.Client.Where("user_id = ?", UserId).Order("date DESC").Find(&record)
	}
	return
}

func (t *RevenueShare) CreateRevenueShareForUser(User UserRecord, Rate int64) (err error) {
	// create rate
	record := RevenueShareRecord{}
	record.UserId = User.Id
	record.Rate = Rate
	record.Date, _ = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	err = mysql.Client.Create(&record).Error
	return
}
