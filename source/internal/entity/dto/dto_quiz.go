package dto

import "source/internal/entity/dto/datatable"

type PayloadQuizSubmit struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Illustration string     `json:"illustration"`
	Category     int64      `json:"category"`
	ContentType  int64      `json:"content_type"`
	Questions    []Question `json:"questions"`
}

type Question struct {
	ID         int64    `json:"id"`
	Question   string   `json:"question"`
	Background string   `json:"background"`
	Answer     []Answer `json:"answer"`
}

type Answer struct {
	ID      int64  `json:"id"`
	Img     string `json:"img"`
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

type QuizIndex struct {
	QuizFilterPostData
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
	Type        []int  `query:"f_type" form:"f_type" json:"f_type"`
}

type QuizAdd struct {
	QuizFilterPostData
	Channels int64 `query:"f_channels" json:"f_channels" form:"f_channels"`
}

type QuizFilterPayload struct {
	datatable.Request
	PostData *QuizFilterPostData `query:"postData"`
}

type QuizFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}
