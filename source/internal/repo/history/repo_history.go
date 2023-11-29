package history

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
	"time"
)

type history struct {
	Db    *gorm.DB
	Cache caching.Cache
}

type RepoHistory interface {
	FindByUser(userID int64, input *FilterByUserInput) (totalRecord int64, histories []model.History, err error)
	FindByFilter(input *FilterInput) (totalRecord int64, histories []model.History, err error)
	FindByID(ID, userID int64) (record model.History)
	GetByID(ID int64) (record model.History)
	GetAllHistoryPage() (pages []string)
	Save(record *model.HistoryModel) (err error)
	LoadHistoriesByObject(input *InputObject) (histories []model.History, err error)
}

func NewHistoryRepo(db *gorm.DB, cache caching.Cache) *history {
	return &history{Db: db, Cache: cache}
}

func (t *history) Save(record *model.HistoryModel) (err error) {
	return
}

func (t *history) FindByID(ID, userID int64) (record model.History) {
	t.Db.Where("id = ? AND user_id = ?", ID, userID).Take(&record)
	return
}

func (t *history) GetByID(ID int64) (record model.History) {
	t.Db.Where("id = ?", ID).Take(&record)
	return
}

type FilterInput struct {
	// filter
	Condition  map[string]interface{}
	StartDate  time.Time
	EndDate    time.Time
	Search     string
	ObjectPage string
	ObjectID   string
	// order and pagination
	Order string
	Page  int
}

func (t *history) FindByFilter(input *FilterInput) (totalRecord int64, histories []model.History, err error) {
	t.Db.
		// Debug().
		// Where("user_id = ? AND app = 'FE' AND creator_id = ?", userID, userID).
		Where(input.Condition).
		Scopes(
			t.setFilterQuerySearch(input.Search),
			t.setFilterQueryObjectId(input.ObjectID),
			t.setFilterQueryObjectPage(input.ObjectPage),
			t.setFilterQueryDate(input.StartDate, input.EndDate),
		).
		Model(&histories).Select("id").Count(&totalRecord).
		Order(input.Order).
		Scopes(mysql.Paginate(mysql.Deps{Page: input.Page})).
		Select("*").Find(&histories)
	return
}

type FilterByUserInput struct {
	// filter
	Condition map[string]interface{}
	StartDate time.Time
	EndDate   time.Time
	Search    string
	// order and pagination
	Order string
	Page  int
}

func (t *history) FindByUser(userID int64, input *FilterByUserInput) (totalRecord int64, histories []model.History, err error) {
	t.Db.
		// Debug().
		Where("user_id = ? AND app = 'FE' AND creator_id = ?", userID, userID).
		Where(input.Condition).
		Scopes(
			t.setFilterQuerySearch(input.Search),
			t.setFilterQueryDate(input.StartDate, input.EndDate),
		).
		Model(&histories).Select("id").Count(&totalRecord).
		Order(input.Order).
		Scopes(mysql.Paginate(mysql.Deps{Page: input.Page})).
		Select("*").Find(&histories)
	return
}

func (t *history) setFilterQuerySearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}
		return db.Where("title LIKE ? OR object_name like ? OR new_data like ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
}

func (t *history) setFilterQueryDate(startDate, endDate time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if startDate.IsZero() || endDate.IsZero() {
			return db
		}
		return db.Where("created_at >= ? AND created_at <= ?", startDate, endDate.AddDate(0, 0, 1))
	}
}

func (t *history) setFilterQueryObjectPage(objectPage string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if objectPage == "" {
			return db
		}
		return db.Where("page = ?", objectPage)
	}
}

func (t *history) setFilterQueryObjectId(objectId string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if objectId == "" {
			return db
		}
		return db.Where("object_id = ?", objectId)
	}
}

func (t *history) GetAllHistoryPage() (pages []string) {
	var histories []model.History
	t.Db.Group("page").Find(&histories)
	if len(histories) == 0 {
		return
	}

	for _, history := range histories {
		pages = append(pages, history.Page)
	}
	return
}

type InputObject struct {
	Id         string
	ObjectType string
}

func (t *history) LoadHistoriesByObject(input *InputObject) (histories []model.History, err error) {
	err = t.Db.Model(model.History{}).Where("object_id = ? and detail_type like ?", input.Id, input.ObjectType+"%").
		Order("created_at DESC").
		Find(&histories).Error
	return
}
