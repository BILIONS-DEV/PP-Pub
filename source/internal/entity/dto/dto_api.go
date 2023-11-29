package dto

import "source/internal/errors"

type PayloadAPIListDomainByBidderStatus struct {
	BidderID int64  `json:"bidder_id" form:"bidder_id"`
	Bidder   string `json:"bidder" form:"bidder"`
	Status   string `json:"status" form:"status"`
}

func (t *PayloadAPIListDomainByBidderStatus) Validate() (err error) {
	if t.Bidder == "" && t.BidderID == 0 {
		err = errors.New("Require bidder or bidder_id")
	}
	return
}

type ResponseAPIListDomainByBidderStatus struct {
	DomainID int64  `json:"domain_id"`
	Domain   string `json:"domain"`
}

type PayloadGetDomainBidderByBidderAndDomains struct {
	Bidder  string `json:"bidder" form:"bidder"`
	Domains string `json:"domain" form:"domain"`
}

func (t *PayloadGetDomainBidderByBidderAndDomains) Validate() (err error) {
	if t.Bidder == "" {
		err = errors.New("Bidder is required ")
	}
	return
}

type ResponseGetDomainBidderByBidderAndDomains struct {
	Domain string `json:"domain"`
	Bidder string `json:"bidder"`
	Status string `json:"status"`
}

type PayloadUpdateAdsTxtForDomains struct {
	AdsTxt string `json:"ads_txt" form:"ads_txt"`
	Domain string `json:"domain" form:"domain"`
}

type UpdateAdsTxtForDomains struct {
	AdsTxt string `json:"ads_txt" form:"ads_txt"`
	Domain string `json:"domain" form:"domain"`
}

type ResponseUpdateAdsTxtForDomains struct {
	Domain string `json:"domain"`
	Bidder string `json:"bidder"`
	Status string `json:"status"`
}

type PayloadUpdateAdsTxtForPublisher struct {
	AdsTxt string `json:"ads_txt" form:"ads_txt"`
	Email  string `json:"email" form:"email"`
}

type PayloadGetInvoiceByPublisher struct {
	Publishers []string `json:"publishers" form:"publishers"`
}

type PayloadGetReportAgency struct {
	Email string `json:"email" form:"email"`
	Month string `json:"month" form:"month"`
}
