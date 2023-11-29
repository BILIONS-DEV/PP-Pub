package billing

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
)

type userBilling struct {
	DB    *gorm.DB
	Cache caching.Cache
}

type RepoUserBilling interface {
	GetByUserID(userID int64) (billing model.UserBilling)
}

func NewUserBillingRepo(DB *gorm.DB) *userBilling {
	return &userBilling{DB: DB}
}

func (t *userBilling) GetByUserID(userID int64) (billing model.UserBilling) {
	t.DB.Where("user_id = ?", userID).Find(&billing)
	return
}
