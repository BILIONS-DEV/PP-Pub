package model

import (
	"fmt"
	"source/apps/frontend/lang"
	"source/core/technology/mysql"
	"source/pkg/utility"
	"time"
)

type UserForgetPassword struct{}

type UserForgetPasswordRecord struct {
	mysql.TableUserForgetPassword
}

func (UserForgetPasswordRecord) TableName() string {
	return mysql.Tables.UserForgetPassword
}

func (t *UserForgetPassword) GetByUserId(userId int64, email string) (row UserForgetPasswordRecord) {
	mysql.Client.Where("user_id = ? and email = ?", userId, email).Find(&row)
	return
}

func (t *UserForgetPassword) GetByUuidEmail(uuid, email string, lang lang.Translation) (err error) {
	var row UserForgetPasswordRecord
	err = mysql.Client.Where("uuid = ? and email = ?", uuid, email).Find(&row).Error
	present := time.Now()
	if row.Id == 0 || row.IsUsed == 1 {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.UserError.LinkValid.ToString())
		}
		return
	} else {
		flag := inTimeSpan(row.CreatedDate, row.ExpiredTime, present)
		if !flag {
			err = fmt.Errorf(lang.Errors.UserError.LinkOutDate.ToString())
			return
		}
	}
	return
}

func (t *UserForgetPassword) GetRecord(uuid, email string) (row UserForgetPasswordRecord) {
	mysql.Client.Table("user_forget_password").Where("uuid = ? and email = ?", uuid, email).Find(&row)
	return
}

func (t *UserForgetPassword) Handle(row *UserForgetPasswordRecord) (err error) {
	data := t.GetByUserId(row.UserId, row.Email)
	data.Uuid = row.Uuid
	data.CreatedDate = row.CreatedDate
	data.ExpiredTime = row.ExpiredTime
	data.IsUsed = row.IsUsed
	if data.Id == 0 {
		err = mysql.Client.Create(&row).Error
	} else {
		err = mysql.Client.Table("user_forget_password").Updates(&data).Where("id = ?", data.Id).Error
	}
	return
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func (t *UserForgetPassword) UpdateLinkUsed(row UserForgetPasswordRecord) (err error) {
	row.IsUsed = 1
	err = mysql.Client.Updates(&row).Error
	return
}
