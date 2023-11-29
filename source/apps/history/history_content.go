package history

import (
	"encoding/json"
	"source/core/technology/mysql"
)

type Content struct {
	Detail    DetailContent
	CreatorId int64
	RecordOld mysql.TableContent
	RecordNew mysql.TableContent
}

func (t *Content) Page() string {
	return "Content"
}

type DetailContent int

const (
	DetailContentFE DetailContent = iota + 1
	DetailContentBE
)

func (t DetailContent) String() string {
	switch t {
	case DetailContentFE:
		return "content_fe"
	case DetailContentBE:
		return "content_be"
	}
	return ""
}

func (t DetailContent) App() string {
	switch t {
	case DetailContentFE:
		return "FE"
	case DetailContentBE:
		return "BE"
	}
	return ""
}

func (t *Content) Type() TYPEHistory {
	return TYPEHistoryContent
}

func (t *Content) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *Content) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailContentFE:
		return t.getHistoryContentFE()
	case DetailContentBE:
		return t.getHistoryContentBE()
	}
	return mysql.TableHistory{}
}

func (t *Content) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailContentFE.String():
		return t.compareDataContentFE(history)
	case DetailContentBE.String():
		return t.compareDataContentBE(history)
	}
	return []ResponseCompare{}
}

func (t *Content) getRootRecord() (record mysql.TableContent) {
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

type contentFE struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Video       *string   `json:"video,omitempty"`
	Thumb       *string   `json:"thumb,omitempty"`
	Channel     *string   `json:"channel,omitempty"`
	Keyword     *[]string `json:"keyword,omitempty"`
	VideoType   *string   `json:"video_type,omitempty"`
}

func (t *Content) getHistoryContentFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := contentFE{}
	newData := contentFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Content,
		ObjectId:   t.getRootRecord().Id,
		ObjectName: t.getRootRecord().Title,
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
		history.Title = "Add Content"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Content"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Content"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *contentFE) MakeData(record mysql.TableContent) {
	rec.Title = &record.Title
	rec.Description = &record.ContentDesc
	rec.Video = &record.VideoUrl
	rec.Thumb = &record.Thumb

	var channel mysql.TableChannels
	mysql.Client.Find(&channel, record.Channels)
	rec.Channel = &channel.Name

	var keywords []string
	for _, keyword := range record.Keywords {
		keywords = append(keywords, keyword.Keyword)
	}
	if len(keywords) > 0 {
		rec.Keyword = &keywords
	}

	var videoType string
	if record.VideoType == 1 {
		videoType = "Public"
	} else {
		videoType = "Private"
	}
	rec.VideoType = &videoType

}

func (t *Content) compareDataContentFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew contentFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Title
	res, err := makeResponseCompare("Title", recordOld.Title, recordNew.Title, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Description
	res, err = makeResponseCompare("Account Type", recordOld.Description, recordNew.Description, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Video
	res, err = makeResponseCompare("Video", recordOld.Video, recordNew.Video, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Thumb
	res, err = makeResponseCompare("Thumb", recordOld.Thumb, recordNew.Thumb, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Channel
	res, err = makeResponseCompare("Channel", recordOld.Channel, recordNew.Channel, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Keyword
	res, err = makeResponseCompare("Keyword", pointerArrayStringToString(recordOld.Keyword), pointerArrayStringToString(recordNew.Keyword), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Video Type
	res, err = makeResponseCompare("Video Type", recordOld.VideoType, recordNew.VideoType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Content) getHistoryContentBE() (history mysql.TableHistory) {
	history = t.getHistoryContentFE()
	return
}

func (t *Content) compareDataContentBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareDataContentFE(history)
	return
}
