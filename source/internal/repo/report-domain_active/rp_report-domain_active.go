package report_domain_active

import (
	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

const KeyName = "RepoReportDomainActive"

type RepoReportDomainActive interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportDomainActiveModel, err error)
	FindAllByQuery(query map[string]interface{}) (record []*model.ReportDomainActiveModel, err error)
	Save(record *model.ReportDomainActiveModel) (err error)
	SaveSlice(records []*model.ReportDomainActiveModel) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportDomainActiveModel, err error)
}

type reportDomainActive struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportDomainActiveRP(db *gorm.DB) *reportDomainActive {
	return &reportDomainActive{
		Db: db,
	}
}

type InputIsExists struct {
	MarketerID string
	CampaignID string
	SectionID  string
	Time       string
}

func (t *reportDomainActive) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportDomainActiveModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportDomainActive) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportDomainActiveModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportDomainActive) FindByQuery(query map[string]interface{}) (record *model.ReportDomainActiveModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportDomainActive) FindAllByQuery(query map[string]interface{}) (record []*model.ReportDomainActiveModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&record).Error
	return
}

func (t *reportDomainActive) Save(record *model.ReportDomainActiveModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportDomainActive) SaveSlice(records []*model.ReportDomainActiveModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportDomainActive) FindByDayForReportAff(day string) (records []*model.ReportDomainActiveModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportDomainActiveModel{}).
		Where("time_utc LIKE '" + day + "%'").
		Find(&records).Error
	return
}
