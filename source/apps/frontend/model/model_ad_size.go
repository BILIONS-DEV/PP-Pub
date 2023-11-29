package model

import (
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/pagination"
)

type AdSize struct{}

type AdSizeRecord struct {
	mysql.TableAdSize
}

func (AdSizeRecord) TableName() string {
	return mysql.Tables.AdSize
}

func (t *AdSize) GetById(id int64) (record AdSizeRecord) {
	if id < 1 {
		return
	}
	mysql.Client.First(&record, id)
	return
}

func (t *AdSize) GetAll() (records []AdSizeRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *AdSize) GetAllSizeForLineGoogle() (records []AdSizeRecord) {
	mysql.Client.Where("type = 1").Find(&records)
	return
}

func (t *AdSize) GetAllSizeForNative() (records []AdSizeRecord) {
	mysql.Client.Where("for_native = 1").Find(&records)
	return
}

func (t *AdSize) GetSizeAdditional(record AdSizeRecord, typ string) (records []AdSizeRecord) {
	if typ == "sticky_mobile" {
		listSizeStickyMobile := new(AdSize).GetSizeStickyMobile()
		var listId []int64
		for _, adSize := range listSizeStickyMobile {
			if record.Id != adSize.Id {
				listId = append(listId, adSize.Id)
			}
		}
		mysql.Client.Where("id in ?", listId).Find(&records)
		return
	}
	mysql.Client.Where("width <= ? and height <= ? and id != ? and width != 1", record.Width, record.Height, record.Id).Find(&records)
	return
}

func (t *AdSize) CountData(value string) (count int64) {
	mysql.Client.Model(&AdSizeRecord{}).Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *AdSize) CountDataPageEdit(listSelected []int64) (count int64) {
	if len(listSelected) > 0 {
		mysql.Client.Model(&AdSizeRecord{}).Where("id not in ?", listSelected).Count(&count)
	} else {
		mysql.Client.Model(&AdSizeRecord{}).Count(&count)
	}
	return
}

func (t *AdSize) LoadMoreData(key, value string, listSelected []int64) (rows []AdSizeRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination(key, limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("name like ? and id not in ?", "%"+value+"%", listSelected).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.Where("name like ?", "%"+value+"%").Limit(limit).Offset(offset).Find(&rows)
	}
	total := t.CountData(value)
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

func (t *AdSize) LoadMoreDataPageEdit(listSelected []int64) (rows []AdSizeRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination("1", limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("id not in ?", listSelected).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.Limit(limit).Offset(offset).Find(&rows)
	}
	if len(rows) > 10 {
		rows = rows[0:9]
	}
	total := t.CountDataPageEdit(listSelected)
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

func appendSize(a []int64, b []int64) []int64 {
	check := make(map[int64]int64)
	d := append(a, b...)
	res := make([]int64, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}

func (t *AdSize) GetAllApi() (data []payload.ListTargetCheck) {
	mysql.Client.Model(&AdSizeRecord{}).Select("id", "name").Find(&data)
	return
}

func (t *AdSize) GetSizeStickyDesktop() (listSize []AdSizeRecord) {
	listId := []int{10, 2, 1, 3, 4, 7, 9}
	mysql.Client.Model(&AdSizeRecord{}).Where("id in ?", listId).Find(&listSize)
	return
}

func (t *AdSize) GetSizeDefaultStickyDesktop(position int) (listSize []AdSizeRecord) {
	listId := []int{}

	switch position {
	case 1, 6:
		listId = []int{2, 10}
	case 2, 3:
		listId = []int{1, 3, 4}
	}
	mysql.Client.Model(&AdSizeRecord{}).Where("id in ?", listId).Find(&listSize)
	return
}

func (t *AdSize) GetSizeDefaultStickyMobile() (listSize []AdSizeRecord) {
	listId := []int{7, 9, 19, 20}
	mysql.Client.Model(&AdSizeRecord{}).Where("id in ?", listId).Find(&listSize)
	return
}

func (t *AdSize) GetSizeStickyMobile() (listSize []AdSizeRecord) {
	listId := []int{7, 9, 17, 18, 18, 19, 20}
	mysql.Client.Model(&AdSizeRecord{}).Where("id in ?", listId).Find(&listSize)
	return
}
