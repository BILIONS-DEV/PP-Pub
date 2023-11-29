package report_mgid

import (
	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/internal/entity/model"
)

type RepoReportMgid interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportMgidModel, err error)
	Save(record *model.ReportMgidModel) (err error)
	SaveSlice(records []*model.ReportMgidModel) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportMgidModel, err error)
}

type reportMgidRepo struct {
	Db *gorm.DB
	//Cache caching.Cache
	//Kafka *kafka.Client
}

func NewReportMgidRepo(db *gorm.DB) *reportMgidRepo {
	return &reportMgidRepo{Db: db}
}

type InputIsExists struct {
	MarketerID string
	CampaignID string
	SectionID  string
	Time       string
}

func (t *reportMgidRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportMgidModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportMgidRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportMgidModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportMgidRepo) FindByQuery(query map[string]interface{}) (record *model.ReportMgidModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportMgidRepo) Save(record *model.ReportMgidModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportMgidRepo) SaveSlice(records []*model.ReportMgidModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportMgidRepo) FindCostForReportAff(day, campaignID, sectionID string) (revenue float64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportMgidModel{}).
		Select("SUM(spend)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&revenue).Error
	return
}

func (t *reportMgidRepo) FindClickForReportAff(day, campaignID, sectionID string) (clicks int64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportMgidModel{}).
		Select("SUM(clicks)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&clicks).Error
	return
}

func (t *reportMgidRepo) FindConversionForReportAff(day, campaignID, sectionID string) (conversions int64, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportMgidModel{}).
		Select("SUM(conversions)").
		Where("time LIKE '"+day+"%'").
		Where("campaign_id = ?", campaignID).
		Where("section_id = ?", sectionID).
		Find(&conversions).Error
	return
}

func (t *reportMgidRepo) FindByDayForReportAff(day string) (records []*model.ReportMgidModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportMgidModel{}).
		Where("time LIKE '" + day + "%'").
		Find(&records).Error
	return
}
