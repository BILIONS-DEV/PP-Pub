package report_pocpoc

import (
	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/internal/entity/model"
)

type RepoReportPocPoc interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportPocPocModel, err error)
	Save(record *model.ReportPocPocModel) (err error)
	SaveSlice(records []*model.ReportPocPocModel) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportPocPocModel, err error)
}

type reportPocPocRepo struct {
	Db *gorm.DB
	//Cache caching.Cache
	//Kafka *kafka.Client
}

func NewReportPocPocRepo(db *gorm.DB) *reportPocPocRepo {
	return &reportPocPocRepo{Db: db}
}

type InputIsExists struct {
	MarketerID string
	CampaignID string
	SectionID  string
	Time       string
}

func (t *reportPocPocRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportPocPocModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportPocPocRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportPocPocModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportPocPocRepo) FindByQuery(query map[string]interface{}) (record *model.ReportPocPocModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportPocPocRepo) Save(record *model.ReportPocPocModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportPocPocRepo) SaveSlice(records []*model.ReportPocPocModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportPocPocRepo) FindCostForReportAff(day, campaignID, sectionID string) (revenue float64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportPocPocModel{}).
		Select("SUM(spend)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&revenue).Error
	return
}

func (t *reportPocPocRepo) FindClickForReportAff(day, campaignID, sectionID string) (clicks int64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportPocPocModel{}).
		Select("SUM(clicks)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&clicks).Error
	return
}

func (t *reportPocPocRepo) FindConversionForReportAff(day, campaignID, sectionID string) (conversions int64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportPocPocModel{}).
		Select("SUM(conversions)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&conversions).Error
	return
}

func (t *reportPocPocRepo) FindByDayForReportAff(day string) (records []*model.ReportPocPocModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportPocPocModel{}).
		Where("time LIKE '" + day + "%'").
		Find(&records).Error
	return
}
