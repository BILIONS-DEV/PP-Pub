package country

import (
	"source/internal/entity/model"
	"source/internal/repo"
)

type UsecaseCountry interface {
	GetAll() (records []model.CountryModel, err error)
}

func NewCountryUC(repos *repo.Repositories) *CountryUC {
	return &CountryUC{repos: repos}
}

type CountryUC struct {
	repos *repo.Repositories
}

func (t *CountryUC) GetAll() (records []model.CountryModel, err error) {
	records, err = t.repos.Country.FindAll()
	return
}
