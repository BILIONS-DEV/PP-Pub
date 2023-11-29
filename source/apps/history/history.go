package history

import (
	"errors"
	"fmt"
	"reflect"
	"source/core/technology/mysql"
	"strconv"
	"strings"
)

type HistorySchema interface {
	Page() string
	Type() TYPEHistory
	Action() mysql.TYPEObjectType
	Data() mysql.TableHistory
	CompareData(history mysql.TableHistory) (res []ResponseCompare)
}

type TYPEHistory int

const (
	TYPEHistoryUser TYPEHistory = iota + 1
	TYPEHistoryInventory
	TYPEHistoryAdTag
	TYPEHistoryBidder
	TYPEHistoryLineItem
	TYPEHistoryIdentity
	TYPEHistoryFloor
	TYPEHistoryGAM
	TYPEHistoryABTest
	TYPEHistoryBlocking
	TYPEHistoryChannel
	TYPEHistoryContent
	TYPEHistoryPlaylist
	TYPEHistoryTemplate
	TYPEHistoryBidderTemplate
	TYPEHistoryModuleUserId
	TYPEHistoryBlockedPage
)

func PushHistory(object interface{}) (err error) {
	// Tạo value cho interface với reflect
	rvRec := reflect.ValueOf(object)

	// Xử lý pointer cho value
	if rvRec.Kind() != reflect.Ptr {
		err = errors.New("object require is pointer")
		return
	}

	// Get tableName
	if schema, ok := rvRec.Interface().(HistorySchema); ok {
		data := schema.Data()
		data.Page = schema.Page()
		// Thử compare data nếu như không có thay đổi nào bỏ qua history
		responses, err := CompareDataHistory(data)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if len(responses) == 0 {
			return errors.New("nothing changes")
		}
		// Create history từ data
		err = mysql.Client.Create(&data).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		err = errors.New("missing interface for struct")
		return
	}
	return
}

type ResponseCompare struct {
	Action  string
	Text    string
	OldData string
	NewData string
}

func CompareDataHistory(history mysql.TableHistory) (res []ResponseCompare, err error) {
	typeHistory := stringToTypeHistory(history.Object)
	var historySchema HistorySchema
	switch typeHistory {
	case TYPEHistoryUser:
		historySchema = &User{}
		break
	case TYPEHistoryInventory:
		historySchema = &Inventory{}
		break
	case TYPEHistoryAdTag:
		historySchema = &AdTag{}
		break
	case TYPEHistoryBidder:
		historySchema = &Bidder{}
		break
	case TYPEHistoryLineItem:
		historySchema = &LineItem{}
		break
	case TYPEHistoryIdentity:
		historySchema = &Identity{}
		break
	case TYPEHistoryFloor:
		historySchema = &Floor{}
		break
	case TYPEHistoryGAM:
		historySchema = &GAM{}
		break
	case TYPEHistoryABTest:
		historySchema = &ABTest{}
		break
	case TYPEHistoryBlocking:
		historySchema = &Blocking{}
		break
	case TYPEHistoryChannel:
		historySchema = &Channel{}
		break
	case TYPEHistoryContent:
		historySchema = &Content{}
		break
	case TYPEHistoryPlaylist:
		historySchema = &Playlist{}
		break
	case TYPEHistoryTemplate:
		historySchema = &Template{}
		break
	case TYPEHistoryBidderTemplate:
		historySchema = &BidderTemplate{}
		break
	case TYPEHistoryModuleUserId:
		historySchema = &ModuleUserId{}
		break
	case TYPEHistoryBlockedPage:
		historySchema = &BlockedPage{}
		break
	default:
		return res, errors.New("not found type history")
	}
	isAction := false
	responses := historySchema.CompareData(history)
	for _, reps := range responses {
		if reps.Action != "none" {
			isAction = true
		}
	}
	if isAction {
		res = responses
	}
	return
}

type Tabler interface {
	TableName() string
}

func DefaultCompareData(object interface{}, detailType interface{}) (res []ResponseCompare, err error) {
	// Tạo value cho interface với reflect
	rvRec := reflect.ValueOf(object)

	// Get tableName
	if schema, ok := rvRec.Interface().(Tabler); ok {
		table := schema.TableName()
		typeHistory := stringToTypeHistory(table)
		var historySchema HistorySchema
		switch typeHistory {
		case TYPEHistoryUser:
			recordNew, ok := object.(mysql.TableUser)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailUser)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &User{
				Detail:    detail,
				RecordOld: mysql.TableUser{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryInventory:
			recordNew, ok := object.(mysql.TableInventory)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailInventory)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Inventory{
				Detail:    detail,
				RecordOld: mysql.TableInventory{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryAdTag:
			recordNew, ok := object.(mysql.TableInventoryAdTag)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailAdTag)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &AdTag{
				Detail:    detail,
				RecordOld: mysql.TableInventoryAdTag{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryBidder:
			recordNew, ok := object.(mysql.TableBidder)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailBidder)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Bidder{
				Detail:    detail,
				RecordOld: mysql.TableBidder{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryLineItem:
			recordNew, ok := object.(mysql.TableLineItem)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailLineItem)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &LineItem{
				Detail:    detail,
				RecordOld: mysql.TableLineItem{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryIdentity:
			recordNew, ok := object.(mysql.TableIdentity)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailIdentity)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Identity{
				Detail:    detail,
				RecordOld: mysql.TableIdentity{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryFloor:
			recordNew, ok := object.(mysql.TableFloor)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailFloor)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Floor{
				Detail:    detail,
				RecordOld: mysql.TableFloor{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryGAM:
			recordNew, ok := object.(mysql.TableGamNetwork)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailGAM)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &GAM{
				Detail:    detail,
				RecordOld: mysql.TableGamNetwork{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryABTest:
			recordNew, ok := object.(mysql.TableAbTesting)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailABTest)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &ABTest{
				Detail:    detail,
				RecordOld: mysql.TableAbTesting{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryBlocking:
			recordNew, ok := object.(mysql.TableBlocking)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailBlocking)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Blocking{
				Detail:    detail,
				RecordOld: mysql.TableBlocking{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryChannel:
			recordNew, ok := object.(mysql.TableChannels)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailChannel)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Channel{
				Detail:    detail,
				RecordOld: mysql.TableChannels{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryContent:
			recordNew, ok := object.(mysql.TableContent)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailContent)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Content{
				Detail:    detail,
				RecordOld: mysql.TableContent{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryPlaylist:
			recordNew, ok := object.(mysql.TablePlaylist)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailPlaylist)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Playlist{
				Detail:    detail,
				RecordOld: mysql.TablePlaylist{},
				RecordNew: recordNew,
			}
			break
		case TYPEHistoryTemplate:
			recordNew, ok := object.(mysql.TableTemplate)
			if !ok {
				err = errors.New("object isn't table " + table)
				return
			}
			detail, ok := detailType.(DetailTemplate)
			if !ok {
				err = errors.New("detail type not accept")
				return
			}
			historySchema = &Template{
				Detail:    detail,
				RecordOld: mysql.TableTemplate{},
				RecordNew: recordNew,
			}
			break
		default:
			return res, errors.New("not found type history")
		}
		res = historySchema.CompareData(historySchema.Data())
	} else {
		err = errors.New("missing interface for struct")
		return
	}
	return
}

func stringToTypeHistory(object string) (typ TYPEHistory) {
	object = strings.TrimSpace(object)
	if object == mysql.Tables.AbTesting {
		typ = TYPEHistoryABTest
	} else if object == mysql.Tables.InventoryAdTag {
		typ = TYPEHistoryAdTag
	} else if object == mysql.Tables.Bidder || object == mysql.Tables.RlsConnectionMCM {
		typ = TYPEHistoryBidder
	} else if object == mysql.Tables.Blocking {
		typ = TYPEHistoryBlocking
	} else if object == mysql.Tables.Channels {
		typ = TYPEHistoryChannel
	} else if object == mysql.Tables.Content {
		typ = TYPEHistoryContent
	} else if object == mysql.Tables.LineItem {
		typ = TYPEHistoryLineItem
	} else if object == mysql.Tables.Floor {
		typ = TYPEHistoryFloor
	} else if object == mysql.Tables.GamNetwork {
		typ = TYPEHistoryGAM
	} else if object == mysql.Tables.Identity {
		typ = TYPEHistoryIdentity
	} else if object == mysql.Tables.Inventory || object == mysql.Tables.InventoryConfig || object == mysql.Tables.InventoryConnectionDemand {
		typ = TYPEHistoryInventory
	} else if object == mysql.Tables.Playlist {
		typ = TYPEHistoryPlaylist
	} else if object == mysql.Tables.Template {
		typ = TYPEHistoryTemplate
	} else if object == mysql.Tables.User || object == mysql.Tables.UserBilling {
		typ = TYPEHistoryUser
	} else if object == mysql.Tables.BidderTemplate {
		typ = TYPEHistoryBidderTemplate
	} else if object == mysql.Tables.ModuleUserId {
		typ = TYPEHistoryModuleUserId
	} else if object == mysql.Tables.Rule {
		typ = TYPEHistoryBlockedPage
	}
	return
}

func makeResponseCompare(text string, recOldData *string, recNewData *string, objectType mysql.TYPEObjectType) (res ResponseCompare, err error) {
	if recOldData == nil && recNewData == nil {
		return res, errors.New("no response")
	}
	var oldData, newData string
	if recOldData != nil {
		oldData = *recOldData
	} else {
		oldData = ""
	}
	if recNewData != nil {
		newData = *recNewData
	} else {
		newData = ""
	}
	var action string
	switch objectType {
	case mysql.TYPEObjectTypeAdd:
		action = "add"
		break
	case mysql.TYPEObjectTypeUpdate:
		action = compareStringData(recOldData, recNewData)
	case mysql.TYPEObjectTypeDel:
		action = "delete"
	}
	res = ResponseCompare{
		Action:  action,
		Text:    text,
		OldData: oldData,
		NewData: newData,
	}
	return
}

func compareStringData(oldData *string, newData *string) (action string) {
	if oldData == nil && newData != nil {
		action = "add"
		return
	}
	if oldData != nil && newData == nil {
		action = "delete"
		return
	}
	if oldData == nil && newData == nil {
		action = "none"
		return
	}

	if oldData != nil && newData != nil {
		if *oldData != *newData {
			action = "update"
			return
		}
		if *oldData == *newData {
			action = "none"
			return
		}
	}

	return
}

func pointerFloatToString(data *float64) *string {
	if data != nil {
		str := fmt.Sprintf("%f", *data)
		return &str
	}
	return nil
}

func pointerIntToString(data *int) *string {
	if data != nil {
		str := strconv.Itoa(*data)
		return &str
	}
	return nil
}

func pointerInt64ToString(data *int64) *string {
	if data != nil {
		str := strconv.FormatInt(*data, 10)
		return &str
	}
	return nil
}

func pointerArrayStringToString(data *[]string) *string {
	if data != nil {
		str := strings.Join(*data, ", ")
		return &str
	}
	return nil
}
