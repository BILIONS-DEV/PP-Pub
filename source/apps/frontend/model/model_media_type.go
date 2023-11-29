package model

import "source/core/technology/mysql"

type MediaType struct{}

type MediaTypeRecord struct {
	mysql.TableMediaType
}

func (MediaTypeRecord) TableName() string {
	return mysql.Tables.MediaType
}

func (t *MediaType) GetById(id int64) (record MediaTypeRecord) {
	mysql.Client.First(&record, id)
	return
}

func (t *MediaType) GetAll() (records []MediaTypeRecord) {
	mysql.Client.Find(&records)
	return
}
