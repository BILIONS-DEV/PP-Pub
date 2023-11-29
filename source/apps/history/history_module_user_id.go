package history

import (
	"encoding/json"
	"errors"
	"source/core/technology/mysql"
)

type ModuleUserId struct {
	Detail    DetailModuleUserId
	CreatorId int64
	RecordOld mysql.TableModuleUserId
	RecordNew mysql.TableModuleUserId
}

func (t *ModuleUserId) Page() string {
	return "Module UserID"
}

type DetailModuleUserId int

const (
	DetailModuleUserIdBE DetailModuleUserId = iota + 1
)

func (t DetailModuleUserId) String() string {
	switch t {
	case DetailModuleUserIdBE:
		return "module_user_id_be"
	}
	return ""
}

func (t DetailModuleUserId) App() string {
	switch t {
	case DetailModuleUserIdBE:
		return "BE"
	}
	return ""
}

func (t *ModuleUserId) Type() TYPEHistory {
	return TYPEHistoryModuleUserId
}

func (t *ModuleUserId) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *ModuleUserId) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailModuleUserIdBE:
		return t.getHistoryModuleUserIdBE()
	}
	return mysql.TableHistory{}
}

func (t *ModuleUserId) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailModuleUserIdBE.String():
		return t.compareDataModuleUserIdBE(history)
	}
	return []ResponseCompare{}
}

func (t *ModuleUserId) getRootRecord() (record mysql.TableModuleUserId) {
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

type moduleUserIdBE struct {
	Name    *string              `json:"name,omitempty"`
	Params  *[]moduleUserIdParam `json:"params,omitempty"`
	Storage *[]moduleUserIdParam `json:"storage,omitempty"`
}

type moduleUserIdParam struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template"`
}

func (t *ModuleUserId) getHistoryModuleUserIdBE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := moduleUserIdBE{}
	newData := moduleUserIdBE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.ModuleUserId,
		ObjectId:   t.getRootRecord().Id,
		ObjectName: t.getRootRecord().Name,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserId,
	}
	var bNewData, bOldData []byte
	if t.RecordNew.Id != 0 {
		newData.MakeData(t.RecordNew)
		bNewData, _ = json.Marshal(newData)
	}
	if t.RecordOld.Id != 0 {
		oldData.MakeData(t.RecordOld)
		bOldData, _ = json.Marshal(oldData)
	}
	if t.Action() == mysql.TYPEObjectTypeAdd {
		history.Title = "Add Module User Id"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Module User Id"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Module User Id"
		history.OldData = string(bOldData)
	}
	return
}

func (t *ModuleUserId) compareDataModuleUserIdBE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew moduleUserIdBE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Name
	res, err := makeResponseCompare("Name", recordOld.Name, recordNew.Name, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Params
	responseParam, err := t.makeResponseCompareParam(recordOld.Params, recordNew.Params)
	if err == nil {
		responses = append(responses, responseParam...)
	}
	// Storage
	responseStorage, err := t.makeResponseCompareParam(recordOld.Storage, recordNew.Storage)
	if err == nil {
		responses = append(responses, responseStorage...)
	}

	return
}

func (rec *moduleUserIdBE) MakeData(record mysql.TableModuleUserId) {
	rec.Name = &record.Name

	var params []moduleUserIdParam
	_ = json.Unmarshal([]byte(record.Params), &params)
	rec.Params = &params

	var storage []moduleUserIdParam
	_ = json.Unmarshal([]byte(record.Storage), &storage)
	rec.Storage = &storage
}

func (t *ModuleUserId) makeResponseCompareParam(oldData, newData *[]moduleUserIdParam) (responses []ResponseCompare, err error) {
	if oldData == nil && newData == nil {
		err = errors.New("no response")
		return
	}
	mapOldData := make(map[string]moduleUserIdParam)
	mapNewData := make(map[string]moduleUserIdParam)

	if oldData != nil {
		for _, param := range *oldData {
			mapOldData[param.Name] = param
		}
	}
	if newData != nil {
		for _, param := range *newData {
			mapNewData[param.Name] = param
		}
	}

	for name, paramOld := range mapOldData {
		if paramNew, ok := mapNewData[name]; ok {
			// Type
			res, err := makeResponseCompare("Params > "+name+" > Type", &paramOld.Type, &paramNew.Type, mysql.TYPEObjectTypeUpdate)
			if err == nil {
				responses = append(responses, res)
			}
			// Template
			res, err = makeResponseCompare("Params > "+name+" > Template", &paramOld.Template, &paramNew.Template, mysql.TYPEObjectTypeUpdate)
			if err == nil {
				responses = append(responses, res)
			}

			// Xóa các param update
			delete(mapNewData, name)
		} else {
			// Nếu chỉ tồn tại oldData tức param này đã bị xóa toàn bộ các newData là nil
			// Type
			res, err := makeResponseCompare("Params > "+name+" > Type", &paramOld.Type, nil, mysql.TYPEObjectTypeDel)
			if err == nil {
				responses = append(responses, res)
			}
			// Template
			res, err = makeResponseCompare("Params > "+name+" > Template", &paramOld.Template, nil, mysql.TYPEObjectTypeDel)
			if err == nil {
				responses = append(responses, res)
			}
		}
	}
	// Các dataNew còn lại đều là các param add mới
	for name, paramNew := range mapNewData {
		// Type
		res, err := makeResponseCompare("Params > "+name+" > Type", nil, &paramNew.Type, mysql.TYPEObjectTypeAdd)
		if err == nil {
			responses = append(responses, res)
		}
		// Template
		res, err = makeResponseCompare("Params > "+name+" > Template", nil, &paramNew.Template, mysql.TYPEObjectTypeAdd)
		if err == nil {
			responses = append(responses, res)
		}
	}
	return
}
