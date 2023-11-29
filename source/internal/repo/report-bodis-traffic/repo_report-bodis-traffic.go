package report_bodis_traffic

import (
	"context"
	"encoding/json"
	kafka2 "github.com/segmentio/kafka-go"

	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

type RepoReportBodisTraffic interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportBodisTrafficModel, err error)
	Save(record *model.ReportBodisTrafficModel) (err error)
	SaveSlice(records []*model.ReportBodisTrafficModel) (err error)
	PushMessage(messages []model.TrackingAdsMessage) (err error)
}

type reportBodisRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportBodisTrafficRepo(db *gorm.DB, ka *kafka.Client) *reportBodisRepo {
	return &reportBodisRepo{Db: db, Kafka: ka}
}

type InputIsExists struct {
	VisitID       string `gorm:"column:visit_id"`
	Time          string `gorm:"column:time"`
	TrafficSource string `gorm:"column:traffic_source"`
	Campaign      string `gorm:"column:campaign"`
	SelectionID   string `gorm:"column:selection_id"`
}

func (t *reportBodisRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportBodisTrafficModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportBodisRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportBodisTrafficModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportBodisRepo) FindByQuery(query map[string]interface{}) (record *model.ReportBodisTrafficModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		First(&record).Error
	return
}

func (t *reportBodisRepo) Save(record *model.ReportBodisTrafficModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportBodisRepo) SaveSlice(records []*model.ReportBodisTrafficModel) (err error) {
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportBodisRepo) PushMessage(messages []model.TrackingAdsMessage) (err error) {
	var msgs []kafka2.Message
	for _, message := range messages {
		message.Partner = "bodis"
		value, _ := json.Marshal(message)
		msgs = append(msgs, kafka2.Message{
			Value: value,
		})
	}
	_ = t.Kafka.Writer.WriteMessages(context.Background(), msgs...)
	return
}
