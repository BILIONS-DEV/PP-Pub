package history

import (
	"encoding/json"
	"source/core/technology/mysql"
)

type ABTest struct {
	Detail    DetailABTest
	CreatorId int64
	RecordOld mysql.TableAbTesting
	RecordNew mysql.TableAbTesting
}

func (t *ABTest) Page() string {
	return "A/B Test"
}

type DetailABTest int

const (
	DetailABTestFE DetailABTest = iota + 1
)

func (t DetailABTest) String() string {
	switch t {
	case DetailABTestFE:
		return "ab_test_fe"
	}
	return ""
}

func (t DetailABTest) App() string {
	switch t {
	case DetailABTestFE:
		return "FE"
	}
	return ""
}

func (t *ABTest) Type() TYPEHistory {
	return TYPEHistoryABTest
}

func (t *ABTest) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *ABTest) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailABTestFE:
		return t.getHistoryABTestFE()
	}
	return mysql.TableHistory{}
}

func (t *ABTest) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailABTestFE.String():
		return t.compareDataABTestFE(history)
	}
	return []ResponseCompare{}
}

func (t *ABTest) getRootRecord() (record mysql.TableAbTesting) {
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

type abTestFE struct {
	Name              *string  `json:"name,omitempty"`
	Description       *string  `json:"description,omitempty"`
	TestType          *string  `json:"test_type,omitempty"`
	DynamicPriceFloor *string  `json:"dynamic_price_floor,omitempty"`
	HardPriceFloor    *float64 `json:"hard_price_floor,omitempty"`
	Status            *string  `json:"status,omitempty"`
	StartDate         *string  `json:"start_date,omitempty"`
	EndDate           *string  `json:"end_date,omitempty"`
	Target            target   `json:"target,omitempty"`
}

func (t *ABTest) getHistoryABTestFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := abTestFE{}
	newData := abTestFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.AbTesting,
		ObjectName: t.getRootRecord().Name,
		ObjectId:   t.getRootRecord().Id,
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
		history.Title = "Add A/B Test"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update A/B Test"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete A/B Test"
		history.OldData = string(bOldData)
	}
	return
}

func (t *ABTest) compareDataABTestFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew abTestFE
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
	// Test Type
	res, err = makeResponseCompare("Test Type", recordOld.TestType, recordNew.TestType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Dynamic Price Floor
	res, err = makeResponseCompare("Dynamic Price Floor", recordOld.DynamicPriceFloor, recordNew.DynamicPriceFloor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Hard Price Floor
	res, err = makeResponseCompare("Hard Price Floor", pointerFloatToString(recordOld.HardPriceFloor), pointerFloatToString(recordNew.HardPriceFloor), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Status
	res, err = makeResponseCompare("Status", recordOld.Status, recordNew.Status, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Start Date
	res, err = makeResponseCompare("Start Date", recordOld.StartDate, recordNew.StartDate, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// End Date
	res, err = makeResponseCompare("End Date", recordOld.EndDate, recordNew.EndDate, history.ObjectType)
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

func (rec *abTestFE) MakeData(record mysql.TableAbTesting) {
	rec.Name = &record.Name
	rec.Description = &record.Description
	testType := record.TestType.String()
	rec.TestType = &testType
	if record.TestType == mysql.TYPETestTypeDynamicHardPriceFloor {
		var dynamicFloorRule mysql.TableFloor
		mysql.Client.Where("id = ? ", record.DynamicFloorPrice).Find(&dynamicFloorRule)
		rec.DynamicPriceFloor = &dynamicFloorRule.Name
		rec.HardPriceFloor = &record.HardPriceFloor
		status := record.Status.String()
		rec.Status = &status
		var startDate, endDate string
		layoutISO := "01/02/2006"
		if record.StartDate.Valid {
			startDate = record.StartDate.Time.Format(layoutISO)
		} else {
			startDate = "Immediately"
		}
		if record.EndDate.Valid {
			endDate = record.EndDate.Time.Format(layoutISO)
		} else {
			endDate = "Unlimited"
		}
		rec.StartDate = &startDate
		rec.EndDate = &endDate
	}
}
