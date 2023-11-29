package report_outbrain

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

type RepoReportOutBrain interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportOutBrainModel, err error)
	Save(record *model.ReportOutBrainModel) (err error)
	SaveSlice(records []*model.ReportOutBrainModel) (err error)
	PushMessage(messages []model.TrackingAdsMessage) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportOutBrainModel, err error)
	GetToken() (token string)
}

type reportOutBrainRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportOutBrainRepo(db *gorm.DB, ka *kafka.Client) *reportOutBrainRepo {
	return &reportOutBrainRepo{Db: db, Kafka: ka}
}

type InputIsExists struct {
	MarketerID string
	CampaignID string
	SectionID  string
	Time       string
}

func (t *reportOutBrainRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportOutBrainModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportOutBrainRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportOutBrainModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportOutBrainRepo) FindByQuery(query map[string]interface{}) (record *model.ReportOutBrainModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportOutBrainRepo) Save(record *model.ReportOutBrainModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportOutBrainRepo) SaveSlice(records []*model.ReportOutBrainModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportOutBrainRepo) PushMessage(messages []model.TrackingAdsMessage) (err error) {
	var msgs []kafka2.Message
	for _, message := range messages {
		value, _ := json.Marshal(message)
		msgs = append(msgs, kafka2.Message{
			Value: value,
		})
	}
	err = t.Kafka.Writer.WriteMessages(context.Background(), msgs...)
	return
}

func (t *reportOutBrainRepo) FindCostForReportAff(day, campaignID, sectionID string) (revenue float64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportOutBrainModel{}).
		Select("SUM(spend)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&revenue).Error
	return
}

func (t *reportOutBrainRepo) FindClickForReportAff(day, campaignID, sectionID string) (clicks int64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportOutBrainModel{}).
		Select("SUM(clicks)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&clicks).Error
	return
}

func (t *reportOutBrainRepo) FindConversionForReportAff(day, campaignID, sectionID string) (conversions int64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportOutBrainModel{}).
		Select("SUM(conversions)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&conversions).Error
	return
}

func (t *reportOutBrainRepo) FindByDayForReportAff(day string) (records []*model.ReportOutBrainModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportOutBrainModel{}).
		Select("id", "marketer_id", "campaign_id", "campaign_name", "section_id", "SUM(spend) as spend", "SUM(clicks) AS clicks", "SUM(conversions) AS conversions").
		Where("time LIKE '" + day + "%'").
		Group("campaign_id, section_id").
		Find(&records).Error
	return
}

func (t *reportOutBrainRepo) GetToken() (token string) {
	var record model.TsTokenModel
	err := t.Db.
		//Debug().
		Where("traffic_source = 'outbrain'").
		Find(&record).Error
	if err != nil {
		return
	}
	token = record.Token
	return
}
