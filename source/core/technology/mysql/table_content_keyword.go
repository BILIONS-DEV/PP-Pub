package mysql

type TableContentKeyword struct {
	Id        int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	ContentId int64  `json:"content_id"`
	Keyword   string `json:"keyword"`
}
