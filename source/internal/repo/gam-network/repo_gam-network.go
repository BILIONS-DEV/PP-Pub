package gam_network

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
)

type RepoGamNetwork interface {
	FindByID(ID int64) (record *model.GamNetworkModel, err error)
	FindByUser(userID int64) (records []model.GamNetworkModel, err error)
	FindAllAdmin() (records []model.GamNetworkModel, err error)
}

type gamNetworkRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewGamNetworkRepo(db *gorm.DB, cache caching.Cache) *gamNetworkRepo {
	return &gamNetworkRepo{Db: db, Cache: cache}
}

func (t *gamNetworkRepo) FindByID(ID int64) (record *model.GamNetworkModel, err error) {
	err = t.Db.
		First(&record, ID).Error
	return
}

func (t *gamNetworkRepo) FindByUser(userID int64) (records []model.GamNetworkModel, err error) {
	err = t.Db.
		Where("user_id = ?", userID).
		Find(&records).Error
	return
}

func (t *gamNetworkRepo) FindAllAdmin() (records []model.GamNetworkModel, err error) {
	err = t.Db.
		//Debug().
		//Where("network_id = 325081995").
		Where("user_id = 0").
		Group("network_id").
		Find(&records).Error
	return
}
