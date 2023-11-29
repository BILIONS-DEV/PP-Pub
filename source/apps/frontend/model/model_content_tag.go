package model

import "source/core/technology/mysql"

type ContentTag struct{}

type ContentTagRecord struct {
	mysql.TableContentTag
}

func (ContentTagRecord) TableName() string {
	return mysql.Tables.ContentTag
}

func (t *ContentTag) GetByContent(contentId int64) (records []ContentTagRecord) {
	mysql.Client.Where("content_id = ?", contentId).Find(&records)
	return
}

func (t *ContentTag) DeleteTagByContent(contentId int64) {
	mysql.Client.Where(ContentTagRecord{mysql.TableContentTag{ContentId: contentId}}).Delete(ContentTagRecord{})
	return
}