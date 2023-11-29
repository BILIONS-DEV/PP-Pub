package ads

import (
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/repo"
	adsRepo "source/internal/repo/ads"
)

type UsecaseAds interface {
	Filter(input *InputFilterUC) (records []model.AdsModel, total int64, err error)
	Submit(domain string) (err error)
	ChangeActionAd(input *ChangeActionAdUC) (err error)
}

func NewAdsUsecase(repos *repo.Repositories) *adsUC {
	return &adsUC{repos: repos}
}

type adsUC struct {
	repos *repo.Repositories
}

type InputFilterUC struct {
	datatable.Request
	User        model.User
	QuerySearch interface{}
	Status      interface{}
	Inventory   string
	OrderColumn int
	OrderDir    string
}

func (t *adsUC) Filter(input *InputFilterUC) (records []model.AdsModel, total int64, err error) {
	// => khởi tạo input của Repo
	inputRepo := adsRepo.InputFilter{
		User:      input.User,
		Status:    input.Status,
		Inventory: input.Inventory,
		Search:    input.QuerySearch,
		Page:      (input.Start / input.Length) + 1,
		Limit:     input.Length,
		Order:     input.OrderString(),
	}
	return t.repos.Ads.Filter(&inputRepo)
}

func (t *adsUC) Submit(domain string) (err error) {
	return
}

type ChangeActionAdUC struct {
	UserLogin model.User
	AdID      int64
	Action    string
	Inventory string
	Placement string
}

func (t *adsUC) ChangeActionAd(input *ChangeActionAdUC) (err error) {
	return t.repos.Ads.ChangeActionAd(adsRepo.ChangeActionAdRP{
		UserLogin: input.UserLogin,
		Email:     input.UserLogin.Email,
		AdID:      input.AdID,
		Action:    input.Action,
		Inventory: input.Inventory,
		Placement: input.Placement,
	})
}
