package model

import (
	"source/core/technology/mysql"
	"source/pkg/ajax"
)

type PageCollapse struct{}

type PageCollapseRecord struct {
	mysql.TablePageCollapse
}

func (PageCollapseRecord) TableName() string {
	return mysql.Tables.PageCollapse
}

func (t *PageCollapse) HandleCollapse(record PageCollapseRecord) (errs []ajax.Error) {
	err, data := t.CheckExist(record)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	if data.Id != 0 {
		record.Id = data.Id
		err = mysql.Client.Updates(&record).Where("id = ?", data.Id).Error
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
	} else {
		err = mysql.Client.Create(&record).Error
		if err != nil {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
	}
	return
}

func (t *PageCollapse) CheckExist(record PageCollapseRecord) (err error, data *PageCollapseRecord) {
	switch record.PageType {
	case "add":
		err = mysql.Client.Where("page_collapse = ? and box_collapse = ? and user_id = ? and page_type = ?", record.PageCollapse, record.BoxCollapse, record.UserId, record.PageType).Find(&data).Error
		return
	case "edit":
		err = mysql.Client.Where("page_collapse = ? and box_collapse = ? and user_id = ? and page_type = ? and page_id = ?", record.PageCollapse, record.BoxCollapse, record.UserId, record.PageType, record.PageId).Find(&data).Error
		return
	}
	return
}
