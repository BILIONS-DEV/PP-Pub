package history

import (
	"encoding/json"
	"fmt"
	"source/core/technology/mysql"
	"strconv"
)

type GAM struct {
	Detail    DetailGAM
	CreatorId int64
	RecordOld mysql.TableGamNetwork
	RecordNew mysql.TableGamNetwork
}

func (t *GAM) Page() string {
	return "GAM"
}

type DetailGAM int

const (
	DetailGAMSignInFE DetailGAM = iota + 1
	DetailGAMSetupFE
	DetailGAMSignInBE
	DetailGAMSetupBE
)

func (t DetailGAM) String() string {
	switch t {
	case DetailGAMSignInFE:
		return "gam_sign_in_fe"
	case DetailGAMSetupFE:
		return "gam_setup_fe"
	case DetailGAMSignInBE:
		return "gam_sign_in_be"
	case DetailGAMSetupBE:
		return "gam_setup_be"
	}
	return ""
}

func (t DetailGAM) App() string {
	switch t {
	case DetailGAMSetupFE, DetailGAMSignInFE:
		return "FE"
	case DetailGAMSetupBE, DetailGAMSignInBE:
		return "BE"
	}
	return ""
}

func (t *GAM) Type() TYPEHistory {
	return TYPEHistoryGAM
}

func (t *GAM) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *GAM) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailGAMSignInFE:
		return t.getHistoryGAMSignInFE()
	case DetailGAMSetupFE:
		return t.getHistoryGAMSetupFE()
	case DetailGAMSignInBE:
		return t.getHistoryGAMSignInBE()
	case DetailGAMSetupBE:
		return t.getHistoryGAMSetupBE()
	}
	return mysql.TableHistory{}
}

func (t *GAM) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailGAMSignInFE.String():
		return t.compareDataGAMSignInFE(history)
	case DetailGAMSetupFE.String():
		return t.compareDataGAMSetupFE(history)
	case DetailGAMSignInBE.String():
		return t.compareDataGAMSignInBE(history)
	case DetailGAMSetupBE.String():
		return t.compareDataGAMSetupBE(history)
	}
	return []ResponseCompare{}
}

func (t *GAM) getRootRecord() (record mysql.TableGamNetwork) {
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

type gamFE struct {
	Id       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	TimeZone string `json:"time_zone"`
}

func (t *GAM) getHistoryGAMSignInFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := gamFE{}
	newData := gamFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.GamNetwork,
		ObjectId:   t.getRootRecord().Id,
		ObjectName: t.getRootRecord().NetworkName,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserId,
	}
	var bNewData, bOldData []byte
	if t.RecordNew.Id != 0 {
		newData = gamFE{
			Id:       t.RecordNew.NetworkId,
			Name:     t.RecordNew.NetworkName,
			Currency: t.RecordNew.CurrencyCode,
			TimeZone: t.RecordNew.TimeZone,
		}
		bNewData, _ = json.Marshal(newData)
	}
	if t.RecordOld.Id != 0 {
		oldData = gamFE{
			Id:       t.RecordOld.NetworkId,
			Name:     t.RecordOld.NetworkName,
			Currency: t.RecordOld.CurrencyCode,
			TimeZone: t.RecordOld.TimeZone,
		}
		bOldData, _ = json.Marshal(oldData)
	}
	if t.Action() == mysql.TYPEObjectTypeAdd {
		history.Title = "Add GAM"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update GAM"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete GAM"
		history.OldData = string(bOldData)
	}
	return
}

func (t *GAM) getHistoryGAMSetupFE() (history mysql.TableHistory) {
	// Validate
	if t.RecordNew.Id == 0 {
		fmt.Println("GAM new is require")
		return
	}

	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.GamNetwork,
		ObjectName: t.getRootRecord().NetworkName,
		ObjectId:   t.getRootRecord().Id,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserId,
	}

	// Đưa dữ liệu vào struct để tạo json
	title := ""
	if t.RecordOld.Status != t.RecordNew.Status {
		if t.RecordNew.Status == mysql.StatusGamPending {
			title = "Unselected GAM"
		} else if t.RecordNew.Status == mysql.StatusGamSelected {
			title = "Selected GAM"
		}
	} else if t.RecordOld.ApiAccess != t.RecordNew.ApiAccess {
		title = "Check Api Access"
	} else if t.RecordOld.PushLineItem != t.RecordNew.PushLineItem {
		if t.RecordOld.PushLineItem == 0 {
			title = "Push Line Item"
		} else if t.RecordOld.PushLineItem == 2 {
			title = "Repush Line Item"
		}
	}
	history.Title = title
	return
}

func (t *GAM) compareDataGAMSignInFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew gamFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Network Id
	networkOld := strconv.FormatInt(recordOld.Id, 10)
	networkNew := strconv.FormatInt(recordOld.Id, 10)
	res, err := makeResponseCompare("Network Id", &networkOld, &networkNew, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Name
	res, err = makeResponseCompare("Name", &recordOld.Name, &recordNew.Name, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Currency
	res, err = makeResponseCompare("Currency", &recordOld.Currency, &recordNew.Currency, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// TimeZone
	res, err = makeResponseCompare("Time Zone", &recordOld.TimeZone, &recordNew.TimeZone, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *GAM) compareDataGAMSetupFE(history mysql.TableHistory) (responses []ResponseCompare) {
	// Xử lý compare từng row

	// Title
	responses = append(responses, ResponseCompare{
		Action:  "update",
		Text:    history.Title,
		OldData: "",
		NewData: "",
	})
	return
}

func (t *GAM) getHistoryGAMSignInBE() (history mysql.TableHistory) {
	history = t.getHistoryGAMSignInFE()
	return
}

func (t *GAM) getHistoryGAMSetupBE() (history mysql.TableHistory) {
	history = t.getHistoryGAMSetupFE()
	return
}

func (t *GAM) compareDataGAMSignInBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareDataGAMSignInFE(history)
	return
}

func (t *GAM) compareDataGAMSetupBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareDataGAMSetupFE(history)
	return
}
