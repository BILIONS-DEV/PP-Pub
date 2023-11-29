package bidder

import (
	"source/internal/entity/model"
	"source/internal/repo"
)

type UsecaseBidder interface {
	GetById(id int64) (record *model.BidderModel, err error)
	GetAllBidderGoogle() (records []*model.BidderModel, err error)
}

func NewBidderUC(repos *repo.Repositories) *bidderUC {
	return &bidderUC{repos: repos}
}

type bidderUC struct {
	repos *repo.Repositories
}

func (t *bidderUC) GetById(id int64) (record *model.BidderModel, err error) {
	record, err = t.repos.Bidder.FindByID(id)
	return
}

func (t *bidderUC) GetAllBidderGoogle() (records []*model.BidderModel, err error) {
	return t.repos.Bidder.FindAllByQuery(model.BidderModel{
		BidderTemplateId: 1,
	}, "user_id = 0", "pub_id != ''")
}
