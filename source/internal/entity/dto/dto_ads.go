package dto

import (
	"errors"
	"source/internal/entity/dto/datatable"
	"source/pkg/utility"
)

type PayloadAdsIndex struct {
	datatable.Request
	OrderColumn int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string   `query:"order_dir" json:"order_dir" form:"order_dir"`
	QuerySearch string   `query:"q" json:"q" form:"q"`
	Status      []string `query:"status" form:"status" json:"status"`
	Inventory   string   `query:"inventory" form:"inventory" json:"inventory"`
}

type PayloadAdsIndexPost struct {
	datatable.Request
	OrderColumn int                      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string                   `query:"order_dir" json:"order_dir" form:"order_dir"`
	PostData    *PayloadAdsIndexPostData `query:"postData"`
}

type PayloadAdsIndexPostData struct {
	QuerySearch interface{} `query:"q" json:"q" form:"q"`
	Status      interface{} `query:"status[]" json:"status" form:"status[]"`
	Inventory   string      `query:"inventory" json:"inventory" form:"inventory"`
}

type ResponseAdsDatatable struct {
	AdPreview   string `json:"ad_preview"`
	AdInfo      string `json:"ad_info"`
	Headline    string `json:"headline"`
	SiteName    string `json:"site_name"`
	ClickUrl    string `json:"click_url"`
	Target      string `json:"target"`
	Impressions string `json:"impressions"`
	Action      string `json:"action"`
}

func (t *PayloadAdsIndexPost) Validate() (errs []Error) {
	// if t.PostData == 0 {
	// 	errs = append(errs, Error{Message: "permission denied"})
	// }
	return
}

type PayloadChangeActionAd struct {
	Action    string `json:"action,omitempty"`
	ID        int64  `json:"id,omitempty"`
	Inventory string `json:"inventory,omitempty"`
	Placement string `json:"placement,omitempty"`
}

func (t *PayloadChangeActionAd) Validate() (errs []error) {
	if t.ID == 0 {
		errs = append(errs, errors.New("Ad id is required"))
	}
	if utility.ValidateString(t.Action) == "" {
		errs = append(errs, errors.New("Action is required"))
	}
	return
}
