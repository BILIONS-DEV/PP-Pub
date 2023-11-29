package history

import (
	"encoding/json"
	"errors"
	"source/core/technology/mysql"
	"strings"
)

type Playlist struct {
	Detail    DetailPlaylist
	CreatorId int64
	RecordOld mysql.TablePlaylist
	RecordNew mysql.TablePlaylist
}

func (t *Playlist) Page() string {
	return "Playlist"
}

type DetailPlaylist int

const (
	DetailPlaylistFE DetailPlaylist = iota + 1
	DetailPlaylistBE
)

func (t DetailPlaylist) String() string {
	switch t {
	case DetailPlaylistFE:
		return "playlist_fe"
	case DetailPlaylistBE:
		return "playlist_be"
	}
	return ""
}

func (t DetailPlaylist) App() string {
	switch t {
	case DetailPlaylistFE:
		return "FE"
	case DetailPlaylistBE:
		return "BE"
	}
	return ""
}

func (t *Playlist) Type() TYPEHistory {
	return TYPEHistoryPlaylist
}

func (t *Playlist) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}
func (t *Playlist) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailPlaylistFE:
		return t.getHistoryPlaylistFE()
	case DetailPlaylistBE:
		return t.getHistoryPlaylistBE()
	}
	return mysql.TableHistory{}
}

func (t *Playlist) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailPlaylistFE.String():
		return t.compareDataPlaylistFE(history)
	case DetailPlaylistBE.String():
		return t.compareDataPlaylistFE(history)
	}
	return []ResponseCompare{}
}

func (t *Playlist) getRootRecord() (record mysql.TablePlaylist) {
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

type playlistFE struct {
	Name              *string           `json:"name,omitempty"`
	Description       *string           `json:"description,omitempty"`
	OrderingMethod    *string           `json:"ordering_method"`
	VideosLimit       *int64            `json:"videos_limit"`
	ChannelsAndVideos ChannelsAndVideos `json:"channels_and_videos"`
}

type ChannelsAndVideos struct {
	Language map[string][]string `json:"language"`
	Channels map[string][]string `json:"channels"`
	Category map[string][]string `json:"category"`
	Keywords map[string][]string `json:"keywords"`
	Videos   map[string][]string `json:"videos"`
}

func (t *Playlist) getHistoryPlaylistFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := playlistFE{}
	newData := playlistFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Playlist,
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
		history.Title = "Add Playlist"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Playlist"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Playlist"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *playlistFE) MakeData(record mysql.TablePlaylist) {
	rec.Name = &record.Name
	rec.Description = &record.Description
	rec.OrderingMethod = &record.OrderingMethod
	rec.VideosLimit = &record.VideosLimit

	// Tạo map
	mapLanguage := make(map[string][]string)
	mapChannel := make(map[string][]string)
	mapCategory := make(map[string][]string)
	mapKeyword := make(map[string][]string)
	mapVideo := make(map[string][]string)

	for _, channelAndVideo := range record.ChannelsAndVideos {
		typ := "include"
		if channelAndVideo.Type == 1 {
			typ = "include"
		} else if channelAndVideo.Type == 2 {
			typ = "exclude"
		}
		if channelAndVideo.LanguageId != 0 {
			if channelAndVideo.LanguageId == -1 {
				mapLanguage[typ] = append(mapLanguage[typ], "all")
			} else {
				var language mysql.TableLanguage
				mysql.Client.Find(&language, channelAndVideo.LanguageId)
				mapLanguage[typ] = append(mapLanguage[typ], language.LanguageName)
			}
		}
		if channelAndVideo.ChannelsId != 0 {
			if channelAndVideo.ChannelsId == -1 {
				mapChannel[typ] = append(mapChannel[typ], "all")
			} else {
				var channels mysql.TableChannels
				mysql.Client.Find(&channels, channelAndVideo.ChannelsId)
				mapChannel[typ] = append(mapChannel[typ], channels.Name)
			}
		}
		if channelAndVideo.CategoryId != 0 {
			if channelAndVideo.CategoryId == -1 {
				mapCategory[typ] = append(mapCategory[typ], "all")
			} else {
				var category mysql.TableCategory
				mysql.Client.Find(&category, channelAndVideo.CategoryId)
				mapCategory[typ] = append(mapCategory[typ], category.Name)
			}
		}
		if channelAndVideo.ContentKeywordId != 0 {
			if channelAndVideo.ContentKeywordId == -1 {
				mapKeyword[typ] = append(mapKeyword[typ], "all")
			} else {
				var keyword mysql.TableContentKeyword
				mysql.Client.Find(&keyword, channelAndVideo.ContentKeywordId)
				mapKeyword[typ] = append(mapKeyword[typ], keyword.Keyword)
			}
		}
		if channelAndVideo.ContentId != 0 {
			if channelAndVideo.ContentId == -1 {
				mapVideo[typ] = append(mapVideo[typ], "all")
			} else {
				var content mysql.TableContent
				mysql.Client.Find(&content, channelAndVideo.ContentId)
				mapVideo[typ] = append(mapVideo[typ], content.Title)
			}
		}
	}

	rec.ChannelsAndVideos.Language = mapLanguage
	rec.ChannelsAndVideos.Channels = mapChannel
	rec.ChannelsAndVideos.Category = mapCategory
	rec.ChannelsAndVideos.Keywords = mapKeyword
	rec.ChannelsAndVideos.Videos = mapVideo
}

func (t *Playlist) compareDataPlaylistFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew playlistFE
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
	// Ordering Method
	res, err = makeResponseCompare("Ordering Method", recordOld.OrderingMethod, recordNew.OrderingMethod, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Videos Limit
	res, err = makeResponseCompare("Videos Limit", pointerInt64ToString(recordOld.VideosLimit), pointerInt64ToString(recordNew.VideosLimit), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Channels And Videos Language
	responseLanguage, err := t.makeResponseCompareChannelsAndVideos(recordOld.ChannelsAndVideos.Language, recordNew.ChannelsAndVideos.Language, "Language")
	if err == nil {
		responses = append(responses, responseLanguage...)
	}
	// Channels And Videos Channels
	responseChannels, err := t.makeResponseCompareChannelsAndVideos(recordOld.ChannelsAndVideos.Channels, recordNew.ChannelsAndVideos.Channels, "Channels")
	if err == nil {
		responses = append(responses, responseChannels...)
	}
	// Channels And Videos Category
	responseCategory, err := t.makeResponseCompareChannelsAndVideos(recordOld.ChannelsAndVideos.Category, recordNew.ChannelsAndVideos.Category, "Category")
	if err == nil {
		responses = append(responses, responseCategory...)
	}
	// Channels And Videos Keywords
	responseKeywords, err := t.makeResponseCompareChannelsAndVideos(recordOld.ChannelsAndVideos.Keywords, recordNew.ChannelsAndVideos.Keywords, "Keywords")
	if err == nil {
		responses = append(responses, responseKeywords...)
	}
	// Channels And Videos Videos
	responseVideos, err := t.makeResponseCompareChannelsAndVideos(recordOld.ChannelsAndVideos.Videos, recordNew.ChannelsAndVideos.Videos, "Videos")
	if err == nil {
		responses = append(responses, responseVideos...)
	}
	return
}

func (t *Playlist) makeResponseCompareChannelsAndVideos(oldData, newData map[string][]string, detail string) (responses []ResponseCompare, err error) {
	if len(oldData) == 0 && len(newData) == 0 {
		err = errors.New("no data channels and videos")
		return
	}
	// Trường hợp chỉ có newData là trường hợp add
	if len(oldData) == 0 && len(newData) > 0 {
		for typ, target := range newData {
			res, err := makeResponseCompare("Channels And Videos > "+strings.Title(detail)+" > "+strings.Title(typ), nil, pointerArrayStringToString(&target), mysql.TYPEObjectTypeAdd)
			if err == nil {
				responses = append(responses, res)
			}
		}
	}
	// Trường hợp có cả 2 data là trường hợp update
	if len(oldData) > 0 && len(newData) > 0 {
		for typeOld, targetOld := range oldData {
			targetNew, existNewData := newData[typeOld]
			if existNewData {
				// Nếu type include hoặc exclude cũng có trong newData thì đây là trường trường hợp update
				res, err := makeResponseCompare("Channels And Videos > "+strings.Title(detail)+" > "+strings.Title(typeOld), pointerArrayStringToString(&targetOld), pointerArrayStringToString(&targetNew), mysql.TYPEObjectTypeUpdate)
				if err == nil {
					responses = append(responses, res)
				}
			} else {
				// Nếu như không tồn tại cùng type trong newData đây là trường hợp đổi type in ra delete oldData và add NewData
				res, err := makeResponseCompare("Channels And Videos > "+strings.Title(detail)+" > "+strings.Title(typeOld), pointerArrayStringToString(&targetOld), nil, mysql.TYPEObjectTypeDel)
				if err == nil {
					responses = append(responses, res)
				}
				var typeNew string
				if typeOld == "include" {
					typeNew = "exclude"
				} else if typeOld == "exclude" {
					typeNew = "include"
				}
				res, err = makeResponseCompare("Channels And Videos > "+strings.Title(detail)+" > "+strings.Title(typeNew), nil, pointerArrayStringToString(&targetNew), mysql.TYPEObjectTypeAdd)
				if err == nil {
					responses = append(responses, res)
				}
			}
		}
	}
	// Trường hợp chỉ có oldData là trường hợp delete
	if len(oldData) > 0 && len(newData) == 0 {
		for typ, target := range oldData {
			res, err := makeResponseCompare("Channels And Videos > "+strings.Title(detail)+" > "+strings.Title(typ), pointerArrayStringToString(&target), nil, mysql.TYPEObjectTypeDel)
			if err == nil {
				responses = append(responses, res)
			}
		}
	}
	return
}

func (t *Playlist) getHistoryPlaylistBE() (history mysql.TableHistory) {
	history = t.getHistoryPlaylistFE()
	return
}

func (t *Playlist) compareDataPlaylistBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareDataPlaylistFE(history)
	return
}
