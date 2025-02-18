package model

import (
	"source/core/technology/mysql"
)

type UserAccount struct{}

type UserAccountRecord struct {
	mysql.TableUserAccount
}

func (UserAccountRecord) TableName() string {
	return mysql.Tables.UserAccount
}