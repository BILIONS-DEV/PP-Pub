package ad_type

import (
	"source/internal/entity/model"
	"source/internal/repo"
)

type UsecaseAdType interface {
	GetById(id int64) (record *model.AdTypeModel, err error)
	GetByUser(userId int64) (records []model.AdTypeModel, err error)
}

func NewAdTypeUC(repos *repo.Repositories) *adTypeUC {
	return &adTypeUC{repos: repos}
}

type adTypeUC struct {
	repos *repo.Repositories
}

func (t *adTypeUC) GetById(id int64) (record *model.AdTypeModel, err error) {
	record, err =
		t.repos.AdType.FindByID(id)
	return
}

func (t *adTypeUC) GetByUser(userId int64) (record []model.AdTypeModel, err error) {
	record, err =
		t.repos.AdType.FindByUser(userId)
	return
}
