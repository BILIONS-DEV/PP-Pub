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

func (t *UserManager) GetById(id int64) (row UserManagerRecord, err error) {
	mysql.Client.Where("id = ?", id).Find(&row)
	return
}

func (t *UserManager) GetByUser(User UserRecord) (row UserManagerRecord, err error) {
	if User.Id == 0 {
		return
	}
	mysql.Client.Where("id", User.UserManagerId).Where("presenter_id = ?", User.Presenter).Find(&row)
	return
}


// func (t *UserManager) GetManagerByPubId(presenterId int64, pubId int64) (record UserManagerRecord, err error) {
// 	err = mysql.Client.Table(mysql.Tables.UserManager).
// 		Select("user_manager.*").
// 		Joins("JOIN user_manager_pub ON user_manager_pub.user_manager_id = user_manager.id").
// 		Where("user_manager.presenter_id = ?", presenterId).
// 		Where("user_manager_pub.pub_id = ?", pubId).
// 		Find(&record).Error
// 	return
// }
