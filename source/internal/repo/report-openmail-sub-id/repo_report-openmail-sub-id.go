package report_openmail_sub_id

import (
	"context"
	"encoding/json"
	kafka2 "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

type RepoReportOpenMailSubID interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByID(ID int64) (record *model.ReportOpenMailSubIdModel, err error)
	FindByQuery(query map[string]interface{}) (record *model.ReportOpenMailSubIdModel, err error)
	Save(record *model.ReportOpenMailSubIdModel) (err error)
	SaveSlice(records []*model.ReportOpenMailSubIdModel) (err error)
	PushMessage(messages []model.TrackingAdsMessage) (err error)
	FindRevenueForReportAff(day, trafficSource, campaignID, sectionID string) (revenue float64, err error)
	FindByDayForReportAff(day string) (records []*model.ReportOpenMailSubIdModel, err error)
}

type reportOpenMailSubIDRepo struct {
	Db    *gorm.DB
	Kafka *kafka.Client
}

func NewReportOpenMailSubIDRepo(DB *gorm.DB, ka *kafka.Client) *reportOpenMailSubIDRepo {
	return &reportOpenMailSubIDRepo{Db: DB, Kafka: ka}
}

type InputIsExists struct {
	UtcHour  string
	Campaign string
	SubID    string
}

func (t *reportOpenMailSubIDRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if input.UtcHour != "" {
					return db.Where("utc_hour = ?", input.UtcHour)
				}
				return db
			},
			func(db *gorm.DB) *gorm.DB {
				if input.Campaign != "" {
					return db.Where("campaign = ?", input.Campaign)
				}
				return db
			},
			func(db *gorm.DB) *gorm.DB {
				if input.SubID != "" {
					return db.Where("sub_id = ?", input.SubID)
				}
				return db
			},
		).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportOpenMailSubIdModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportOpenMailSubIDRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportOpenMailSubIdModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportOpenMailSubIDRepo) FindByID(ID int64) (record *model.ReportOpenMailSubIdModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		First(&record, ID).Error
	return
}

func (t *reportOpenMailSubIDRepo) FindByQuery(query map[string]interface{}) (record *model.ReportOpenMailSubIdModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		First(&record).Error
	return
}

func (t *reportOpenMailSubIDRepo) Save(record *model.ReportOpenMailSubIdModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportOpenMailSubIDRepo) SaveSlice(records []*model.ReportOpenMailSubIdModel) (err error) {
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportOpenMailSubIDRepo) PushMessage(messages []model.TrackingAdsMessage) (err error) {
	var msgs []kafka2.Message
	for _, message := range messages {
		message.Partner = "system1"
		value, _ := json.Marshal(message)
		msgs = append(msgs, kafka2.Message{
			Value: value,
		})
	}
	err = t.Kafka.Writer.WriteMessages(context.Background(), msgs...)
	return
}

func (t *reportOpenMailSubIDRepo) FindRevenueForReportAff(day, trafficSource, campaignID, sectionID string) (revenue float64, err error) {
	subID := trafficSource + "_" + campaignID + ":" + sectionID
	err = t.Db.
		//Debug().
		Model(model.ReportOpenMailSubIdModel{}).
		Select("SUM(estimated_revenue)").
		Where("utc_hour LIKE '"+day+"%'").
		Where("sub_id = ?", subID).
		Find(&revenue).Error
	return
}

func (t *reportOpenMailSubIDRepo) FindByDayForReportAff(day string) (records []*model.ReportOpenMailSubIdModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportOpenMailSubIdModel{}).
		Select("id", "campaign", "utc_hour", "sub_id", "SUM(searches) as searches", "SUM(clicks) AS clicks", "SUM(estimated_revenue) AS estimated_revenue").
		Where("utc_hour LIKE '" + day + "%'").
		Group("sub_id").
		Find(&records).Error
	return
}
