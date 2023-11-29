package model

import (
	"source/core/technology/mysql"
	"source/pkg/pagination"
)

type Country struct{}

type CountryRecord struct {
	mysql.TableCountry
}

func (CountryRecord) TableName() string {
	return mysql.Tables.Country
}

func (t *Country) GetById(id int64) (record CountryRecord) {
	if id < 1 {
		return
	}
	mysql.Client.First(&record, id)
	return
}

func (t *Country) GetAll() (records []CountryRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *Country) CountData(value string) (count int64) {
	mysql.Client.Model(&CountryRecord{}).Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *Country) CountDataPageEdit(listSelected []int64) (count int64) {
	if len(listSelected) > 0 {
		mysql.Client.Model(&CountryRecord{}).Where("id not in ?", listSelected).Count(&count)
	} else {
		mysql.Client.Model(&CountryRecord{}).Count(&count)
	}
	return
}

func (t *Country) LoadMoreData(key, value string, listSelected []int64) (rows []CountryRecord, isMoreData, lastPage bool) {
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

func (t *Country) LoadMoreDataPageEdit(listSelected []int64) (rows []CountryRecord, isMoreData, lastPage bool) {
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
