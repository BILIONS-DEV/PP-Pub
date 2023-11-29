package history

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"source/core/technology/mysql"
)

type Blocking struct {
	Detail    DetailBlocking
	CreatorId int64
	RecordOld mysql.TableBlocking
	RecordNew mysql.TableBlocking
}

func (t *Blocking) Page() string {
	return "Rule"
}

type DetailBlocking int

const (
	DetailBlockingFE DetailBlocking = iota + 1
)

func (t DetailBlocking) String() string {
	switch t {
	case DetailBlockingFE:
		return "blocking_fe"
	}
	return ""
}

func (t DetailBlocking) App() string {
	switch t {
	case DetailBlockingFE:
		return "FE"
	}
	return ""
}

func (t *Blocking) Type() TYPEHistory {
	return TYPEHistoryBlocking
}

func (t *Blocking) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *Blocking) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailBlockingFE:
		return t.getHistoryBlockingFE()
	}
	return mysql.TableHistory{}
}

func (t *Blocking) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailBlockingFE.String():
		return t.compareDataBlockingFE(history)
	}
	return []ResponseCompare{}
}

func (t *Blocking) getRootRecord() (record mysql.TableBlocking) {
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

type blockingFE struct {
	RestrictionName *string      `json:"restriction_name,omitempty"`
	Inventory       *[]string    `json:"inventory,omitempty"`
	Restrictions    Restrictions `json:"restrictions,omitempty"`
}

type Restrictions struct {
	AdvertiserDomains []string
	Creative          []string
}

func (t *Blocking) getHistoryBlockingFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := blockingFE{}
	newData := blockingFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Blocking,
		ObjectId:   t.getRootRecord().Id,
		ObjectName: t.getRootRecord().RestrictionName,
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
		history.Title = "Add Blocking"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Blocking"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Blocking"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *blockingFE) MakeData(record mysql.TableBlocking) {
	rec.RestrictionName = &record.RestrictionName
	var inventories []string
	for _, inventory := range record.Inventories {
		inventories = append(inventories, inventory.Name)
	}
	if len(inventories) > 0 {
		rec.Inventory = &inventories
	}

	for _, restriction := range record.Restrictions {
		if !govalidator.IsNull(restriction.AdvertiseDomain) {
			rec.Restrictions.AdvertiserDomains = append(rec.Restrictions.AdvertiserDomains, restriction.AdvertiseDomain)
		}
		if !govalidator.IsNull(restriction.CreativeId) {
			rec.Restrictions.Creative = append(rec.Restrictions.Creative, restriction.CreativeId)
		}
	}
}

func (t *Blocking) compareDataBlockingFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew blockingFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Restriction Name
	res, err := makeResponseCompare("Restriction Name", recordOld.RestrictionName, recordNew.RestrictionName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Inventory
	res, err = makeResponseCompare("Inventory", pointerArrayStringToString(recordOld.Inventory), pointerArrayStringToString(recordNew.Inventory), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Restrictions Advertiser Domains
	res, err = makeResponseCompare("Restrictions > Advertiser Domains", pointerArrayStringToString(&recordOld.Restrictions.AdvertiserDomains), pointerArrayStringToString(&recordNew.Restrictions.AdvertiserDomains), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Restrictions Creative
	res, err = makeResponseCompare("Restrictions > Creative", pointerArrayStringToString(&recordOld.Restrictions.Creative), pointerArrayStringToString(&recordNew.Restrictions.Creative), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}
