package gam_network

import (
	"source/internal/entity/model"
	"source/internal/repo"
)

type UsecaseGamNetwork interface {
	GetById(id int64) (record *model.GamNetworkModel, err error)
	GetByUser(userId int64) (records []model.GamNetworkModel, err error)
}

func NewGamNetworkUC(repos *repo.Repositories) *gamNetworkUC {
	return &gamNetworkUC{repos: repos}
}

type gamNetworkUC struct {
	repos *repo.Repositories
}

func (t *gamNetworkUC) GetById(id int64) (record *model.GamNetworkModel, err error) {
	record, err = t.repos.GamNetwork.FindByID(id)
	return
}

func (t *gamNetworkUC) GetByUser(userId int64) (records []model.GamNetworkModel, err error) {
	records, err = t.repos.GamNetwork.FindByUser(userId)
	return
}
