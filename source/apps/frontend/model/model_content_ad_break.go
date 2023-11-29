package model

import "source/core/technology/mysql"

type ContentAdBreak struct{}

type ContentAdBreakRecord struct {
	mysql.TableContentAdBreak
}

func (ContentAdBreakRecord) TableName() string {
	return mysql.Tables.ContentAdBreak
}

func (t *ContentAdBreak) GetByContent(contentId int64) (records []ContentAdBreakRecord) {
	mysql.Client.Where("content_id = ?", contentId).Find(&records)
	return
}

func (t *ContentAdBreak) DeleteAllAdBreakByContent(contentId int64) {
	mysql.Client.Where(ContentAdBreakRecord{mysql.TableContentAdBreak{ContentId: contentId}}).Delete(ContentAdBreakRecord{})
	return
}