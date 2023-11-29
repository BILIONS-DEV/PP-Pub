package api

import (
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	"strings"
)

type UsecaseAPI interface {
	ListDomainByBidderStatus(input InputListDomainByBidderStatus) (domains []model.InventoryModel, err error)
	GetDomainBidderByBidderAndDomains(input InputGetDomainBidderByBidderAndDomains) (domainBiddersStatus []model.RlsBidderSystemInventoryModel, err error)
	UpdateAdsTxtForDomains(inputs []InputUpdateAdsTxtForDomains) (err error)
	UpdateAdsTxtForPublisher(inputs InputUpdateAdsTxtForPublisher) (err error)
	ListTagIDByDomainName(domain string) (tagIDs []int64, err error)
	GetUsersByPermission(permissions []int64) (records []model.User, err error)
}

type apiUsecase struct {
	repos *repo.Repositories
}

func NewApiUsecase(repos *repo.Repositories) *apiUsecase {
	return &apiUsecase{repos: repos}
}

type InputListDomainByBidderStatus struct {
	BidderID int64
	Bidder   string
	Status   string
}

func (t *apiUsecase) ListDomainByBidderStatus(input InputListDomainByBidderStatus) (domains []model.InventoryModel, err error) {
	// Validate
	if input.BidderID == 0 && input.Bidder == "" {
		err = errors.New("bidder or bidder_id is require")
	}

	if input.BidderID == 0 {
		bidder, err := t.repos.Bidder.FindByName(input.Bidder, 0)
		if err != nil {
			return domains, err
		}
		if bidder.ID == 0 {
			err = errors.New("bidder not found")
			return domains, err
		}
		input.BidderID = bidder.ID
	}

	// Chuyển input status về type chuẩn với db
	var status []model.TYPERlsBidderSystemInventoryStatus
	switch input.Status {
	case model.RlsBidderSystemInventoryPending.String():
		status = append(status, model.RlsBidderSystemInventoryPending)
		break
	case model.RlsBidderSystemInventorySubmitted.String():
		status = append(status, model.RlsBidderSystemInventorySubmitted)
		break
	case model.RlsBidderSystemInventoryApproved.String():
		status = append(status, model.RlsBidderSystemInventoryApproved)
		status = append(status, model.RlsBidderSystemInventoryApprovedS2S)
		status = append(status, model.RlsBidderSystemInventoryApprovedClient)
		break
	case model.RlsBidderSystemInventoryRejected.String():
		status = append(status, model.RlsBidderSystemInventoryRejected)
		break
	case model.RlsBidderSystemInventoryQueue.String():
		status = append(status, model.RlsBidderSystemInventoryQueue)
		break
	case model.RlsBidderSystemInventoryNotfound.String():
		status = append(status, model.RlsBidderSystemInventoryNotfound)
		break
	}
	return t.repos.Inventory.FindByBidderStatus(input.BidderID, status)
}

type InputGetDomainBidderByBidderAndDomains struct {
	Bidder  string
	Domains []string
}

func (t *apiUsecase) GetDomainBidderByBidderAndDomains(input InputGetDomainBidderByBidderAndDomains) (domainBiddersStatus []model.RlsBidderSystemInventoryModel, err error) {
	var bidderSystem *model.BidderModel
	if bidderSystem, err = t.repos.Bidder.FindByName(input.Bidder, 0); err != nil || bidderSystem.ID == 0 {
		return
	}

	domainBiddersStatus, err = t.repos.Bidder.ListDomainBidderStatusByBidder(bidderSystem.ID, input.Domains)
	return
}

type InputUpdateAdsTxtForDomains struct {
	Domain string
	AdsTxt string
}

func (t *apiUsecase) UpdateAdsTxtForDomains(inputs []InputUpdateAdsTxtForDomains) (err error) {
	for _, value := range inputs {
		var domains []*model.InventoryModel
		if domains, err = t.repos.Inventory.GetInventoriesByName(value.Domain); err != nil {
			return err
		}
		if len(domains) == 0 {
			continue
		}
		for _, domain := range domains {
			adsTxt := strings.TrimSpace(string(domain.AdsTxtCustomByAdmin) + "\n" + value.AdsTxt)
			if strings.Contains(string(domain.AdsTxtCustomByAdmin), value.AdsTxt) {
				continue
			}
			if err = t.repos.Inventory.UpdateAdsTxtDomain(domain.ID, adsTxt); err != nil {
				return
			}
		}
	}
	return
}

type InputUpdateAdsTxtForPublisher struct {
	Email  string
	AdsTxt string
}

func (t *apiUsecase) UpdateAdsTxtForPublisher(input InputUpdateAdsTxtForPublisher) (err error) {
	Publisher := t.repos.User.GetByEmail(input.Email)
	if Publisher.ID == 0 {
		err = errors.New("Publisher does not exist!")
		return
	}

	adsTxt := strings.TrimSpace(string(Publisher.AdsTxtCustomByAdmin))
	if adsTxt == "" {
		adsTxt = input.AdsTxt
	} else {
		adsTxt = adsTxt + "\n" + input.AdsTxt
	}
	arrayAdsTxt := strings.Split(adsTxt, "\n")
	for index, value := range arrayAdsTxt {
		arrayAdsTxt[index] = strings.TrimSpace(value)
	}
	uniqueAdsTxt := []string{}
	tempMap := map[string]bool{}

	for _, num := range arrayAdsTxt {
		if _, ok := tempMap[num]; !ok {
			tempMap[num] = true
			uniqueAdsTxt = append(uniqueAdsTxt, num)
		}
	}
	newAdsTxt := strings.Join(uniqueAdsTxt, "\n")
	if err = t.repos.User.UpdateAdsTxtPublisher(Publisher.ID, newAdsTxt); err != nil {
		return
	}
	if err = t.repos.User.ResetCacheAll(Publisher.ID); err != nil {
		return
	}
	return
}

// get all tags by domain name
func (t *apiUsecase) ListTagIDByDomainName(domain string) (tagIDs []int64, err error) {
	var records []model.InventoryAdTagModel
	if records, err = t.repos.Inventory.ListTagIdByDomainName(domain); err != nil || len(records) == 0 {
		return
	}
	for _, value := range records {
		tagIDs = append(tagIDs, value.ID)
	}
	return
}

func (t *apiUsecase) GetUsersByPermission(permissions []int64) (records []model.User, err error) {
	return t.repos.User.GetUsersByPermission(permissions)
}
