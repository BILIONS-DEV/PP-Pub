package model

import (
	"source/core/technology/mysql"
)

type CronJobBlockedPage struct{}

type CronJobBlockedPageRecord struct {
	mysql.TableCronJobBlockedPage
}

func (CronJobBlockedPageRecord) TableName() string {
	return mysql.Tables.CronJobBlockedPage
}

func (t *CronJobBlockedPage) Save(record *CronJobBlockedPageRecord) (err error) {
	err = mysql.Client.
		//Debug().
		Save(&record).
		Error
	return
}

func (t *CronJobBlockedPage) UnBlockAllByUser(userID int64) (err error) {
	err = mysql.Client.
		//Debug().
		Model(&CronJobBlockedPageRecord{}).
		Where("user_id = ?", userID).
		Update("type", mysql.TYPECronjobBlockedPageUnBlock).
		Error
	return
}
