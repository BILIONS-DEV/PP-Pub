package model

import (
	"source/core/technology/mysql"
)

type ManagerSub struct{}

type ManagerSubRecord struct {
	mysql.TableManagerSub
}

func (ManagerSubRecord) TableName() string {
	return mysql.Tables.ManagerSub
}

func (t *ManagerSub) GetById(id int64) (row ManagerSubRecord, err error) {
	mysql.Client.Where("id = ?", id).Find(&row)
	return
}

func (t *ManagerSub) GetByUser(User UserRecord) (row ManagerSubRecord, err error) {
	if User.Id == 0 {
		return
	}
	mysql.Client.Where("id", User.ManagerSubId).Where("presenter_id = ?", User.Presenter).Find(&row)
	return
}


// func (t *ManagerSub) GetManagerByPubId(presenterId int64, pubId int64) (record ManagerSubRecord, err error) {
// 	err = mysql.Client.Table(mysql.Tables.ManagerSub).
// 		Select("user_manager.*").
// 		Joins("JOIN user_manager_pub ON user_manager_pub.user_manager_id = user_manager.id").
// 		Where("user_manager.presenter_id = ?", presenterId).
// 		Where("user_manager_pub.pub_id = ?", pubId).
// 		Find(&record).Error
// 	return
// }
