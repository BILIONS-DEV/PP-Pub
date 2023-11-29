package country

import (
	"gorm.io/gorm"
	"source/internal/entity/model"
)

type RepoCountry interface {
	FindAll() (records []model.CountryModel, err error)
	GetAllCountryOutbrain() (records []model.CountryOutbrainModel, err error)
}

func NewCountryRepo(DB *gorm.DB) *countryRepo {
	return &countryRepo{Db: DB}
}

type countryRepo struct {
	Db *gorm.DB
}

func (t *countryRepo) FindAll() (records []model.CountryModel, err error) {
	err = t.Db.Find(&records).Error
	return
}

func (t *countryRepo) GetAllCountryOutbrain() (records []model.CountryOutbrainModel, err error) {
	err = t.Db.Find(&records).Error
	return
}
