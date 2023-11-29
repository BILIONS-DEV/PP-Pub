package dto

import (
	"source/internal/entity/dto/datatable"
)

type PayloadReportAffSelectionIndex struct {
	datatable.Request
	StartDate      string   `query:"f_start_day" json:"f_start_day" form:"f_start_day"`
	EndDate        string   `query:"f_end_day" json:"f_end_day" form:"f_end_day"`
	Campaigns      []string `query:"f_campaign" json:"f_campaign" form:"f_campaign"`
	TrafficSources []string `query:"f_traffic_source" json:"f_traffic_source" form:"f_traffic_source"`
	PublisherID    []string `query:"f_publisher_id" json:"f_publisher_id" form:"f_publisher_id"`
	SectionID      []string `query:"f_section_id" json:"f_section_id" form:"f_section_id"`
	GroupBy        []string `query:"f_group_by" json:"f_group_by" form:"f_group_by"`
}

type PayloadReportAffSelectionIndexPost struct {
	datatable.Request
	UserID   int64                                   `json:"user_id"`
	PostData *PayloadReportAffSelectionIndexPostData `query:"postData"`
}

type PayloadReportAffSelectionIndexPostData struct {
	QuerySearch    string      `query:"f_q" json:"f_q" form:"f_q"`
	StartDate      string      `query:"f_start_day" json:"f_start_day" form:"f_start_day"`
	EndDate        string      `query:"f_end_day" json:"f_end_day" form:"f_end_day"`
	Campaigns      interface{} `query:"f_campaign[]" json:"f_campaign" form:"f_campaign[]"`
	TrafficSources interface{} `query:"f_traffic_source[]" json:"f_traffic_source" form:"f_traffic_source[]"`
	PublisherID    interface{} `query:"f_publisher_id[]" json:"f_publisher_id" form:"f_publisher_id[]"`
	SectionID      interface{} `query:"f_section_id[]" json:"f_section_id" form:"f_section_id[]"`
	GroupBy        interface{} `query:"f_group_by[]" json:"f_group_by" form:"f_group_by[]"`
}

func (t *PayloadReportAffSelectionIndexPost) Validate() (errs []Error) {
	return
}
