package report_aff_pixel

import (
	"strings"
	"time"

	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

type RepoReportAffPixel interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindOneByQuery(query map[string]interface{}, queryRaw ...string) (record *model.ReportAffPixelModel, err error)
	FindAllByQuery(query map[string]interface{}) (records []*model.ReportAffPixelModel, err error)
	FindForReportAff(day string) (records []*model.ReportAffPixelModel, err error)
	FindForReportAffSelection(day string) (records []*model.ReportAffPixelModel, err error)
	FindAllByGroup(group ...string) (records []*model.ReportAffPixelModel, err error)
	Save(record *model.ReportAffPixelModel) (err error)
	SaveSlice(records []*model.ReportAffPixelModel) (err error)
	FindCampaignIDForAdsense(redirectID int64) (campaignID string, err error)
	FindRedirectID(campaignID string) (redirectID string, err error)
	FindForReportAllTrafficSourceIDs(trafficSource, account string) (ids []string, err error)
}

type reportAffPixelRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportAffPixelRepo(db *gorm.DB, cache caching.Cache, ka *kafka.Client) *reportAffPixelRepo {
	return &reportAffPixelRepo{Db: db, Cache: cache, Kafka: ka}
}

func (t *reportAffPixelRepo) setFilterDate(startDate, endDate string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var listDate []string
		layout := "2006-01-02"
		start, _ := time.Parse(layout, startDate)
		listDate = append(listDate, start.Format(layout))
		if startDate != endDate {
			for {
				// Add thêm 1 ngày
				start = start.AddDate(0, 0, 1)
				listDate = append(listDate, start.Format(layout))
				if start.Format(layout) == endDate {
					break
				}
			}
		}
		return db.Where("date in ?", listDate)
	}
}

type InputIsExists struct {
	VisitID       string `gorm:"column:visit_id"`
	Time          string `gorm:"column:time"`
	TrafficSource string `gorm:"column:traffic_source"`
	Campaign      string `gorm:"column:campaign"`
	SelectionID   string `gorm:"column:selection_id"`
}

func (t *reportAffPixelRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportAffPixelModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportAffPixelRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportAffPixelModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportAffPixelRepo) FindOneByQuery(query map[string]interface{}, queryRaw ...string) (record *model.ReportAffPixelModel, err error) {
	err = t.Db.
		//Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(queryRaw) > 0 {
				for _, raw := range queryRaw {
					db.Where(raw)
				}
			}
			return db
		}).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportAffPixelRepo) FindAllByQuery(query map[string]interface{}) (records []*model.ReportAffPixelModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&records).Error
	return
}

func (t *reportAffPixelRepo) FindForReportAff(day string) (records []*model.ReportAffPixelModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Select(
			"id",
			"uid",
			"traffic_source",
			"demand_source",
			"campaign_id",
			"redirect_id",
			"user_id",
			"style_id",
			"layout_id",
			"layout_version",
			"device",
			"geo",
			"section_id",
			"section_name",
			"SUM(imp_quantumdex) AS imp_quantumdex",
			"SUM(system_traffic) AS system_traffic",
			"SUM(pre_bot_traffic) AS pre_bot_traffic",
			"SUM(bot_traffic) AS bot_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
			"SUM(click2) AS click2",
			"SUM(click3) AS click3",
			"SUM(click4) AS click4",
			"SUM(estimate_revenue) AS estimate_revenue",
			"SUM(pre_estimate_revenue) AS pre_estimate_revenue",
			"SUM(cost) AS cost",
		).
		Where("time_conversion LIKE '%" + day + "%'").
		//Where("traffic_source LIKE 'quantumdex'").
		Group("traffic_source,redirect_id,section_id,style_id,layout_id,layout_version,device,geo").
		Find(&records).Error
	return
}

func (t *reportAffPixelRepo) FindForReportAffSelection(day string) (records []*model.ReportAffPixelModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Select(
			"id",
			"uid",
			"traffic_source",
			"demand_source",
			"campaign_id",
			"section_id",
			"section_name",
			"publisher_id",
			"publisher_name",
			"SUM(system_traffic) AS system_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
		).
		Where("time_conversion LIKE '%" + day + "%'").
		Group("traffic_source,campaign_id,section_id,publisher_id").
		Find(&records).Error
	return
}

func (t *reportAffPixelRepo) FindAllByGroup(group ...string) (records []*model.ReportAffPixelModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Group(strings.Join(group, ",")).
		Find(&records).Error
	return
}

func (t *reportAffPixelRepo) Save(record *model.ReportAffPixelModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportAffPixelRepo) SaveSlice(records []*model.ReportAffPixelModel) (err error) {
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportAffPixelRepo) FindCampaignIDForAdsense(redirectID int64) (campaignID string, err error) {
	err = t.Db.
		//Debug().
		Model(&model.ReportAffPixelModel{}).
		Select("campaign_id").
		Where("redirect_id = ?", redirectID).
		Where("campaign_id != 'unknow' and campaign_id != ''").
		Last(&campaignID).Error
	return
}

func (t *reportAffPixelRepo) FindRedirectID(campaignID string) (redirectID string, err error) {
	err = t.Db.
		//Debug().
		Model(&model.ReportAffPixelModel{}).
		Select("redirect_id").
		Where("campaign_id = ?", campaignID).
		Where("time_conversion < '2022-12-01'").
		Where("redirect_id != 'unknow' and redirect_id != ''").
		Last(&redirectID).Error
	return
}

func (t *reportAffPixelRepo) FindForReportAllTrafficSourceIDs(trafficSource, account string) (ids []string, err error) {
	err = t.Db.
		//Debug().
		Model(&model.ReportAffPixelModel{}).
		Select("campaign_id").
		Where("traffic_source = ?", trafficSource).
		Where("account = ?", account).
		Where("campaign_id != 'unknow' and campaign_id != ''").
		Where("time_conversion > ?", time.Now().UTC().AddDate(0, 0, -2).Format("2006-01-02 15:04:05")).
		Group("campaign_id").
		Find(&ids).Error
	return
}
