package history

import (
	"encoding/json"
	"source/core/technology/mysql"
)

type Identity struct {
	Detail    DetailIdentity
	CreatorId int64
	RecordOld mysql.TableIdentity
	RecordNew mysql.TableIdentity
}

func (t *Identity) Page() string {
	return "Identity"
}

type DetailIdentity int

const (
	DetailIdentityFE DetailIdentity = iota + 1
)

func (t DetailIdentity) String() string {
	switch t {
	case DetailIdentityFE:
		return "identity_fe"
	}
	return ""
}

func (t DetailIdentity) App() string {
	switch t {
	case DetailIdentityFE:
		return "FE"
	}
	return ""
}

func (t *Identity) Type() TYPEHistory {
	return TYPEHistoryIdentity
}

func (t *Identity) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *Identity) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailIdentityFE:
		return t.getHistoryIdentityFE()
	}
	return mysql.TableHistory{}
}

func (t *Identity) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailIdentityFE.String():
		return t.compareDataIdentityFE(history)
	}
	return []ResponseCompare{}
}

func (t *Identity) getRootRecord() (record mysql.TableIdentity) {
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

type identityFE struct {
	Name          *string        `json:"name,omitempty"`
	Description   *string        `json:"description,omitempty"`
	SyncDelay     *int           `json:"sync_delay,omitempty"`
	AuctionDelay  *int           `json:"auction_delay,omitempty"`
	Priority      *int           `json:"priority,omitempty"`
	Status        *string        `json:"status,omitempty"`
	UserIdModules []userIdModule `json:"user_id_modules"`
	Target        target         `json:"target"`
}

type userIdModule struct {
	Name    string                 `json:"name"`
	Params  *[]paramModuleUserId   `json:"params,omitempty"`
	Storage *[]storageModuleUserId `json:"storage,omitempty"`
}

type paramModuleUserId struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template"`
}

type storageModuleUserId struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template"`
}

func (t *Identity) getHistoryIdentityFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := identityFE{}
	newData := identityFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Identity,
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
		history.Title = "Add Identity"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Identity"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Identity"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *identityFE) MakeData(record mysql.TableIdentity) {
	rec.Name = &record.Name
	rec.Description = &record.Description
	rec.SyncDelay = &record.SyncDelay
	rec.AuctionDelay = &record.AuctionDelay
	rec.Priority = &record.Priority
	status := record.Status.String()
	rec.Status = &status

	for _, identityUserIdModule := range record.UserIdModules {
		var params []paramModuleUserId
		var storage []storageModuleUserId
		if identityUserIdModule.Params != "" {
			err := json.Unmarshal([]byte(identityUserIdModule.Params), &params)
			if err != nil {
				return
			}
		}
		if identityUserIdModule.Storage != "" {
			err := json.Unmarshal([]byte(identityUserIdModule.Storage), &storage)
			if err != nil {
				return
			}
		}
		var userIdModule userIdModule
		userIdModule.Name = identityUserIdModule.Name
		if len(params) > 0 {
			userIdModule.Params = &params
		}
		if len(storage) > 0 {
			userIdModule.Storage = &storage
		}
		rec.UserIdModules = append(rec.UserIdModules, userIdModule)
	}

	// Xử lý target
	rec.Target = makeTarget(record.Targets)
	// Loại các target không dùng trong type google
	rec.Target.Format = nil
	rec.Target.Size = nil
	rec.Target.AdTag = nil
	rec.Target.Geography = nil
	rec.Target.Device = nil
}

func (t *Identity) compareDataIdentityFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew identityFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Name
	res, err := makeResponseCompare("Name", recordOld.Name, recordNew.Name, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Description
	res, err = makeResponseCompare("Description", recordOld.Description, recordNew.Description, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Sync Delay
	res, err = makeResponseCompare("Sync Delay", pointerIntToString(recordOld.SyncDelay), pointerIntToString(recordNew.SyncDelay), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Auction Delay
	res, err = makeResponseCompare("Auction Delay", pointerIntToString(recordOld.AuctionDelay), pointerIntToString(recordNew.AuctionDelay), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Priority
	res, err = makeResponseCompare("Priority", pointerIntToString(recordOld.Priority), pointerIntToString(recordNew.Priority), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Status
	res, err = makeResponseCompare("Status", recordOld.Status, recordNew.Status, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// UserId Modules
	resModules, err := t.makeResponseUserIdModules(recordOld.UserIdModules, recordNew.UserIdModules)
	if err == nil {
		responses = append(responses, resModules...)
	}
	// Target Domain
	res, err = makeResponseCompare("Target > Domain", pointerArrayStringToString(recordOld.Target.Domain), pointerArrayStringToString(recordNew.Target.Domain), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Format
	res, err = makeResponseCompare("Target > Format", pointerArrayStringToString(recordOld.Target.Format), pointerArrayStringToString(recordNew.Target.Format), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Size
	res, err = makeResponseCompare("Target > Size", pointerArrayStringToString(recordOld.Target.Size), pointerArrayStringToString(recordNew.Target.Size), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Ad Tag
	res, err = makeResponseCompare("Target > Ad Tag", pointerArrayStringToString(recordOld.Target.AdTag), pointerArrayStringToString(recordNew.Target.AdTag), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Geography
	res, err = makeResponseCompare("Target > Geography", pointerArrayStringToString(recordOld.Target.Geography), pointerArrayStringToString(recordNew.Target.Geography), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Device
	res, err = makeResponseCompare("Target > Device", pointerArrayStringToString(recordOld.Target.Device), pointerArrayStringToString(recordNew.Target.Device), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Identity) makeResponseUserIdModules(oldData, newData []userIdModule) (responses []ResponseCompare, err error) {
	mapModuleOld := make(map[string]userIdModule)
	mapModuleNew := make(map[string]userIdModule)

	for _, module := range oldData {
		mapModuleOld[module.Name] = module
	}

	for _, module := range newData {
		mapModuleNew[module.Name] = module
	}

	for name, moduleOld := range mapModuleOld {
		if moduleNew, ok := mapModuleNew[name]; ok {
			// Param
			if moduleOld.Params != nil || moduleNew.Params != nil {
				mapParamModuleOld := make(map[string]paramModuleUserId)
				mapParamModuleNew := make(map[string]paramModuleUserId)

				if moduleOld.Params != nil {
					for _, param := range *moduleOld.Params {
						mapParamModuleOld[param.Name] = param
					}
				}
				if moduleNew.Params != nil {
					for _, param := range *moduleNew.Params {
						mapParamModuleNew[param.Name] = param
					}
				}
				for nameParam, paramOld := range mapParamModuleOld {
					// Kiểm tra nếu có tồn tại trong paramNew là update
					if paramNew, ok := mapParamModuleNew[nameParam]; ok {
						// Type
						res, err := makeResponseCompare("UserId Module > "+name+" > Param > "+nameParam+" > Type", &paramOld.Type, &paramNew.Type, mysql.TYPEObjectTypeUpdate)
						if err == nil {
							responses = append(responses, res)
						}
						// Template
						res, err = makeResponseCompare("UserId Module > "+name+" > Param > "+nameParam+" > Template", &paramOld.Template, &paramNew.Template, mysql.TYPEObjectTypeUpdate)
						if err == nil {
							responses = append(responses, res)
						}
						// Loại các param update để lại các param add mới trong mapParammoduleNew
						delete(mapParamModuleNew, nameParam)
					} else {
						// Nếu như param không tồn tại trong paramNew tức param này bị xóa
						// Type
						res, err := makeResponseCompare("UserId Module > "+name+" > Param > "+nameParam+" > Type", &paramOld.Type, nil, mysql.TYPEObjectTypeDel)
						if err == nil {
							responses = append(responses, res)
						}
						// Template
						res, err = makeResponseCompare("UserId Module > "+name+" > Param > "+nameParam+" > Template", &paramOld.Template, nil, mysql.TYPEObjectTypeDel)
						if err == nil {
							responses = append(responses, res)
						}
					}
				}
				// Sau khi đã loại các param đã tồn tại trong param old còn lại sẽ là các param mới
				for nameParam, paramNew := range mapParamModuleNew {
					// Type
					res, err := makeResponseCompare("UserId Module > "+name+" > Param > "+nameParam+" > Type", nil, &paramNew.Type, mysql.TYPEObjectTypeAdd)
					if err == nil {
						responses = append(responses, res)
					}
					// Template
					res, err = makeResponseCompare("UserId Module > "+name+" > Param > "+nameParam+" > Template", nil, &paramNew.Template, mysql.TYPEObjectTypeAdd)
					if err == nil {
						responses = append(responses, res)
					}
				}
			}
			// Storage
			if moduleOld.Storage != nil || moduleNew.Storage != nil {
				mapStorageModuleOld := make(map[string]storageModuleUserId)
				mapStorageModuleNew := make(map[string]storageModuleUserId)

				if moduleOld.Storage != nil {
					for _, storage := range *moduleOld.Storage {
						mapStorageModuleOld[storage.Name] = storage
					}
				}
				if moduleNew.Storage != nil {
					for _, storage := range *moduleNew.Storage {
						mapStorageModuleNew[storage.Name] = storage
					}
				}
				for nameStorage, storageOld := range mapStorageModuleOld {
					// Kiểm tra nếu có tồn tại trong storageNew là update
					if storageNew, ok := mapStorageModuleNew[nameStorage]; ok {
						// Type
						res, err := makeResponseCompare("UserId Module > "+name+" > Storage > "+nameStorage+" > Type", &storageOld.Type, &storageNew.Type, mysql.TYPEObjectTypeUpdate)
						if err == nil {
							responses = append(responses, res)
						}
						// Template
						res, err = makeResponseCompare("UserId Module > "+name+" > Storage > "+nameStorage+" > Template", &storageOld.Template, &storageNew.Template, mysql.TYPEObjectTypeUpdate)
						if err == nil {
							responses = append(responses, res)
						}
						// Loại các storage update để lại các storage add mới trong mapStorageModuleNew
						delete(mapStorageModuleNew, nameStorage)
					} else {
						// Nếu như param không tồn tại trong storageNew tức param này bị xóa
						// Type
						res, err := makeResponseCompare("UserId Module > "+name+" > Storage > "+nameStorage+" > Type", &storageOld.Type, nil, mysql.TYPEObjectTypeDel)
						if err == nil {
							responses = append(responses, res)
						}
						// Template
						res, err = makeResponseCompare("UserId Module > "+name+" > Storage > "+nameStorage+" > Template", &storageOld.Template, nil, mysql.TYPEObjectTypeDel)
						if err == nil {
							responses = append(responses, res)
						}
					}
				}
				// Sau khi đã loại các param đã tồn tại trong param old còn lại sẽ là các param mới
				for _, storage := range mapStorageModuleNew {
					// Type
					res, err := makeResponseCompare("UserId Module > "+name+" > Storage > "+storage.Name+" > Type", nil, &storage.Type, mysql.TYPEObjectTypeAdd)
					if err == nil {
						responses = append(responses, res)
					}
					// Template
					res, err = makeResponseCompare("UserId Module > "+name+" > Storage > "+storage.Name+" > Template", nil, &storage.Template, mysql.TYPEObjectTypeAdd)
					if err == nil {
						responses = append(responses, res)
					}
				}
			}
			// Xóa các module update
			delete(mapModuleNew, name)
		} else {
			// Nếu chỉ tồn tại module old tức module này đã bị xóa toàn bộ các newData là nil
			// Param
			if moduleOld.Params != nil {
				for _, param := range *moduleOld.Params {
					// Type
					res, err := makeResponseCompare("UserId Module > "+name+" > Param > "+param.Name+" > Type", &param.Type, nil, mysql.TYPEObjectTypeDel)
					if err == nil {
						responses = append(responses, res)
					}
					// Template
					res, err = makeResponseCompare("UserId Module > "+name+" > Param > "+param.Name+" > Template", &param.Template, nil, mysql.TYPEObjectTypeDel)
					if err == nil {
						responses = append(responses, res)
					}
				}
			}
			// Storage
			if moduleOld.Storage != nil {
				for _, storage := range *moduleOld.Storage {
					// Type
					res, err := makeResponseCompare("UserId Module > "+name+" > Storage > "+storage.Name+" > Type", &storage.Type, nil, mysql.TYPEObjectTypeDel)
					if err == nil {
						responses = append(responses, res)
					}
					// Template
					res, err = makeResponseCompare("UserId Module > "+name+" > Storage > "+storage.Name+" > Template", &storage.Template, nil, mysql.TYPEObjectTypeDel)
					if err == nil {
						responses = append(responses, res)
					}
				}
			}
		}
	}
	// Các moduleNew còn lại đều là các module add mới
	for name, moduleNew := range mapModuleNew {
		// Param
		if moduleNew.Params != nil {
			for _, param := range *moduleNew.Params {
				// Type
				res, err := makeResponseCompare("UserId Module > "+name+" > Param > "+param.Name+" > Type", nil, &param.Type, mysql.TYPEObjectTypeAdd)
				if err == nil {
					responses = append(responses, res)
				}
				// Template
				res, err = makeResponseCompare("UserId Module > "+name+" > Param > "+param.Name+" > Template", nil, &param.Template, mysql.TYPEObjectTypeAdd)
				if err == nil {
					responses = append(responses, res)
				}
			}
		}
		// Storage
		if moduleNew.Storage != nil {
			for _, storage := range *moduleNew.Storage {
				// Type
				res, err := makeResponseCompare("UserId Module > "+name+" > Storage > "+storage.Name+" > Type", nil, &storage.Type, mysql.TYPEObjectTypeAdd)
				if err == nil {
					responses = append(responses, res)
				}
				// Template
				res, err = makeResponseCompare("UserId Module > "+name+" > Storage > "+storage.Name+" > Template", nil, &storage.Template, mysql.TYPEObjectTypeAdd)
				if err == nil {
					responses = append(responses, res)
				}
			}
		}
	}
	return
}
