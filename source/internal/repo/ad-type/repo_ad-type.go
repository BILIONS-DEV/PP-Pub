package ad_type

import (
	"gorm.io/gorm"
	"source/internal/entity/model"
)

type RepoAdType interface {
	FindByID(id int64) (record *model.AdTypeModel, err error)
	FindByUser(userID int64) (records []model.AdTypeModel, err error)
}

func NewAdTypeRepo(DB *gorm.DB) *adTypeRepo {
	return &adTypeRepo{Db: DB}
}

type adTypeRepo struct {
	Db *gorm.DB
}

func (t *adTypeRepo) FindByID(id int64) (record *model.AdTypeModel, err error) {
	err = t.Db.Find(&record, id).Error
	return
}

func (t *adTypeRepo) FindByUser(userID int64) (records []model.AdTypeModel, err error) {
	err = t.Db.Where("user_id = ?", userID).Find(&records).Error
	return
}
