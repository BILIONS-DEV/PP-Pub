package model

import (
	"github.com/lib/pq"
	"time"
)

func (CronjobModel) TableName() string {
	return "cronjob"
}

type CronjobModel struct {
	ID           int64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Type         TYPECronJob   `gorm:"column:type" json:"type"`
	ListObjectID pq.Int64Array `gorm:"column:list_object_id;type:integer[]" json:"list_object_id"`
	Log          string        `gorm:"column:log" json:"log"`
	Error        string        `gorm:"column:error" json:"error"`
	Status       StatusCronJob `gorm:"column:status" json:"status"`
	CreatedAt    time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

type TYPECronJob string

const (
	TYPECronJobCreateKeyValueGAM TYPECronJob = "CREATE_KEY_VALUE_GAM"
)

type StatusCronJob string

const (
	StatusCronJobQueue   StatusCronJob = "queue"
	StatusCronJobPending StatusCronJob = "pending"
	StatusCronJobSuccess StatusCronJob = "success"
	StatusCronJobError   StatusCronJob = "error"
)
