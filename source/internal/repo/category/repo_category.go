package category

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
)

type category struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func (t *category) FindAll() (records []model.Category, err error) {
	err = t.Db.Find(&records).Error
	return
}

func NewCategoryRepo(db *gorm.DB, cache caching.Cache) *category {
	return &category{Db: db, Cache: cache}
}
