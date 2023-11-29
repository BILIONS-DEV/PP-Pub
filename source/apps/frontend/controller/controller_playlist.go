package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/pkg/ajax"
	"strconv"
)

type Playlist struct{}

type AssignPlaylist struct {
	assign.Schema
	Tags   []model.InventoryAdTagRecord
	Params payload.PlaylistIndex
}

func (t *Playlist) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylist)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.PlaylistIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignPlaylist{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Playlist")
	return ctx.Render("playlist/index", assigns, view.LAYOUTMain)
}

func (t *Playlist) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylist)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.PlaylistFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.Playlist).GetByFilters(&inputs, userLogin.Id)
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

type AssignPlaylistAdd struct {
	assign.Schema
	ListBoxCollapse []string
	// Keywords        []model.ContentKeywordRecord
	// Contents        []model.ContentRecord
	// Channels        []model.ChannelsRecord
	// Categories      []model.CategoryRecord
}

func (t *Playlist) Add(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylistAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignPlaylistAdd{Schema: assign.Get(ctx)}
	// assigns.Keywords = new(model.ContentKeyword).GetByUserId(user.Id)
	// assigns.Channels = new(model.Channels).GetAll(user.Id)
	// assigns.Categories = new(model.Category).GetAll()
	// assigns.Contents = new(model.Content).HandleContents(new(model.Content).GetAll(user.Id))
	assigns.ListBoxCollapse = new(model.Playlist).GetListBoxCollapse(userLogin.Id, 0, "playlist", "add")
	assigns.Title = config.TitleWithPrefix("Add Playlist")
	return ctx.Render("playlist/add", assigns, view.LAYOUTMain)
}

func (t *Playlist) AddPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylistAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.PlaylistCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Playlist).Create(inputs, userLogin.Id, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

type AssignPlaylistEdit struct {
	assign.Schema
	ListBoxCollapse []string
	// Tags            []model.InventoryAdTagRecord
	Row model.PlaylistRecord

	// playlist config
	Videos     []model.ContentRecord
	Channels   []model.ChannelsRecord
	Categories []model.CategoryRecord
	Language   []model.LanguageRecord
	Keywords   []model.ContentKeywordRecord

	VideosOfPlaylist   []int64
	ChannelsOfPlaylist []int64
	CategoryOfPlaylist []int64
	LanguageOfPlaylist []int64
	KeywordsOfPlaylist []int64

	TotalVideos  int
	TypeVideos   int64
	TypeChannels int64
	TypeCategory int64
	TypeLanguage int64
	TypeKeywords int64
}

func (t *Playlist) Edit(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylistEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignPlaylistEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idPlaylist, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	recPlaylist, err := new(model.Playlist).GetById(idPlaylist, userLogin.Id)
	if err != nil {
		return err
	}
	assigns.Row = recPlaylist

	PlaylistConfig := new(model.PlaylistConfig).GetPlaylistConfigPlaylist(idPlaylist)
	mapLanguage := make(map[int64]int)
	mapChannels := make(map[int64]int)
	mapCategory := make(map[int64]int)
	mapVideos := make(map[int64]int)
	mapKeywords := make(map[int64]int)
	for _, config := range PlaylistConfig {
		if config.LanguageId != 0 {
			mapLanguage[config.LanguageId] = 1
			assigns.TypeLanguage = config.Type
		}
		if config.ChannelsId != 0 {
			mapChannels[config.ChannelsId] = 1
			assigns.TypeChannels = config.Type
		}
		if config.CategoryId != 0 {
			mapCategory[config.CategoryId] = 1
			assigns.TypeCategory = config.Type
		}
		if config.ContentId != 0 {
			mapVideos[config.ContentId] = 1
			assigns.TypeVideos = config.Type
		}
		if config.ContentKeywordId != 0 {
			mapKeywords[config.ContentKeywordId] = 1
			assigns.TypeKeywords = config.Type
		}
	}

	// Append list sau khi đã lọc bỏ những id trùng nhau với map
	for languageId, _ := range mapLanguage {
		assigns.LanguageOfPlaylist = append(assigns.LanguageOfPlaylist, languageId)
	}
	for channelsId, _ := range mapChannels {
		assigns.ChannelsOfPlaylist = append(assigns.ChannelsOfPlaylist, channelsId)
	}
	for categoryId, _ := range mapCategory {
		assigns.CategoryOfPlaylist = append(assigns.CategoryOfPlaylist, categoryId)
	}
	for contentId, _ := range mapVideos {
		assigns.VideosOfPlaylist = append(assigns.VideosOfPlaylist, contentId)
	}
	for keywordId, _ := range mapKeywords {
		assigns.KeywordsOfPlaylist = append(assigns.KeywordsOfPlaylist, keywordId)
	}

	// recordsRlPlaylistContents := new(model.RlPlaylistContent).GetByPlaylistId(idPlaylist)
	// var contentSelected []model.ContentRecord
	// for _, recordsRlPlaylistContent := range recordsRlPlaylistContents {
	// 	content, _ := new(model.Content).GetById(recordsRlPlaylistContent.ContentId, userLogin.Id)
	// 	if content.Id == 0 {
	// 		continue
	// 	}
	// 	contentSelected = append(contentSelected, content)
	// }
	// assigns.ContentSelected = new(model.Content).HandleContents(contentSelected)
	assigns.Videos = new(model.Content).GetAllVideo(userLogin.Id)
	assigns.TotalVideos = len(assigns.Videos)
	assigns.Channels = new(model.Channels).GetAll(userLogin.Id)
	assigns.Categories = new(model.Category).GetAll()
	assigns.Language = new(model.Language).GetAll()
	assigns.Keywords = new(model.ContentKeyword).GetByUserId(userLogin.Id)
	assigns.ListBoxCollapse = new(model.Playlist).GetListBoxCollapse(userLogin.Id, idPlaylist, "playlist", "edit")
	assigns.Title = config.TitleWithPrefix("Edit Playlist")
	return ctx.Render("playlist/edit", assigns, view.LAYOUTMain)
}

func (t *Playlist) EditPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylistEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.PlaylistCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Playlist).Update(inputs, userLogin.Id, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

func (t *Playlist) View(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylistView)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignPlaylistEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idPlaylist, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	recPlaylist, err := new(model.Playlist).GetById(idPlaylist, userLogin.Id)
	if err != nil {
		return err
	}
	assigns.Row = recPlaylist

	PlaylistConfig := new(model.PlaylistConfig).GetPlaylistConfigPlaylist(idPlaylist)
	mapLanguage := make(map[int64]int)
	mapChannels := make(map[int64]int)
	mapCategory := make(map[int64]int)
	mapVideos := make(map[int64]int)
	mapKeywords := make(map[int64]int)
	for _, config := range PlaylistConfig {
		if config.LanguageId != 0 {
			mapLanguage[config.LanguageId] = 1
			assigns.TypeLanguage = config.Type
		}
		if config.ChannelsId != 0 {
			mapChannels[config.ChannelsId] = 1
			assigns.TypeChannels = config.Type
		}
		if config.CategoryId != 0 {
			mapCategory[config.CategoryId] = 1
			assigns.TypeCategory = config.Type
		}
		if config.ContentId != 0 {
			mapVideos[config.ContentId] = 1
			assigns.TypeVideos = config.Type
		}
		if config.ContentKeywordId != 0 {
			mapKeywords[config.ContentKeywordId] = 1
			assigns.TypeKeywords = config.Type
		}
	}

	// Append list sau khi đã lọc bỏ những id trùng nhau với map
	for languageId, _ := range mapLanguage {
		assigns.LanguageOfPlaylist = append(assigns.LanguageOfPlaylist, languageId)
	}
	for channelsId, _ := range mapChannels {
		assigns.ChannelsOfPlaylist = append(assigns.ChannelsOfPlaylist, channelsId)
	}
	for categoryId, _ := range mapCategory {
		assigns.CategoryOfPlaylist = append(assigns.CategoryOfPlaylist, categoryId)
	}
	for contentId, _ := range mapVideos {
		assigns.VideosOfPlaylist = append(assigns.VideosOfPlaylist, contentId)
	}
	for keywordId, _ := range mapKeywords {
		assigns.KeywordsOfPlaylist = append(assigns.KeywordsOfPlaylist, keywordId)
	}

	// recordsRlPlaylistContents := new(model.RlPlaylistContent).GetByPlaylistId(idPlaylist)
	// var contentSelected []model.ContentRecord
	// for _, recordsRlPlaylistContent := range recordsRlPlaylistContents {
	// 	content, _ := new(model.Content).GetById(recordsRlPlaylistContent.ContentId, userLogin.Id)
	// 	if content.Id == 0 {
	// 		continue
	// 	}
	// 	contentSelected = append(contentSelected, content)
	// }
	// assigns.ContentSelected = new(model.Content).HandleContents(contentSelected)
	assigns.Videos = new(model.Content).GetAllVideo(recPlaylist.UserId)
	assigns.TotalVideos = len(assigns.Videos)
	assigns.Channels = new(model.Channels).GetAll(recPlaylist.UserId)
	assigns.Categories = new(model.Category).GetAll()
	assigns.Language = new(model.Language).GetAll()
	assigns.Keywords = new(model.ContentKeyword).GetByUserId(recPlaylist.UserId)
	assigns.ListBoxCollapse = new(model.Playlist).GetListBoxCollapse(userLogin.Id, idPlaylist, "playlist", "edit")
	assigns.Title = config.TitleWithPrefix("View Playlist")
	return ctx.Render("playlist/view", assigns, view.LAYOUTMain)
}

func (t *Playlist) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylistDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	m := new(model.Playlist)
	notify := m.Delete(inputs.Id, userLogin.Id, userAdmin)
	return ctx.JSON(notify)
}

func (t *Playlist) Collapse(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlaylistCollapse)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := model.PageCollapseRecord{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	inputs.UserId = userLogin.Id
	inputs.PageCollapse = "playlist"
	errs := new(model.PageCollapse).HandleCollapse(inputs)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}
