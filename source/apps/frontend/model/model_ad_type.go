package model

import (
	"source/core/technology/mysql"
	"source/pkg/pagination"
)

type AdType struct{}

type AdTypeRecord struct {
	mysql.TableAdType
}

func (AdTypeRecord) TableName() string {
	return mysql.Tables.AdType
}

func (t *AdType) GetById(id int64) (record AdTypeRecord) {
	if id < 1 {
		return
	}
	mysql.Client.First(&record, id)
	return
}

func (t *AdType) GetAll(userLogin, userAdmin UserRecord) (records []AdTypeRecord) {
	if userAdmin.IsFound() && (userAdmin.Permission == mysql.UserPermissionAdmin || userAdmin.Permission == mysql.UserPermissionSale) {
		mysql.Client.Find(&records)
	} else if userLogin.Permission == mysql.UserPermissionManagedService {
		mysql.Client.Where("id = 1").Find(&records)
	} else {
		mysql.Client.Find(&records)
	}
	return
}

func (t *AdType) GetAllTypeVideo() (records []AdTypeRecord) {
	mysql.Client.Where(AdTypeRecord{mysql.TableAdType{Type: mysql.TypeVideo}}).Find(&records)
	return
}

func (t *AdType) CountData(value string) (count int64) {
	mysql.Client.Model(&AdTypeRecord{}).Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *AdType) CountDataPageEdit(listSelected []int64) (count int64) {
	if len(listSelected) > 0 {
		mysql.Client.Model(&AdTypeRecord{}).Where("id not in ?", listSelected).Count(&count)
	} else {
		mysql.Client.Model(&AdTypeRecord{}).Count(&count)
	}
	return
}

func (t *AdType) LoadMoreData(key, value string, listSelected []int64) (rows []AdTypeRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination(key, limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("name like ? and id not in ?", "%"+value+"%", listSelected).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.Where("name like ?", "%"+value+"%").Limit(limit).Offset(offset).Find(&rows)
	}
	total := t.CountData(value)
	totalPages := int(total) / limit
	if (int(total) % limit) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	if page >= totalPages || len(rows) == 0 {
		isMoreData = false
		lastPage = true
	}
	return
}

func (t *AdType) LoadMoreDataPageEdit(listSelected []int64) (rows []AdTypeRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination("1", limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("id not in ?", listSelected).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.Limit(limit).Offset(offset).Find(&rows)
	}
	if len(rows) > 10 {
		rows = rows[0:9]
	}
	total := t.CountDataPageEdit(listSelected)
	totalPages := int(total) / limit
	if (int(total) % limit) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	if page >= totalPages || len(rows) == 0 {
		isMoreData = false
		lastPage = true
	}
	return
}
