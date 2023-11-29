package history

import (
	"encoding/json"
	"source/core/technology/mysql"
	"source/pkg/utility"
)

type BlockedPage struct {
	Detail    DetailBlockedPage
	CreatorId int64
	RecordOld mysql.TableRule
	RecordNew mysql.TableRule
}

func (t *BlockedPage) Page() string {
	return "Rule"
}

type DetailBlockedPage int

const (
	DetailBlockedPageFE DetailBlockedPage = iota + 1
)

func (t DetailBlockedPage) String() string {
	switch t {
	case DetailBlockedPageFE:
		return "blocked_page_fe"
	}
	return ""
}

func (t DetailBlockedPage) App() string {
	switch t {
	case DetailBlockedPageFE:
		return "FE"
	}
	return ""
}

func (t *BlockedPage) Type() TYPEHistory {
	return TYPEHistoryBlockedPage
}

func (t *BlockedPage) Action() mysql.TYPEObjectType {
	if t.RecordOld.ID == 0 && t.RecordNew.ID != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.ID != 0 && t.RecordNew.ID == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *BlockedPage) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailBlockedPageFE:
		return t.getHistoryBlockedPageFE()
	}
	return mysql.TableHistory{}
}

func (t *BlockedPage) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailBlockedPageFE.String():
		return t.compareDataBlockedPageFE(history)
	}
	return []ResponseCompare{}
}

func (t *BlockedPage) getRootRecord() (record mysql.TableRule) {
	switch t.Action() {
	case mysql.TYPEObjectTypeAdd:
		return t.RecordNew
	case mysql.TYPEObjectTypeUpdate:
		return t.RecordNew
	case mysql.TYPEObjectTypeDel:
		return t.RecordOld
	}
	return
}

type blockedPageFE struct {
	Name        *string  `json:"name,omitempty"`
	BlockedPage []string `json:"blocked_page"`
}

func (t *BlockedPage) getHistoryBlockedPageFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := blockedPageFE{}
	newData := blockedPageFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Rule,
		ObjectId:   t.getRootRecord().ID,
		ObjectName: t.getRootRecord().Name,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserID,
	}
	var bNewData, bOldData []byte
	if t.RecordNew.ID != 0 {
		newData.MakeData(t.RecordNew)
		bNewData, _ = json.Marshal(newData)
	}
	if t.RecordOld.ID != 0 {
		oldData.MakeData(t.RecordOld)
		bOldData, _ = json.Marshal(oldData)
	}
	if t.Action() == mysql.TYPEObjectTypeAdd {
		history.Title = "Add Blocked Page"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Blocked Page"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Blocked Page"
		history.OldData = string(bOldData)
	}
	return
}

func (t *BlockedPage) compareDataBlockedPageFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew blockedPageFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Name
	res, err := makeResponseCompare("Name", recordOld.Name, recordNew.Name, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}

	// Page
	for _, page := range recordOld.BlockedPage {
		// Nếu như page không có trong page new thì page đã xóa
		if !utility.InArray(page, recordNew.BlockedPage, false) {
			res, err := makeResponseCompare("Page", &page, nil, mysql.TYPEObjectTypeDel)
			if err == nil {
				responses = append(responses, res)
			}
		}
	}
	for _, page := range recordNew.BlockedPage {
		// Nếu như page không có trong page old thì page được add mới
		if !utility.InArray(page, recordOld.BlockedPage, false) {
			res, err := makeResponseCompare("Page", nil, &page, mysql.TYPEObjectTypeAdd)
			if err == nil {
				responses = append(responses, res)
			}
		}
	}
	return
}

func (rec *blockedPageFE) MakeData(record mysql.TableRule) {
	rec.Name = &record.Name
	for _, blockedPage := range record.BlockedPages {
		if !utility.InArray(blockedPage.Page, rec.BlockedPage, false) {
			rec.BlockedPage = append(rec.BlockedPage, blockedPage.Page)
		}
	}
}
