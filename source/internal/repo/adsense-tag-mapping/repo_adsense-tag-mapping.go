package adsense_tag_mapping

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
)

type RepoAdsenseTagMapping interface {
	Migrate()
	Filter(input *InputFilter) (totalRecord int64, records []*model.AdsenseTagMappingModel, err error)
	FindByID(ID int64) (record *model.AdsenseTagMappingModel, err error)
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	Save(record *model.AdsenseTagMappingModel) (err error)
	Delete(ID int64) (err error)
}

type adsenseTagMappingRepo struct {
	Db *gorm.DB
}

func NewAdsenseTagMappingRepo(DB *gorm.DB) *adsenseTagMappingRepo {
	return &adsenseTagMappingRepo{Db: DB}
}

type InputIsExists struct {
	Key   string
	Value string
}

func (t *adsenseTagMappingRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Scopes(func(db *gorm.DB) *gorm.DB {
			if input.Value != "" {
				return db.Joins("Value", t.Db.Where("custom_value.name = ?", input.Value))
			}
			return db
		}).
		Select("ID").
		Where("custom_key.name = ?", input.Key)
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.AdsenseTagMappingModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

type InputFilter struct {
	UserID int64

	Search string
	Offset int
	Limit  int
	Order  string
}

func (t *adsenseTagMappingRepo) Filter(input *InputFilter) (totalRecord int64, records []*model.AdsenseTagMappingModel, err error) {
	var table = model.AdsenseTagMappingModel{}.TableName()
	err = t.Db.
		//Debug().
		Scopes(
			t.setFilterQuerySearch(input.Search),
		).
		Model(&records).Distinct(table + `.id`).Count(&totalRecord).
		Order(input.Order).
		Scopes(mysql.Paginate(mysql.Deps{Offset: input.Offset, Limit: input.Limit})).
		Select("*").Group(table + ".id").Find(&records).Error
	return
}

func (t *adsenseTagMappingRepo) setFilterQuerySearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}
		return db.Where("key_name LIKE ?", "%"+search+"%")
	}
}

func (t *adsenseTagMappingRepo) FindByID(ID int64) (record *model.AdsenseTagMappingModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		First(&record, ID).Error
	return
}

func (t *adsenseTagMappingRepo) Save(record *model.AdsenseTagMappingModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *adsenseTagMappingRepo) Delete(ID int64) (err error) {
	err = t.Db.
		//Debug().
		Select("Value").
		Delete(&model.AdsenseTagMappingModel{ID: ID}).Error
	return
}

func (t *adsenseTagMappingRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.Debug().AutoMigrate(
		&model.AdsenseTagMappingModel{},
		&model.ValueModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}
