package history

import (
	"encoding/json"
	"source/core/technology/mysql"
)

type Floor struct {
	Detail    DetailFloor
	CreatorId int64
	RecordOld mysql.TableFloor
	RecordNew mysql.TableFloor
}

func (t *Floor) Page() string {
	return "Rule"
}

type DetailFloor int

const (
	DetailFloorFE DetailFloor = iota + 1
)

func (t DetailFloor) String() string {
	switch t {
	case DetailFloorFE:
		return "floor_fe"
	}
	return ""
}

func (t DetailFloor) App() string {
	switch t {
	case DetailFloorFE:
		return "FE"
	}
	return ""
}

func (t *Floor) Type() TYPEHistory {
	return TYPEHistoryFloor
}

func (t *Floor) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *Floor) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailFloorFE:
		return t.getHistoryFloorFE()
	}
	return mysql.TableHistory{}
}

func (t *Floor) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailFloorFE.String():
		return t.compareDataFloorFE(history)
	}
	return []ResponseCompare{}
}

func (t *Floor) getRootRecord() (record mysql.TableFloor) {
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

type floorFE struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Priority    *int     `json:"priority,omitempty"`
	Status      *string  `json:"status,omitempty"`
	FloorType   *string  `json:"floor_type,omitempty"`
	FloorValue  *float64 `json:"floor_value,omitempty"`
	Target      target
}

func (t *Floor) getHistoryFloorFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := floorFE{}
	newData := floorFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Floor,
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
		history.Title = "Add Floor"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Floor"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Floor"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *floorFE) MakeData(record mysql.TableFloor) {
	rec.Name = &record.Name
	rec.Description = &record.Description
	rec.Priority = &record.Priority

	// Xử lý target
	rec.Target = makeTarget(record.Targets)

	var floorType string
	if record.FloorType == 1 {
		floorType = "Dynamic Price Floor"
		// Bỏ các target không dùng trong dynamic
		rec.Target.Format = nil
		rec.Target.Size = nil
		rec.Target.AdTag = nil
		rec.Target.Geography = nil
		rec.Target.Device = nil

	} else if record.FloorType == 2 {
		floorType = "Hard Price Floor"
		rec.FloorValue = &record.FloorValue
	}
	rec.FloorType = &floorType

}

func (t *Floor) compareDataFloorFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew floorFE
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
	// FloorType
	res, err = makeResponseCompare("Floor Type", recordOld.FloorType, recordNew.FloorType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Floor Value
	res, err = makeResponseCompare("Floor Value", pointerFloatToString(recordOld.FloorValue), pointerFloatToString(recordNew.FloorValue), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
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
