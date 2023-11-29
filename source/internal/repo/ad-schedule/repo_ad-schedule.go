package ad_schedule

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
)

type RepoAdSchedule interface {
	Migrate()
	DB() *gorm.DB
	Filter(input *InputFilter) (totalRecord int64, records []*model.AdScheduleModel, err error)
	FindByID(ID int64) (record *model.AdScheduleModel, err error)
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	Save(record *model.AdScheduleModel) (err error)
	EmptyConfigs(configs []model.AdScheduleConfigModel) (err error)
	DeleteByID(ID int64) error
}

type adScheduleRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewAdScheduleRepo(DB *gorm.DB, cache caching.Cache) *adScheduleRepo {
	return &adScheduleRepo{Db: DB, Cache: cache}
}

type InputIsExists struct {
	Name string
}

func (t *adScheduleRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		Debug().
		Select("ID").
		Where("name = ?", input.Name)
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.AdScheduleModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

type InputFilter struct {
	UserID                int64
	ClientTypes           []string
	AdBreakTypes          []string
	AdBreakConfigTypeJoin []string
	Search                string
	Page                  int
	Order                 string
}

func (t *adScheduleRepo) Filter(input *InputFilter) (totalRecord int64, records []*model.AdScheduleModel, err error) {
	var table = model.AdScheduleModel{}.TableName()
	err = t.Db.
		Debug().
		//Where(input.Condition).
		//Joins(input.JoinTable, t.Db.Where(input.JoinCondition)).
		Scopes(
			t.setFilterCondition(input.UserID, input.ClientTypes, input.AdBreakTypes),
			t.setFilterQuerySearch(input.Search),
			t.setFilterJoin(input.AdBreakConfigTypeJoin),
		).
		Model(&records).Distinct(table + `.id`).Count(&totalRecord).
		Order(input.Order).
		Scopes(mysql.Paginate(mysql.Deps{Page: input.Page, Limit: 10})).
		Preload("AdBreakConfigs").
		Preload("AdBreakConfigs.AdTagUrls").
		Select("*").Group(table + ".id").Find(&records).Error
	return
}

func (t *adScheduleRepo) setFilterCondition(userID int64, listClientType, listAdBreakType []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var condition = make(map[string]interface{})
		//=> convert input sang dữ liệu condition để query
		condition["user_id"] = userID
		if len(listClientType) > 0 {
			condition["client_type"] = listClientType
		}
		if len(listAdBreakType) > 0 {
			condition["ad_break_type"] = listAdBreakType
		}
		return db.Where(condition)
	}
}

func (t *adScheduleRepo) setFilterJoin(listAdBreakConfigType []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//=> convert input sang dữ liệu join
		var joinTable string
		var joinCondition = make(map[string]interface{})
		if len(listAdBreakConfigType) > 0 {
			joinTable = "AdBreakConfigs"
			joinCondition["AdBreakConfigs.config_type"] = listAdBreakConfigType
		}
		return db.Joins(joinTable, t.Db.Where(joinCondition))
	}
}

func (t *adScheduleRepo) setFilterQuerySearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}
		return db.Where("title LIKE ?", "%"+search+"%")
	}
}

func (t *adScheduleRepo) FindByID(ID int64) (record *model.AdScheduleModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload("AdBreakConfigs").
		Preload("AdBreakConfigs.AdTagUrls").
		First(&record, ID).Error
	return
}

func (t *adScheduleRepo) Save(record *model.AdScheduleModel) (err error) {
	err = t.Db.
		Debug().
		Save(record).Error
	return
}

func (t *adScheduleRepo) EmptyConfigs(configs []model.AdScheduleConfigModel) (err error) {
	err = t.Db.
		//Debug().
		Select("AdTagUrls").Delete(&configs).Error
	return
}

func (t *adScheduleRepo) DeleteByID(ID int64) error {
	return t.Db.Delete(&model.AdScheduleModel{}, ID).Error
}

func (t *adScheduleRepo) DB() *gorm.DB {
	return t.Db
}
func (t *adScheduleRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.Debug().AutoMigrate(
		&model.AdScheduleModel{},
		&model.AdScheduleConfigModel{},
		&model.AdScheduleConfigAdTagUrlModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}
