package model

import (
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/pagination"
)

type Category struct{}

type CategoryRecord struct {
	mysql.TableCategory
}

func (CategoryRecord) TableName() string {
	return mysql.Tables.Category
}

func (t *Category) GetAll() (records []CategoryRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *Category) GetById(id int64) (records CategoryRecord) {
	mysql.Client.Find(&records, id)
	return
}

func (t *Category) GetCategoryNameById(id int64) string {
	var records CategoryRecord
	mysql.Client.Find(&records, id)
	return records.Name
}

func (t *Category) LoadMoreData(key, value string, userID int64, filterTarget payload.FilterTarget, listSelected []int64) (rows []CategoryRecord, isMoreData, lastPage bool) {
	limit := 40
	page, offset := pagination.Pagination(key, limit)
	Query := mysql.Client.Limit(limit)
	if len(listSelected) > 0 {
		Query.Where("name like ? and id not in ?", "%"+value+"%", listSelected)
	} else {
		Query.Where("name like ?", "%"+value+"%")
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
		categoriesID := []int64{}
		for _, value := range channels {
			categoriesID = append(categoriesID, value.Category)
		}

		if categoriesID != nil {
			Query.Where("id in ?", categoriesID)
		}
	} else {
		channels := []ChannelsRecord{}
		mysql.Client.Model(&ChannelsRecord{}).Where("user_id = ?", userID).Find(&channels)

		if len(channels) == 0 {
			return
		}
		categoriesID := []int64{}
		for _, value := range channels {
			categoriesID = append(categoriesID, value.Category)
		}

		if categoriesID != nil {
			Query.Where("id in ?", categoriesID)
		}
	}

	Query.Offset(offset).Find(&rows)

	total := int64(0)
	Query.Offset(0).Limit(10000000).Count(&total)
	// total := t.CountData(value)
	totalPages := int(total) / limit
	if (int(total) % limit) != 0 {
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

func (t *Category) CountData(value string) (count int64) {
	mysql.Client.Model(&ChannelsRecord{}).Where("name like ?", "%"+value+"%").Count(&count)
	return
}