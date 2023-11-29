package mysql

type TableContentQuestion struct {
	Id             int64  `json:"id"`
	ContentId      int64  `json:"content_id"`
	Title          string `json:"title"`
	Type           int64  `json:"type"`
	BackgroundType int64  `json:"background_type"`
	Background     string `json:"background"`
	PictureType    int64  `json:"picture_type"`
	Picture        string `json:"picture"`
	Answers        string `json:"answers"`
}
