package report_tiktok

import (
	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

const KeyName = "RepoReportTikTok"

type RepoReportTikTok interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportTikTokModel, err error)
	FindAllByQuery(query map[string]interface{}) (record []*model.ReportTikTokModel, err error)
	Save(record *model.ReportTikTokModel) (err error)
	SaveSlice(records []*model.ReportTikTokModel) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportTikTokModel, err error)
}

type reportTikTok struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportTikTokRP(db *gorm.DB) *reportTikTok {
	return &reportTikTok{
		Db: db,
	}
}

type InputIsExists struct {
	MarketerID string
	CampaignID string
	SectionID  string
	Time       string
}

func (t *reportTikTok) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportTikTokModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportTikTok) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportTikTokModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportTikTok) FindByQuery(query map[string]interface{}) (record *model.ReportTikTokModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportTikTok) FindAllByQuery(query map[string]interface{}) (record []*model.ReportTikTokModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&record).Error
	return
}

func (t *reportTikTok) Save(record *model.ReportTikTokModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportTikTok) SaveSlice(records []*model.ReportTikTokModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportTikTok) FindByDayForReportAff(day string) (records []*model.ReportTikTokModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportTikTokModel{}).
		Select("stat_time_hour", "advertiser_id", "adgroup_id", "campaign_id", "campaign_name", "SUM(impressions) as impressions", "SUM(clicks) as clicks", "SUM(spend) as spend", "time_utc", "redirect_id").
		Where("time_utc LIKE '" + day + "%'").
		Where("redirect_id != 0").
		Group("redirect_id, adgroup_id").
		Find(&records).Error
	return
}
