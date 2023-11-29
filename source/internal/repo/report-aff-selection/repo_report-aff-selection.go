package reportAffSelectionSelection

import (
	"fmt"
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

type RepoReportAffSelection interface {
	Migrate()
	Filter(input *InputFilter) (totalRecord int64, records []*model.ReportAffSelectionModel, recordTotal *model.ReportAffSelectionModel, err error)
	FindOneByQuery(query map[string]interface{}) (record *model.ReportAffSelectionModel, err error)
	FindAllByQuery(query map[string]interface{}) (records []*model.ReportAffSelectionModel, err error)
	FindAllByGroup(group ...string) (records []*model.ReportAffSelectionModel, err error)
	Save(record *model.ReportAffSelectionModel) (err error)
	SaveSlice(records []*model.ReportAffSelectionModel) (err error)
	UpdateCampaignName(campaignID string, campaignName string) (err error)
}

type reportAffSelectionRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
}

func NewReportAffSelectionRepo(db *gorm.DB, cache caching.Cache, ka *kafka.Client) *reportAffSelectionRepo {
	return &reportAffSelectionRepo{Db: db, Cache: cache, Kafka: ka}
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
	TrafficSources interface{}
	PublisherID    interface{}
	GroupBy        interface{}
	SectionID      interface{}
}

func (t *reportAffSelectionRepo) Filter(input *InputFilter) (totalRecord int64, records []*model.ReportAffSelectionModel, recordTotal *model.ReportAffSelectionModel, err error) {
	err = t.Db.
		//Debug().
		Select(
			"date",
			"traffic_source",
			"section_id", "section_name",
			"campaign_id",
			"campaign_name",
			"partner",
			"publisher_id",
			"SUM(system_traffic) AS system_traffic",
			"SUM(impressions) AS impressions",
			"SUM(click) AS click",
			"SUM(revenue) AS revenue",
		).
		Scopes(
			t.setFilterCondition(input.Campaigns, input.TrafficSources, input.PublisherID, input.SectionID),
			t.setFilterGroupBy(input.GroupBy),
			t.setFilterDate(input.StartDate, input.EndDate),
		).
		Model(&model.ReportAffSelectionModel{}).
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
			"SUM(revenue) AS revenue",
		).
		Scopes(
			t.setFilterCondition(input.Campaigns, input.TrafficSources, input.PublisherID, input.SectionID),
			t.setFilterDate(input.StartDate, input.EndDate),
		).
		Model(&model.ReportAffSelectionModel{}).
		Find(&recordTotal).Error
	return
}

func (t *reportAffSelectionRepo) setFilterGroupBy(groupBy interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if groupBy == nil {
			return db
		}
		switch groupBy.(type) {
		case string:
			if groupBy != "" {
				group := ""
				if fmt.Sprintf("%v", groupBy) == "campaign" {
					group = "campaign_id"
				} else if fmt.Sprintf("%v", groupBy) == "publisher" {
					group = "publisher_id"
				} else {
					group = fmt.Sprintf("%v", groupBy)
				}
				return db.Group(group)
			}
		case []interface{}:
			listGroup := groupBy.([]interface{})
			var groups []string
			for _, group := range listGroup {
				if fmt.Sprintf("%v", group) == "campaign" {
					groups = append(groups, fmt.Sprintf("%v", "campaign_id"))
				} else if fmt.Sprintf("%v", group) == "publisher" {
					groups = append(groups, fmt.Sprintf("%v", "publisher_id"))
				} else {
					groups = append(groups, fmt.Sprintf("%v", group))
				}
			}
			return db.Group(strings.Join(groups, ","))
		}
		return db
	}
}

func (t *reportAffSelectionRepo) setFilterCondition(campaigns interface{}, trafficSources interface{}, publisherID interface{}, sectionID interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var condition = make(map[string]interface{})
		//=> convert input sang dữ liệu condition để query
		if campaigns != nil {
			condition["campaign_id"] = campaigns
		}
		if trafficSources != nil {
			condition["traffic_source"] = trafficSources
		}
		if publisherID != nil {
			condition["publisher_id"] = publisherID
		}
		if sectionID != nil {
			condition["section_id"] = sectionID
		}
		return db.Where(condition)
	}
}

func (t *reportAffSelectionRepo) setFilterDate(startDate, endDate string) func(db *gorm.DB) *gorm.DB {
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

func (t *reportAffSelectionRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportAffSelectionModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportAffSelectionRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportAffSelectionModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportAffSelectionRepo) FindOneByQuery(query map[string]interface{}) (record *model.ReportAffSelectionModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		First(&record).Error
	return
}

func (t *reportAffSelectionRepo) FindAllByQuery(query map[string]interface{}) (records []*model.ReportAffSelectionModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&records).Error
	return
}

func (t *reportAffSelectionRepo) FindAllByGroup(group ...string) (records []*model.ReportAffSelectionModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Group(strings.Join(group, ",")).
		Find(&records).Error
	return
}

func (t *reportAffSelectionRepo) Save(record *model.ReportAffSelectionModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportAffSelectionRepo) SaveSlice(records []*model.ReportAffSelectionModel) (err error) {
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportAffSelectionRepo) UpdateCampaignName(campaignID string, campaignName string) (err error) {
	err = t.Db.
		//Debug().
		Model(&model.ReportAffSelectionModel{}).
		Where("campaign_id = ?", campaignID).
		Update("campaign_name", campaignName).
		Error
	return
}
