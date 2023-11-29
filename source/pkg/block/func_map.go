package block

import (
	"fmt"
	"source/core/technology/mysql"
	"text/template"
)

func FunMaps() template.FuncMap {
	funcMap := template.FuncMap{
		"IncIndex":     IncIndex,
		"GetEmailById": GetEmailById,
	}
	return funcMap
}

type UserRecord struct {
	Email string
}

func GetEmailById(userId int64) string {
	fmt.Println("GetEmailById")
	var user UserRecord
	mysql.Client.Table(mysql.Tables.User).Select("email").First(&user, userId)
	return user.Email
}

func IncIndex(number int) int {
	return number + 1
}
