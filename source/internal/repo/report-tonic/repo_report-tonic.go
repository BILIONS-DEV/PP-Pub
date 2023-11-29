package report_tonic

import (
	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

type RepoReportTonic interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportTonicModel, err error)
	FindAllByQuery(query map[string]interface{}) (record []*model.ReportTonicModel, err error)
	Save(record *model.ReportTonicModel) (err error)
	SaveSlice(records []*model.ReportTonicModel) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportTonicModel, err error)
}

type reportTonic struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportTonic(db *gorm.DB) *reportTonic {
	return &reportTonic{Db: db}
}

type InputIsExists struct {
	MarketerID string
	CampaignID string
	SectionID  string
	Time       string
}

func (t *reportTonic) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportTonicModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportTonic) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportTonicModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportTonic) FindByQuery(query map[string]interface{}) (record *model.ReportTonicModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportTonic) FindAllByQuery(query map[string]interface{}) (record []*model.ReportTonicModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&record).Error
	return
}

func (t *reportTonic) Save(record *model.ReportTonicModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportTonic) SaveSlice(records []*model.ReportTonicModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportTonic) FindByDayForReportAff(day string) (records []*model.ReportTonicModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportTonicModel{}).
		Select("id", "date", "campaign_id", "campaign_name", "SUM(revenue_usd) as revenue_usd", "SUM(revenue_usd) as revenue_usd", "subid1", "timestamp").
		Where("date LIKE '" + day + "%'").
		Group("subid1").
		Find(&records).Error
	return
}
