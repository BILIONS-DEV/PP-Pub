package reportAff

import (
	"fmt"
	"github.com/google/uuid"
	"source/infrastructure/mysql"
	"strings"
	"time"

	// "source/infrastructure/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

type RepoReportAff interface {
	Migrate()
	Filter(input *InputFilter) (totalRecord int64, records []*model.ReportAffModel, recordTotal *model.ReportAffModel, err error)
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindOneByQuery(query map[string]interface{}) (record *model.ReportAffModel, err error)
	FindAllByQuery(query map[string]interface{}) (records []*model.ReportAffModel, err error)
	FindAllByGroup(group ...string) (records []*model.ReportAffModel, err error)
	Save(record *model.ReportAffModel) (err error)
	SaveSlice(records []*model.ReportAffModel) (err error)
	PushCache(record model.PixelAff) (err error)
	FindAllCache() (records []model.PixelAffForAll, err error)
	PushCacheCheck(key string, record model.PixelAffCheck) (err error)
	FindCacheCheck(key string) (record model.PixelAffCheck, err error)
	PushBinCache(setName, key, name string, value interface{}) (err error)
	DeleteCache(key string) (err error)
	UpdateCampaignForReport(campaignID string, campaignName string, userID int64) (err error)
	FindForBlockSectionOutBrain(endDate string) (records []*model.ReportAffModel, err error)
	FindReportTotalForBlock(campaignID, endDate string) (record *model.ReportAffModel, err error)
	Delete(record *model.ReportAffModel) (err error)
}

type reportAffRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportAffRepo(db *gorm.DB, cache caching.Cache, ka *kafka.Client) *reportAffRepo {
	return &reportAffRepo{Db: db, Cache: cache, Kafka: ka}
}

type InputFilter struct {
	UserID int64

	Search string
	Offset int
	Limit  int
	Order  string

	StartDate      string
	EndDate        string
	Campaigns      interface{}
	SectionID      interface{}
	TrafficSources interface{}
	DemandSources  interface{}
	GroupBy        interface{}
}

func (t *reportAffRepo) Filter(input *InputFilter) (totalRecord int64, records []*model.ReportAffModel, recordTotal *model.ReportAffModel, err error) {
	err = t.Db.
		// Debug().
		Select(
			"date",
			"traffic_source",
			"section_id", "section_name",
			"campaign_id",
			"campaign_name",
			"partner",
			"SUM(system_traffic) AS system_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
			"SUM(click_adsense) AS click_adsense",
			"SUM(cost) AS cost",
			"SUM(revenue) AS revenue",
			"SUM(estimated_revenue) AS estimated_revenue",
			"SUM(pre_estimated_revenue) AS pre_estimated_revenue",
			"SUM(supply_click) AS supply_click",
			"SUM(supply_conversions) AS supply_conversions",
		).
		Where("campaign_id != 'unknow'").
		Scopes(
			t.setFilterCondition(input),
			t.setFilterGroupBy(input.GroupBy),
			t.setFilterDate(input.StartDate, input.EndDate),
		).
		Model(&model.ReportAffModel{}).
		Count(&totalRecord).
		Order(input.Order).
		Scopes(mysql.Paginate(mysql.Deps{Offset: input.Offset, Limit: input.Limit})).
		Find(&records).Error
	err = t.Db.
		//Debug().
		Select(
			"SUM(system_traffic) AS system_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
			"SUM(click_adsense) AS click_adsense",
			"SUM(cost) AS cost",
			"SUM(revenue) AS revenue",
			"SUM(supply_click) AS supply_click",
			"SUM(supply_conversions) AS supply_conversions",
			"SUM(estimated_revenue) AS estimated_revenue",
			"SUM(pre_estimated_revenue) AS pre_estimated_revenue",
		).
		Where("campaign_id != 'unknow'").
		Scopes(
			t.setFilterCondition(input),
			t.setFilterDate(input.StartDate, input.EndDate),
		).
		Model(&model.ReportAffModel{}).
		Find(&recordTotal).Error
	return
}

func (t *reportAffRepo) setFilterGroupBy(groupBy interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if groupBy == nil {
			return db
		}
		switch groupBy.(type) {
		case string:
			if groupBy != "" {
				group := ""
				if fmt.Sprintf("%v", groupBy) == "campaign" {
					group = "campaign_id, traffic_source"
				} else if fmt.Sprintf("%v", groupBy) == "section_id" {
					group = "section_id, traffic_source"
				} else {
					group = fmt.Sprintf("%v", groupBy)
				}
				return db.Group(group)
			}
		case []interface{}:
			listGroup := groupBy.([]interface{})
			var groups []string
			var flag bool
			for _, group := range listGroup {
				if fmt.Sprintf("%v", group) == "campaign" {
					groups = append(groups, fmt.Sprintf("%v", "campaign_id"))
					if !flag {
						groups = append(groups, fmt.Sprintf("%v", "traffic_source"))
						flag = true
					}
				} else if fmt.Sprintf("%v", group) == "section_id" {
					groups = append(groups, fmt.Sprintf("%v", "section_id"))
					if !flag {
						groups = append(groups, fmt.Sprintf("%v", "traffic_source"))
						flag = true
					}
				} else {
					groups = append(groups, fmt.Sprintf("%v", group))
				}
			}
			return db.Group(strings.Join(groups, ","))
		}
		return db
	}
}

func (t *reportAffRepo) setFilterCondition(input *InputFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var condition = make(map[string]interface{})
		// => convert input sang dữ liệu condition để query
		if input.Campaigns != nil {
			condition["campaign_id"] = input.Campaigns
		}
		if input.TrafficSources != nil {
			condition["traffic_source"] = input.TrafficSources
		}
		if input.DemandSources != nil {
			condition["partner"] = input.DemandSources
		}
		if input.SectionID != nil {
			condition["section_id"] = input.SectionID
		}
		return db.Where(condition)
	}
}

func (t *reportAffRepo) setFilterDate(startDate, endDate string) func(db *gorm.DB) *gorm.DB {
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

func (t *reportAffRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		// Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportAffModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportAffRepo) Migrate() {
	// os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportAffModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportAffRepo) FindOneByQuery(query map[string]interface{}) (record *model.ReportAffModel, err error) {
	err = t.Db.
		//Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		First(&record).Error
	return
}

func (t *reportAffRepo) FindAllByQuery(query map[string]interface{}) (records []*model.ReportAffModel, err error) {
	err = t.Db.
		// Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&records).Error
	return
}

func (t *reportAffRepo) FindAllByGroup(group ...string) (records []*model.ReportAffModel, err error) {
	err = t.Db.
		// Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Group(strings.Join(group, ",")).
		Find(&records).Error
	return
}

func (t *reportAffRepo) Save(record *model.ReportAffModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportAffRepo) Delete(record *model.ReportAffModel) (err error) {
	err = t.Db.
		//Debug().
		Delete(record).Error
	return
}

func (t *reportAffRepo) UpdateCampaignForReport(campaignID string, campaignName string, userID int64) (err error) {
	err = t.Db.
		// Debug().
		Model(&model.ReportAffModel{}).
		Where("campaign_id = ?", campaignID).
		Updates(model.ReportAffModel{
			UserID:       userID,
			CampaignName: campaignName,
		}).
		Error
	return
}

func (t *reportAffRepo) SaveSlice(records []*model.ReportAffModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		// Debug().
		Save(&records).Error
	return
}

func (t *reportAffRepo) PushCache(record model.PixelAff) (err error) {
	key := uuid.NewString()
	return t.Cache.SetWithTTL(key, record, 3600, record.SetName())
}

func (t *reportAffRepo) FindAllCache() (records []model.PixelAffForAll, err error) {
	err = t.Cache.GetAll(&records, model.PixelAff{}.SetName())
	return
}

func (t *reportAffRepo) PushCacheCheck(key string, record model.PixelAffCheck) (err error) {
	return t.Cache.SetWithTTL(key, record, 36000, record.SetName())
}

func (t *reportAffRepo) PushBinCache(setName, key, name string, value interface{}) (err error) {
	return t.Cache.SetBinsWithTTL(setName, key, 36000, &caching.Bin{
		Name:  name,
		Value: value,
	})
}

func (t *reportAffRepo) FindCacheCheck(key string) (record model.PixelAffCheck, err error) {
	err = t.Cache.Get(key, &record, record.SetName())
	return
}

func (t *reportAffRepo) DeleteCache(key string) (err error) {
	err = t.Cache.Delete(key, model.PixelAff{}.SetName())
	return
}

func (t *reportAffRepo) FindForBlockSectionOutBrain(endDate string) (records []*model.ReportAffModel, err error) {
	err = t.Db.
		// Debug().
		Select(
			"date",
			"traffic_source",
			"section_id",
			"section_name",
			"campaign_id",
			"campaign_name",
			"style_id",
			"partner",
			"SUM(system_traffic) AS system_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
			"SUM(click_adsense) AS click_adsense",
			"SUM(cost) AS cost",
			"SUM(revenue) AS revenue",
			"SUM(estimated_revenue) AS estimated_revenue",
			"SUM(pre_estimated_revenue) AS pre_estimated_revenue",
			"SUM(supply_click) AS supply_click",
			"SUM(supply_conversions) AS supply_conversions",
		).
		Where("traffic_source like 'outbrain'").
		Where("DATE_FORMAT('2017-06-15', date) > DATE_FORMAT('2017-06-15', '2022-09-29')").  // Bắt đầu từ ngày 29/09/2022
		Where("DATE_FORMAT('2017-06-15', date) <= DATE_FORMAT('2017-06-15', '?')", endDate). // Ngày kết thúc lấy report
		Group("campaign_id, section_id").
		Find(&records).
		Error
	return
}

func (t *reportAffRepo) FindReportTotalForBlock(campaignID, endDate string) (record *model.ReportAffModel, err error) {
	err = t.Db.
		//Debug().
		Select(
			"date",
			"traffic_source",
			"section_id",
			"section_name",
			"campaign_id",
			"campaign_name",
			"style_id",
			"partner",
			"SUM(system_traffic) AS system_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
			"SUM(click_adsense) AS click_adsense",
			"SUM(cost) AS cost",
			"SUM(revenue) AS revenue",
			"SUM(supply_click) AS supply_click",
			"SUM(supply_conversions) AS supply_conversions",
			"SUM(estimated_revenue) AS estimated_revenue",
			"SUM(pre_estimated_revenue) AS pre_estimated_revenue",
		).
		Where("campaign_id = ?", campaignID).
		Where("DATE_FORMAT('2017-06-15', date) > DATE_FORMAT('2017-06-15', '2022-09-29')").  // Bắt đầu từ ngày 29/09/2022
		Where("DATE_FORMAT('2017-06-15', date) <= DATE_FORMAT('2017-06-15', '?')", endDate). // Ngày kết thúc lấy report
		Group("campaign_id").
		Find(&record).
		Error
	return
}
