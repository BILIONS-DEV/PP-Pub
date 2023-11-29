package model

import (
	"encoding/json"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
)

type ModuleUserId struct{}

type ModuleUserIdRecord struct {
	mysql.TableModuleUserId
}

func (ModuleUserIdRecord) TableName() string {
	return mysql.Tables.ModuleUserId
}

func (t *ModuleUserId) GetById(id int64) (record ModuleUserIdRecord) {
	mysql.Client.First(&record, id)
	return
}

func (t *ModuleUserId) GetAll() (records []ModuleUserIdRecord, err error) {
	err = mysql.Client.Find(&records).Error
	return
}

func (t *ModuleUserId) ConvertParamToList(param string) (list []payload.ParamModuleUserId, err error) {
	err = json.Unmarshal([]byte(param), &list)
	if err != nil {
		return
	}
	return
}

func (t *ModuleUserId) ConvertStorageToList(storage string) (list []payload.StorageModuleUserId, err error) {
	err = json.Unmarshal([]byte(storage), &list)
	if err != nil {
		return
	}
	return
}

//func (t *ModuleUserId) GetDefaultModule() (rows []ModuleUserIdRecord, list []int64, err error) {
//	err = mysql.Client.Where("name in ('pubCommonId','flocId','criteo','id5Id')").Find(&rows).Error
//	if err != nil {
//		return
//	}
//	err = mysql.Client.Model(&ModuleUserIdRecord{}).Select("id").Where("name in ('pubCommonId','flocId','criteo','id5Id')").Find(&list).Error
//	if err != nil {
//		return
//	}
//	return
//}
