package model

import (
	"source/core/technology/mysql"
)

type UserManager struct{}

type UserManagerRecord struct {
	mysql.TableUserManager
}

func (UserManagerRecord) TableName() string {
	return mysql.Tables.UserManager
}

type UserManagerPub struct{}

type UserManagerPubRecord struct {
	mysql.TableUserManagerPub
}

func (UserManagerPubRecord) TableName() string {
	return mysql.Tables.UserManagerPub
}

func (t *UserManager) GetManagerByPubId(presenterId int64, pubId int64) (record UserManagerRecord, err error) {
	err = mysql.Client.Table(mysql.Tables.UserManager).
		Select("user_manager.*").
		Joins("JOIN user_manager_pub ON user_manager_pub.user_manager_id = user_manager.id").
		Where("user_manager.presenter_id = ?", presenterId).
		Where("user_manager_pub.pub_id = ?", pubId).
		Find(&record).Error
	return
}