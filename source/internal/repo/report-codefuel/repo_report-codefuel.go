package report_codefuel

import (
	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

type RepoReportCodeFuel interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportCodeFuelModel, err error)
	FindAllByQuery(query map[string]interface{}) (record []*model.ReportCodeFuelModel, err error)
	Save(record *model.ReportCodeFuelModel) (err error)
	SaveSlice(records []*model.ReportCodeFuelModel) (err error)
	FindByDayForReportAff(startDate, endDate string) (records []*model.ReportCodeFuelModel, err error)
}

type reportCodeFuel struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportCodeFuel(db *gorm.DB) *reportCodeFuel {
	return &reportCodeFuel{Db: db}
}

type InputIsExists struct {
	MarketerID string
	CampaignID string
	SectionID  string
	Time       string
}

func (t *reportCodeFuel) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportCodeFuelModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportCodeFuel) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportCodeFuelModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportCodeFuel) FindByQuery(query map[string]interface{}) (record *model.ReportCodeFuelModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportCodeFuel) FindAllByQuery(query map[string]interface{}) (record []*model.ReportCodeFuelModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&record).Error
	return
}

func (t *reportCodeFuel) Save(record *model.ReportCodeFuelModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportCodeFuel) SaveSlice(records []*model.ReportCodeFuelModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportCodeFuel) FindByDayForReportAff(startDate, endDate string) (records []*model.ReportCodeFuelModel, err error) {
	err = t.Db.
		Debug().
		Model(model.ReportCodeFuelModel{}).
		Where("time >= ? and time <= ?", startDate, endDate).
		Find(&records).Error
	return
}
