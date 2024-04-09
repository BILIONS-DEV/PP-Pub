package model

import (
	"source/core/technology/mysql"
)

const (
	ChangeStatusDomain = "change_status_domain"
)

type Cronjob struct{}

type CronjobRecord struct {
	mysql.TableCronjob
}

func (t *Cronjob) Save(record CronjobRecord) (err error) {
	err = mysql.Client.Save(&record).Error
	return
}
