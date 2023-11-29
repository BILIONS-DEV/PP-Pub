package model

import "source/core/technology/mysql"

type Currency struct {}

type CurrencyRecord struct {
	mysql.TableCurrency
}

func (t *Currency) GetAll() (records []CurrencyRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *Currency) GetById(id int64) (record CurrencyRecord) {
	mysql.Client.Find(&record, id)
	return
}

func (t *Currency) GetByCode(code string) (record CurrencyRecord) {
	mysql.Client.Where("code = ?", code).Find(&record)
	return
}