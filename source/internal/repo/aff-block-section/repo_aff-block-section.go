package aff_block_section

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
)

type RepoAffBlockSection interface {
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByID(ID int64) (record *model.AffBlockSectionModel, err error)
	FindOneByQuery(query map[string]interface{}) (record *model.AffBlockSectionModel, err error)
	FindAllByQuery(query map[string]interface{}) (records []*model.AffBlockSectionModel, err error)
	Save(record *model.AffBlockSectionModel) (err error)
	SaveSlice(records []*model.AffBlockSectionModel) (err error)
	DeleteByCampaignID(campaignID string) (err error)
}

type affBlockSectionRepo struct {
	Db    *gorm.DB
	Kafka *kafka.Client
}

func NewAffBlockSectionRepo(DB *gorm.DB) *affBlockSectionRepo {
	return &affBlockSectionRepo{Db: DB}
}

type InputIsExists struct {
	UtcHour  string
	Campaign string
	SubID    string
}

func (t *affBlockSectionRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
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
	var record model.AffBlockSectionModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *affBlockSectionRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.AffBlockSectionModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *affBlockSectionRepo) FindByID(ID int64) (record *model.AffBlockSectionModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		First(&record, ID).Error
	return
}

func (t *affBlockSectionRepo) FindOneByQuery(query map[string]interface{}) (record *model.AffBlockSectionModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *affBlockSectionRepo) FindAllByQuery(query map[string]interface{}) (records []*model.AffBlockSectionModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&records).Error
	return
}

func (t *affBlockSectionRepo) Save(record *model.AffBlockSectionModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *affBlockSectionRepo) SaveSlice(records []*model.AffBlockSectionModel) (err error) {
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *affBlockSectionRepo) DeleteByCampaignID(campaignID string) (err error) {
	err = t.Db.
		//Debug().
		Where("campaign_id = ?", campaignID).
		Delete(&model.AffBlockSectionModel{}).Error
	return
}
