package account

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/internal/entity/model"
)

type RepoAccount interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByID(ID int64) (record *model.AccountModel, err error)
	FindByQuery(query map[string]interface{}) (record *model.AccountModel, err error)
	Save(record *model.AccountModel) (err error)
	SaveSlice(records []*model.AccountModel) (err error)
}

type accountRepo struct {
	Db *gorm.DB
}

func NewAccountRepo(DB *gorm.DB) *accountRepo {
	return &accountRepo{Db: DB}
}

type InputIsExists struct {
	UtcHour  string
	Campaign string
	SubID    string
}

func (t *accountRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if input.UtcHour != "" {
					return db.Where("utc_hour = ?", input.UtcHour)
				}
				return db
			},
			func(db *gorm.DB) *gorm.DB {
				if input.Campaign != "" {
					return db.Where("campaign = ?", input.Campaign)
				}
				return db
			},
			func(db *gorm.DB) *gorm.DB {
				if input.SubID != "" {
					return db.Where("sub_id = ?", input.SubID)
				}
				return db
			},
		).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.AccountModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *accountRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.AccountModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *accountRepo) FindByID(ID int64) (record *model.AccountModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		First(&record, ID).Error
	return
}

func (t *accountRepo) FindByQuery(query map[string]interface{}) (record *model.AccountModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		First(&record).Error
	return
}

func (t *accountRepo) Save(record *model.AccountModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *accountRepo) SaveSlice(records []*model.AccountModel) (err error) {
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}
