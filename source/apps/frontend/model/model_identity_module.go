package model

import (
	"encoding/json"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"strconv"
	"strings"
	// "fmt"
)

type IdentityModuleInfo struct{}

type IdentityModuleInfoRecord struct {
	mysql.TableIdentityModuleInfo
}

func (IdentityModuleInfoRecord) TableName() string {
	return mysql.Tables.IdentityModuleInfo
}

func (t *IdentityModuleInfo) MakeRecord(inputs []payload.ModuleInfo, identityId int64) (rows []IdentityModuleInfoRecord, err error) {
	var row IdentityModuleInfoRecord
	for _, val := range inputs {
		row.ModuleId = val.ModuleId
		row.Name = val.ModuleName
		row.IdentityId = identityId
		params, err := json.Marshal(val.Params)
		if err != nil {
			return nil, err
		}
		storages, err := json.Marshal(val.Storage)
		if err != nil {
			return nil, err
		}
		row.Params = string(params)
		row.Storage = string(storages)
		row.AbTesting = val.AbTesting
		row.Volume = val.Volume
		rows = append(rows, row)
	}
	return
}

func (t *IdentityModuleInfo) CreateModuleInfo(record IdentityModuleInfoRecord) (err error) {
	err = mysql.Client.Create(&record).Error
	return
}

func (t *IdentityModuleInfo) DeleteModuleInfo(identityId, moduleId int64) (err error) {
	err = mysql.Client.Where(IdentityModuleInfoRecord{mysql.TableIdentityModuleInfo{IdentityId: identityId, ModuleId: moduleId}}).Delete(&IdentityModuleInfoRecord{}).Error
	return
}

func (t *IdentityModuleInfo) GetModules(identityId int64) (records []IdentityModuleInfoRecord, listIdOfInvent []int64, err error) {
	err = mysql.Client.Where("identity_id = ?", identityId).Order("module_id").Find(&records).Error
	if err != nil {
		return
	}
	for _, rec := range records {
		listIdOfInvent = append(listIdOfInvent, rec.ModuleId)
	}
	return
}

func (t *IdentityModuleInfo) ArrayInt64ToString(identityId []int64) (listIdOfInvent string) {
	if identityId == nil {
		return
	}

	var IDs []string
	for _, i := range identityId {
		IDs = append(IDs, strconv.FormatInt(i, 10))
	}
	listIdOfInvent = strings.Join(IDs, ",")
	return
}

func (t *IdentityModuleInfo) CheckExist(identityId, moduleId int64) (row IdentityModuleInfoRecord) {
	mysql.Client.Where("inventory_id = ? and module_id = ?", identityId, moduleId).Find(&row)
	return
}

func (t *IdentityModuleInfo) MakeRowDefault(identityId int64) {
	listModule := []payload.DefaultModule{
		{
			Id:      1,
			Name:    "flocId",
			Params:  "[{\"name\":\"token\",\"type\":\"string\",\"template\":\"A3dHTSoNUMjjERBLlrvJSelNnwWUCwVQhZ5tNQ+sll7y+LkPPVZXtB77u2y7CweRIxiYaGw GXNlW1/dFp8VMEgIAAAB+eyJvcmlnaW4iOiJodHRwczovL3NoYXJlZGlkLm9yZzo0NDMiLC JmZWF0dXJlIjoiSW50ZXJlc3RDb2hvcnRBUEkiLCJleHBpcnkiOjE2MjYyMjA3OTksImlzU 3ViZG9tYWluIjp0cnVlLCJpc1RoaXJkUGFydHkiOnRydWV9\"}]",
			Storage: "[]",
		},
		{
			Id:      3,
			Name:    "pubCommonId",
			Params:  "[]",
			Storage: "[{\"name\":\"type\",\"type\":\"string\",\"template\":\"cookie\"},{\"name\":\"name\",\"type\":\"string\",\"template\":\"_pubcid\"},{\"name\":\"expires\",\"type\":\"int\",\"template\":\"365\"}]",
		},
		{
			Id:      4,
			Name:    "criteo",
			Params:  "[]",
			Storage: "[]",
		},
		{
			Id:      5,
			Name:    "id5Id",
			Params:  "[{\"name\":\"partner\",\"type\":\"int\",\"template\":\"\"},{\"name\":\"pd\",\"type\":\"string\",\"template\":\"\"}]",
			Storage: "[{\"name\":\"type\",\"type\":\"string\",\"template\":\"html5\"},{\"name\":\"name\",\"type\":\"string\",\"template\":\"id5id\"},{\"name\":\"expires\",\"type\":\"int\",\"template\":\"90\"},{\"name\":\"refreshInSeconds\",\"type\":\"int\",\"template\":\"8*3600\"}]",
		},
	}
	for _, module := range listModule {
		moduleInfo := IdentityModuleInfoRecord{mysql.TableIdentityModuleInfo{
			IdentityId: identityId,
			ModuleId:    module.Id,
			Name:        module.Name,
			Params:      module.Params,
			Storage:     module.Storage,
			AbTesting:   1,
			Volume:      50,
		}}
		mysql.Client.Model(IdentityModuleInfoRecord{}).Create(&moduleInfo)
	}
}

func (t *IdentityModuleInfo) DelModuleInfo(domainId int64) {
	mysql.Client.Model(&IdentityModuleInfoRecord{}).Delete(&IdentityModuleInfoRecord{}, "inventory_id = ?", domainId)
}
