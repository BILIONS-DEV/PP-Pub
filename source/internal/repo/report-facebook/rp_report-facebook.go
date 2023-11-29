package report_facebook

import (
	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

const KeyName = "RepoReportFacebook"

type RepoReportFacebook interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportFacebookModel, err error)
	FindAllByQuery(query map[string]interface{}) (record []*model.ReportFacebookModel, err error)
	Save(record *model.ReportFacebookModel) (err error)
	SaveSlice(records []*model.ReportFacebookModel) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportFacebookModel, err error)
}

type reportFacebook struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportFacebookRP(db *gorm.DB) *reportFacebook {
	return &reportFacebook{
		Db: db,
	}
}

type InputIsExists struct {
	MarketerID string
	CampaignID string
	SectionID  string
	Time       string
}

func (t *reportFacebook) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportFacebookModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportFacebook) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportFacebookModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportFacebook) FindByQuery(query map[string]interface{}) (record *model.ReportFacebookModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportFacebook) FindAllByQuery(query map[string]interface{}) (record []*model.ReportFacebookModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&record).Error
	return
}

func (t *reportFacebook) Save(record *model.ReportFacebookModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportFacebook) SaveSlice(records []*model.ReportFacebookModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportFacebook) FindByDayForReportAff(day string) (records []*model.ReportFacebookModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportFacebookModel{}).
		Where("time_utc LIKE '" + day + "%'").
		Find(&records).Error
	return
}
