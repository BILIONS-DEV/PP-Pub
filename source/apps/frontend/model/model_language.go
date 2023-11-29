package model

import (
	"source/core/technology/mysql"
	"source/pkg/pagination"
)

type Language struct{}

type LanguageRecord struct {
	mysql.TableLanguage
}

func (LanguageRecord) TableName() string {
	return mysql.Tables.Language
}

func (t *Language) GetAll() (records []LanguageRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *Language) GetById(id int64) (records LanguageRecord) {
	mysql.Client.Find(&records, id)
	return
}

func (t *Language) GetLanguageNameById(id int64) (language string) {
	var records LanguageRecord
	mysql.Client.Find(&records, id)
	language = records.LanguageName
	return
}

func (t *Language) GetByLanguageCode(code string) (records LanguageRecord) {
	mysql.Client.Find(&records, code)
	return
}

func (t *Language) LoadMoreData(key, value string, listSelected []int64) (rows []LanguageRecord, isMoreData, lastPage bool) {
	limit := 30
	page, offset := pagination.Pagination(key, limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("language_name like ? and code not in ?", "%"+value+"%", listSelected).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.Where("language_name like ?", "%"+value+"%").Limit(limit).Offset(offset).Find(&rows)
	}

	total := t.CountData(value)
	totalPages := int(total) / 30
	if (int(total) % 30) != 0 {
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

func (t *Language) CountData(value string) (count int64) {
	mysql.Client.Model(&LanguageRecord{}).Where("language_name like ?", "%"+value+"%").Count(&count)
	return
}
