package dto

import (
	"source/internal/entity/dto/datatable"
)

type PayloadReportBodisFilter struct {
	datatable.Request
	StartDate string   `query:"f_start_day" json:"f_start_day" form:"f_start_day"`
	EndDate   string   `query:"f_end_day" json:"f_end_day" form:"f_end_day"`
	SubID     string   `query:"f_subid" json:"f_subid" form:"f_subid"`
	OrderBy   string   `query:"f_orderby" json:"f_orderby" form:"f_orderby"`
	GroupBy   []string `query:"f_group_by" json:"f_group_by" form:"f_group_by"`
}

func (t *PayloadReportBodisFilter) Validate() (err error) {
	return
}

func (t *PayloadReportBodisFilter) ToCondition() map[string]interface{} {
	var condition = make(map[string]interface{})
	if t.SubID != "" {
		condition["subid"] = t.SubID
	}

	return condition
}

// type PayloadGetHistory struct {
// 	Id     string  `json:"id" form:"id"`
// 	Object string `json:"object" form:"object"`
// }

type ReportBodisParkingSearch struct {
	CurrentPage  int                            `json:"current_page"`
	Data         []ReportBodisParkingSearchData `json:"data"`
	FirstPageURL string                         `json:"first_page_url"`
	From         int                            `json:"from"`
	NextPageURL  interface{}                    `json:"next_page_url"`
	Path         string                         `json:"path"`
	PerPage      int                            `json:"per_page"`
	PrevPageURL  interface{}                    `json:"prev_page_url"`
	To           int                            `json:"to"`
}

type ReportBodisParkingSearchData struct {
	VisitID          string  `json:"visit_id"`
	DomainName       string  `json:"domain_name"`
	IPAddress        string  `json:"ip_address"`
	Type             string  `json:"type"`
	ServerDatetime   string  `json:"server_datetime"`
	CountryID        int64   `json:"country_id"`
	PageQuery        string  `json:"page_query"`
	Subids           string  `json:"subids"`
	Visits           int     `json:"visits"`
	Zeroclicks       int     `json:"zeroclicks"`
	Ctrs             int     `json:"ctrs"`
	Clicks           int64   `json:"clicks"`
	EstimatedRevenue float64 `json:"estimated_revenue"`
}
