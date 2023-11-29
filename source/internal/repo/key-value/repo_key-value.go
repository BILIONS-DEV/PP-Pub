package key_value

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
)

type RepoKeyValue interface {
	Migrate()
	Filter(input *InputFilter) (totalRecord int64, records []*model.KeyModel, err error)
	FindByID(ID int64) (record *model.KeyModel, err error)
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	Save(record *model.KeyModel) (err error)
	Delete(ID int64) (err error)
}

type keyValueRepo struct {
	Db *gorm.DB
}

func NewKeyValueRepo(DB *gorm.DB) *keyValueRepo {
	return &keyValueRepo{Db: DB}
}

type InputIsExists struct {
	Key   string
	Value string
}

func (t *keyValueRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
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
	var record model.KeyModel
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

func (t *keyValueRepo) Filter(input *InputFilter) (totalRecord int64, records []*model.KeyModel, err error) {
	var table = model.KeyModel{}.TableName()
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

func (t *keyValueRepo) setFilterQuerySearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}
		return db.Where("key_name LIKE ?", "%"+search+"%")
	}
}

func (t *keyValueRepo) FindByID(ID int64) (record *model.KeyModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		First(&record, ID).Error
	return
}

func (t *keyValueRepo) Save(record *model.KeyModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *keyValueRepo) Delete(ID int64) (err error) {
	err = t.Db.
		//Debug().
		Select("Value").
		Delete(&model.KeyModel{ID: ID}).Error
	return
}

func (t *keyValueRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.Debug().AutoMigrate(
		&model.KeyModel{},
		&model.ValueModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}
