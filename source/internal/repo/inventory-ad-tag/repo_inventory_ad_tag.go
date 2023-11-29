package inventory_ad_tag

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
	// "source/core/technology/mysql"
)

type RepoInventoryAdTag interface {
	FindByID(ID int64) (record *model.InventoryAdTagModel, err error)
}

type inventoryAdTagRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewInventoryAdTagRepo(db *gorm.DB, cache caching.Cache) *inventoryAdTagRepo {
	return &inventoryAdTagRepo{Db: db, Cache: cache}
}

func (t *inventoryAdTagRepo) FindByID(ID int64) (record *model.InventoryAdTagModel, err error) {
	err = t.Db.
		// Debug().
		First(&record, ID).Error
	return
}
