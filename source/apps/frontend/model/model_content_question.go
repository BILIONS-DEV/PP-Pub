package model

import (
	"encoding/json"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
)

type ContentQuestion struct{}

type ContentQuestionRecord struct {
	mysql.TableContentQuestion
}

func (ContentQuestionRecord) TableName() string {
	return mysql.Tables.ContentQuestion
}

func (t *ContentQuestion) GetByContent(contentId int64) (records []ContentQuestionRecord) {
	mysql.Client.Where("content_id = ?", contentId).Find(&records)
	return
}

func (t *ContentQuestion) DeleteAllQuizContent(contentId int64) {
	mysql.Client.Where(ContentQuestionRecord{mysql.TableContentQuestion{ContentId: contentId}}).Delete(ContentQuestionRecord{})
	return
}

func (t *ContentQuestion) HandleQuestion(records []ContentQuestionRecord) (data []payload.Questions) {
	if records != nil {
		for _, record := range records {
			question := payload.Questions{}
			Answers := []string{}
			json.Unmarshal([]byte(record.Answers), &Answers)
			question.Id = record.Id
			question.Title = record.Title
			question.BackgroundType = record.BackgroundType
			question.Background = record.Background
			question.Type = record.Type
			question.PictureType = record.PictureType
			question.Picture = record.Picture
			question.Answer = Answers
			data = append(data, question)
		}
	}
	return
}
