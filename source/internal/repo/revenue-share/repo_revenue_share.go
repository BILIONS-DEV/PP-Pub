package revenue_share

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
)

type revenueShare struct {
	DB    *gorm.DB
	Cache caching.Cache
}

type RepoRevenueShare interface {
	GetByUserID(userID int64) (revenueShare model.RevenueShare)
}

func NewRevenueShareRepo(DB *gorm.DB) *revenueShare {
	return &revenueShare{DB: DB}
}

func (t *revenueShare) GetByUserID(userID int64) (revenueShare model.RevenueShare) {
	t.DB.Where("user_id = ?", userID).Find(&revenueShare)
	return
}
