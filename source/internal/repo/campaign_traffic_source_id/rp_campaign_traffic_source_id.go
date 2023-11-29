package campaign_traffic_source_id

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/internal/entity/model"
)

const KeyName = "RepoCampaignTrafficSourceID"

type RepoCampaignTrafficSourceID interface {
	FindByID(ID int64) (record *model.CampaignTrafficSourceID, err error)
	FindByQuery(query map[string]interface{}) (record *model.CampaignTrafficSourceID, err error)
	Save(record *model.CampaignTrafficSourceID) (err error)
	SaveSlice(records []*model.CampaignTrafficSourceID) (err error)
	FindAllByQuery(query model.CampaignTrafficSourceID, queryRaw ...string) (record []*model.CampaignTrafficSourceID, err error)
	FindOneByQuery(query model.CampaignTrafficSourceID, queryRaw ...string) (record *model.CampaignTrafficSourceID, err error)
}

type campaignTrafficSourceIDRP struct {
	DB *gorm.DB
}

func NewCampaignTrafficSourceIDRP(DB *gorm.DB) *campaignTrafficSourceIDRP {
	return &campaignTrafficSourceIDRP{DB: DB}
}

func (t *campaignTrafficSourceIDRP) Migrate() {
	// os.Exit(1)
	err := t.DB.AutoMigrate(
		&model.CampaignTrafficSourceID{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *campaignTrafficSourceIDRP) FindByID(ID int64) (record *model.CampaignTrafficSourceID, err error) {
	err = t.DB.
		// Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		First(&record, ID).Error
	return
}

func (t *campaignTrafficSourceIDRP) FindByName(name string) (record *model.CampaignTrafficSourceID) {
	t.DB.
		// Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where("name = ?", name).
		Find(&record)
	return
}

func (t *campaignTrafficSourceIDRP) FindByQuery(query map[string]interface{}) (record *model.CampaignTrafficSourceID, err error) {
	err = t.DB.
		// Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		First(&record).Error
	return
}

func (t *campaignTrafficSourceIDRP) Save(record *model.CampaignTrafficSourceID) (err error) {
	err = t.DB.
		// Debug().
		Save(record).Error
	return
}

func (t *campaignTrafficSourceIDRP) SaveSlice(records []*model.CampaignTrafficSourceID) (err error) {
	err = t.DB.
		// Debug().
		Save(&records).Error
	return
}

func (t *campaignTrafficSourceIDRP) FindAllByQuery(query model.CampaignTrafficSourceID, queryRaw ...string) (record []*model.CampaignTrafficSourceID, err error) {
	err = t.DB.
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
		Find(&record).Error
	return
}

func (t *campaignTrafficSourceIDRP) FindOneByQuery(query model.CampaignTrafficSourceID, queryRaw ...string) (record *model.CampaignTrafficSourceID, err error) {
	err = t.DB.
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
