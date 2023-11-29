package gam

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
)

type RepoGam interface {
	FindByID(ID int64) (record *model.GamModel, err error)
}

type gamRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewGamRepo(db *gorm.DB, cache caching.Cache) *gamRepo {
	return &gamRepo{Db: db, Cache: cache}
}

func (t *gamRepo) FindByID(ID int64) (record *model.GamModel, err error) {
	err = t.Db.
		First(&record, ID).Error
	return
}
