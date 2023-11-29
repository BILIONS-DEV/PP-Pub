package payload

import "source/pkg/datatable"

type GamIndex struct {
	GamFilterPostData
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
}

type GamFilterPayload struct {
	datatable.Request
	PostData *GamFilterPostData `query:"postData"`
}
type GamFilterPostData struct {
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
	Page        int    `query:"page" json:"page" form:"page"`
	Limit       int    `query:"limit" json:"limit" form:"limit"`
	Start       int    `query:"start" json:"start" form:"start"`
	Length      int    `query:"length" json:"length" form:"length"`
}

type GamSelectGam struct {
	GamId     int64 `query:"gam_id" json:"gam_id"`
	NetworkId int64 `query:"network_id" json:"network_id"`
	Select    bool  `query:"select" json:"select"`
}

type GamCheckApiAcess struct {
	GamId       int64   `query:"gam_id" json:"gam_id"`
	ListNetwork []int64 `query:"list_network" json:"list_network"`
}
