package history

import (
	"encoding/json"
	"source/core/technology/mysql"
)

type Channel struct {
	Detail    DetailChannel
	CreatorId int64
	RecordOld mysql.TableChannels
	RecordNew mysql.TableChannels
}

func (t *Channel) Page() string {
	return "Channels"
}

type DetailChannel int

const (
	DetailChannelFE DetailChannel = iota + 1
	DetailChannelBE
)

func (t DetailChannel) String() string {
	switch t {
	case DetailChannelFE:
		return "channel_fe"
	case DetailChannelBE:
		return "channel_be"
	}
	return ""
}

func (t DetailChannel) App() string {
	switch t {
	case DetailChannelFE:
		return "FE"
	case DetailChannelBE:
		return "BE"
	}
	return ""
}

func (t *Channel) Type() TYPEHistory {
	return TYPEHistoryChannel
}

func (t *Channel) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *Channel) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailChannelFE:
		return t.getHistoryChannelFE()
	case DetailChannelBE:
		return t.getHistoryChannelBE()
	}
	return mysql.TableHistory{}
}

func (t *Channel) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailChannelFE.String():
		return t.compareDataChannelFE(history)
	case DetailChannelBE.String():
		return t.compareDataChannelFE(history)
	}
	return []ResponseCompare{}
}

func (t *Channel) getRootRecord() (record mysql.TableChannels) {
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

type channelFE struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Category    *string   `json:"category,omitempty"`
	Keyword     *[]string `json:"keyword,omitempty"`
	Language    *string   `json:"language,omitempty"`
}

func (t *Channel) getHistoryChannelFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := channelFE{}
	newData := channelFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Channels,
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
		history.Title = "Add Channel"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Channel"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Channel"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *channelFE) MakeData(record mysql.TableChannels) {
	rec.Name = &record.Name
	rec.Description = &record.Description
	var category mysql.TableCategory
	mysql.Client.Find(&category, record.Category)
	rec.Category = &category.Name

	var keywords []string
	for _, keyword := range record.Keywords {
		keywords = append(keywords, keyword.Keyword)
	}
	if len(keywords) > 0 {
		rec.Keyword = &keywords
	}

	var language mysql.TableLanguage
	mysql.Client.Find(&language, record.Language)
	rec.Language = &language.LanguageName
}

func (t *Channel) compareDataChannelFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew channelFE
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
	// Category
	res, err = makeResponseCompare("Category", recordOld.Category, recordNew.Category, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Keyword
	res, err = makeResponseCompare("Keyword", pointerArrayStringToString(recordOld.Keyword), pointerArrayStringToString(recordNew.Keyword), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Language
	res, err = makeResponseCompare("Language", recordOld.Language, recordNew.Language, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Channel) getHistoryChannelBE() (history mysql.TableHistory) {
	history = t.getHistoryChannelFE()
	return
}

func (t *Channel) compareDataChannelBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareDataChannelFE(history)
	return
}
