package model

import "source/core/technology/mysql"

type ListParamValidate struct{}

type ListParamValidateRecord struct {
	mysql.TableListParamValidate
}

func (ListParamValidateRecord) TableName() string {
	return mysql.Tables.ListParamValidate
}

func (t *ListParamValidate) GetMapCheckParam() (mapCheckParam map[string]int64) {
	mapCheckParam = make(map[string]int64)
	var records []ListParamValidateRecord
	mysql.Client.Find(&records)
	for _, v := range records {
		mapCheckParam[v.Name] = v.Id
	}
	return
}
