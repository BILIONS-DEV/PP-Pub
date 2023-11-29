package model

import (
	"source/core/technology/mysql"
	"source/pkg/pagination"
)

type Device struct{}

type DeviceRecord struct {
	mysql.TableDevice
}

func (DeviceRecord) TableName() string {
	return mysql.Tables.Device
}

func (t *Device) GetById(id int64) (record DeviceRecord) {
	if id < 1 {
		return
	}
	mysql.Client.First(&record, id)
	return
}

func (t *Device) GetAll() (records []DeviceRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *Device) CountData(value string) (count int64) {
	mysql.Client.Model(&DeviceRecord{}).Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *Device) CountDataPageEdit(listSelected []int64) (count int64) {
	if len(listSelected) > 0 {
		mysql.Client.Model(&DeviceRecord{}).Where("id not in ?", listSelected).Count(&count)
	} else {
		mysql.Client.Model(&DeviceRecord{}).Count(&count)
	}
	return
}

func (t *Device) LoadMoreData(key, value string, listSelected []int64) (rows []DeviceRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination(key, limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("name like ? and id not in ?", "%"+value+"%", listSelected).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.Where("name like ?", "%"+value+"%").Limit(limit).Offset(offset).Find(&rows)
	}
	total := t.CountData(value)
	totalPages := int(total) / 10
	if (int(total) % 10) != 0 {
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

func (t *Device) LoadMoreDataPageEdit(listSelected []int64) (rows []DeviceRecord, isMoreData, lastPage bool) {
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
