package model

import (
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/pagination"
)

type ContentKeyword struct{}

type ContentKeywordRecord struct {
	mysql.TableContentKeyword
}

func (ContentKeywordRecord) TableName() string {
	return mysql.Tables.ContentKeyword
}

func (t *ContentKeyword) GetByContent(contentId int64) (records []ContentKeywordRecord) {
	mysql.Client.Where("content_id = ?", contentId).Find(&records)
	return
}

func (t *ContentKeyword) GetByUserId(userId int64) (records []ContentKeywordRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(&records)
	return
}

func (t *ContentKeyword) DeleteKeywordByContent(contentId int64) {
	mysql.Client.Where(ContentKeywordRecord{mysql.TableContentKeyword{ContentId: contentId}}).Delete(ContentKeywordRecord{})
	return
}

func (t *ContentKeyword) LoadMoreData(key, value string, userID int64, filterTarget payload.FilterTarget, listSelected []int64) (rows []ContentKeywordRecord, isMoreData, lastPage bool) {
	limit := 30
	page, offset := pagination.Pagination(key, limit)
	Query := mysql.Client.Limit(limit).Order("keyword ASC")
	if len(listSelected) > 0 {
		Query.Where("keyword like ? and user_id = ? and id not in ?", "%"+value+"%", userID, listSelected)
	} else {
		Query.Where("keyword like ? and user_id = ?", "%"+value+"%", userID)
	}

	if len(filterTarget.Language) > 0 || len(filterTarget.Channels) > 0 || len(filterTarget.ExcludeLanguage) > 0 || len(filterTarget.ExcludeChannels) > 0 {
		channels := []ChannelsRecord{}
		queryChannels := mysql.Client.Model(&ChannelsRecord{})
		if len(filterTarget.Language) > 0 {
			queryChannels.Where("user_id = ? and language in ?", userID, filterTarget.Language)
		}
		if len(filterTarget.ExcludeLanguage) > 0 {
			queryChannels.Where("user_id = ? and language not in ?", userID, filterTarget.ExcludeLanguage)
		}
		if len(filterTarget.Channels) > 0 {
			queryChannels.Where("id in ?", filterTarget.Channels)
		}
		if len(filterTarget.ExcludeChannels) > 0 {
			queryChannels.Where("id not in ?", filterTarget.ExcludeChannels)
		}
		queryChannels.Find(&channels)
		if len(channels) == 0 {
			return
		}
		channelsID := []int64{}
		for _, value := range channels {
			channelsID = append(channelsID, value.Id)
		}

		contents := []ContentRecord{}
		mysql.Client.Model(&ContentRecord{}).Where("channels in ?",  channelsID).Find(&contents)
		if len(contents) == 0 {
			return
		}
		contentID := []int64{}
		for _, value := range contents{
			contentID = append(contentID, value.Id)
		}
		Query.Where("content_id in ?", contentID)
	}

	Query.Offset(offset).Group("keyword").Find(&rows)

	total := int64(0)
	Query.Offset(0).Limit(10000000).Count(&total)
	// total := t.CountData(value, userID)
	totalPages := int(total) / 30
	if (int(total) % 30) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	if page >= totalPages || len(rows) == 0 {
		isMoreData = false
		lastPage = true
	}
	return
}

func (t *ContentKeyword) CountData(value string, userID int64) (count int64) {
	mysql.Client.Model(&ContentKeywordRecord{}).Where("keyword like ? and user_id = ? ", "%"+value+"%", userID).Count(&count)
	return
}
