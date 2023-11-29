package model

import (
	"source/core/technology/mysql"
	"source/pkg/htmlblock"
	"strconv"
	"time"
)

type Notification struct{}

type NotificationRecord struct {
	mysql.TableNotification
}

func (NotificationRecord) TableName() string {
	return mysql.Tables.Notification
}

func (t *Notification) GetById(gamId int64) (rec NotificationRecord) {
	mysql.Client.Last(&rec, gamId)
	return
}

func (t *Notification) GetByUser(userId int64) (rec []NotificationRecord) {
	mysql.Client.Where("user_id = ?", userId).Limit(30).Order("id desc").Find(&rec)
	return
}

func (t *Notification) GetNewByUser(userId int64) (rec int64) {
	mysql.Client.Table(mysql.Tables.Notification).Where("user_id = ? AND status = 1", userId).Count(&rec)
	return
}

type TemplateNotification struct {
	NumberNew int64  `json:"NumberNew"`
	Template  string `json:"Template"`
}

func (t *Notification) GetNotifications(userId int64) (rec TemplateNotification) {
	reccord := t.GetByUser(userId)
	if reccord == nil {
		return
	}

	NumberNew := t.GetNewByUser(userId)
	notifactions := t.MakeResponse(reccord)
	var Template string
	for _, notification := range notifactions {
		if Template == "" {
			Template = notification.Notification
		} else {
			Template = Template + "\n" + notification.Notification
		}
	}

	rec.Template = Template
	rec.NumberNew = NumberNew

	return
}

type NotificationResponse struct {
	NotificationRecord
	Notification string `json:"Notification"`
	TimeString   string `json:"TimeString"`
}

func (t *Notification) MakeResponse(notifications []NotificationRecord) (records []NotificationResponse) {
	if notifications == nil {
		return
	}
	for _, notification := range notifications {
		TimeString := TimeAgo(notification.CreatedAt)
		rec := NotificationResponse{
			NotificationRecord: notification,
			TimeString:         TimeString,
		}
		records = append(records, rec)
	}

	for key, record := range records {
		record.Notification = htmlblock.Render("notification/block.notification.gohtml", record).String()
		records[key] = record
	}
	return
}

func TimeAgo(t time.Time) (result string) {
	// timeNow := time.Now()
	if time.Now().Unix()-t.Unix() < 60 {
		result = "Just now"
		return
	}
	if time.Now().Unix()-t.Unix() < 3600 {
		minute := (time.Now().Unix() - t.Unix())/60
		result = strconv.Itoa(int(minute)) + " minutes ago"
		return
	}

	if time.Now().Unix()-t.Unix() < 84600 {
		minute := (time.Now().Unix() - t.Unix()) / 3600
		result = strconv.Itoa(int(minute)) + " hour ago"
		return
	}

	result = t.Format("2006 Jan 02")
	return
}

func (t *Notification) ReadNotificationsForUser(userId, id int64) (result bool) {
	Query := mysql.Client.Table(mysql.Tables.Notification).Where("user_id = ? AND status = 1", userId)
	if id > 0 {
		Query.Where("id = ?", id)
	}
	err := Query.Update("status", 2).Error
	if err == nil {
		result = true
	}
	return
}

func (t *Notification) CreateNotify(notify NotificationRecord) (err error) {
	mysql.Client.Create(&notify)
	return
}