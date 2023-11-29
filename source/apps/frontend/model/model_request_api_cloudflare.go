package model

import "source/core/technology/mysql"

type RequestApiCloudflare struct{}

type RequestApiCloudflareRecord struct {
	mysql.TableRequestApiCloudFlare
}

func (RequestApiCloudflareRecord) TableName() string {
	return mysql.Tables.RequestApiCloudflare
}

func (t *RequestApiCloudflare) CreateWaiting(uuid string) {
	mysql.Client.Create(&mysql.TableRequestApiCloudFlare{
		Uuid:   uuid,
		Status: mysql.TYPEStatusApiCloudflareWaiting,
	})
}
