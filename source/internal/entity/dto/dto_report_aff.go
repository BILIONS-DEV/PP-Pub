package dto

import (
	"source/internal/entity/dto/datatable"
)

type PayloadReportAffIndex struct {
	datatable.Request
	OrderColumn    int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir       string   `query:"order_dir" json:"order_dir" form:"order_dir"`
	StartDate      string   `query:"f_start_day" json:"f_start_day" form:"f_start_day"`
	EndDate        string   `query:"f_end_day" json:"f_end_day" form:"f_end_day"`
	Campaigns      []string `query:"f_campaign" json:"f_campaign" form:"f_campaign"`
	TrafficSources []string `query:"f_traffic_source" json:"f_traffic_source" form:"f_traffic_source"`
	DemandSources  []string `query:"f_demand_source" json:"f_demand_source" form:"f_demand_source"`
	SectionID      []string `query:"f_section_id" json:"f_section_id" form:"f_section_id"`
	GroupBy        []string `query:"f_group_by" json:"f_group_by" form:"f_group_by"`
}

type PayloadReportAffIndexPost struct {
	datatable.Request
	UserID   int64                          `json:"user_id"`
	PostData *PayloadReportAffIndexPostData `query:"postData"`
}

type PayloadReportAffIndexPostData struct {
	QuerySearch    string      `query:"f_q" json:"f_q" form:"f_q"`
	StartDate      string      `query:"f_start_day" json:"f_start_day" form:"f_start_day"`
	EndDate        string      `query:"f_end_day" json:"f_end_day" form:"f_end_day"`
	Campaigns      interface{} `query:"f_campaign[]" json:"f_campaign" form:"f_campaign[]"`
	TrafficSources interface{} `query:"f_traffic_source[]" json:"f_traffic_source" form:"f_traffic_source[]"`
	DemandSources  interface{} `query:"f_demand_source[]" json:"f_demand_source" form:"f_demand_source[]"`
	SectionID      interface{} `query:"f_section_id[]" json:"f_section_id" form:"f_section_id[]"`
	GroupBy        interface{} `query:"f_group_by[]" json:"f_group_by" form:"f_group_by[]"`
}

func (t *PayloadReportAffIndexPost) Validate() (errs []Error) {
	return
}
