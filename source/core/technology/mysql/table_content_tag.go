package mysql

type TableContentTag struct {
	Id        int64  `json:"id"`
	ContentId int64  `json:"content_id"`
	Tag       string `json:"tag"`
}
