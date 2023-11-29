package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/pkg/ajax"
	"source/pkg/utility"
	"strconv"
)

type Content struct{}

type AssignContent struct {
	assign.Schema
	Categories []model.CategoryRecord
	Channels   []model.ChannelsRecord
	Params     payload.ContentIndex
}

type AssignContentAdd struct {
	assign.Schema
	Categories []model.CategoryRecord
	Channels   []model.ChannelsRecord
	Params     payload.ContentAdd
}

type AssignContentEdit struct {
	assign.Schema
	Row model.ContentRecord
	// Categories       []model.CategoryRecord
	Channels         []model.ChannelsRecord
	Keyword          []model.ContentKeywordRecord
	Tags             []model.ContentTagRecord
	DomainUpload     string
	ListSelectThumb  []string
	CheckImgSelected bool
	ListAdBreak      []model.ContentAdBreakRecord
	AdBreakPreRoll   model.ContentAdBreakRecord
	AdBreakMidRoll   []model.ContentAdBreakRecord
	AdBreakPostRoll  model.ContentAdBreakRecord
}

type AssignContentEditQuiz struct {
	assign.Schema
	Row          model.ContentRecord
	Categories   []model.CategoryRecord
	Tags         []model.ContentTagRecord
	DomainUpload string
	Questions    []payload.Questions
}

func (t *Content) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContent)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.ContentIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignContent{Schema: assign.Get(ctx)}
	// assigns.Categories = new(model.Category).GetAll()
	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("List Content")
	assigns.LANG.Title = "Content"
	return ctx.Render("content/index", assigns, view.LAYOUTMain)
}

func (t *Content) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContent)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.ContentFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.Content).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		return err
	}
	return ctx.JSON(dataTable)
}

func (t *Content) Add(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignContentAdd{Schema: assign.Get(ctx)}
	params := payload.ContentAdd{}
	channels := ctx.Query("f_channels")
	if channels != "" {
		selectChannelsId, err := strconv.ParseInt(channels, 10, 64)
		if err != nil {
			return err
		}
		params.Channels = selectChannelsId
		assigns.Params = params
	}

	assigns.Channels = new(model.Channels).GetAll(userLogin.Id)
	assigns.Title = config.TitleWithPrefix("Add Content")
	return ctx.Render("content/add", assigns, view.LAYOUTMain)
}

func (t *Content) AddPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.ContentCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Content).Create(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

func (t *Content) Edit(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignContentEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	row, err := new(model.Content).GetById(idSearch, userLogin.Id)
	if err != nil || row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.ContentError.NotFound.ToString())
	}
	// Update uuid cho các row chưa có
	if row.Uuid == "" {
		row.RenderUuid()
	}
	// assigns.Categories = new(model.Category).GetAll()
	assigns.Channels = new(model.Channels).GetAll(userLogin.Id)
	assigns.Keyword = new(model.ContentKeyword).GetByContent(row.Id)
	assigns.Tags = new(model.ContentTag).GetByContent(row.Id)
	var domainUpload string
	if utility.IsWindow() {
		// domainUpload = "http://127.0.0.1:8543"
		domainUpload = "https://ul.pubpowerplatform.io"
	} else {
		domainUpload = "https://ul.pubpowerplatform.io"
	}
	assigns.DomainUpload = domainUpload
	assigns.CheckImgSelected = false
	for i := 1; i <= 5; i++ {
		nameThumbDefault := "/assets/img/thumb_" + row.NameFile + "_" + strconv.Itoa(i) + ".png"
		if nameThumbDefault == row.Thumb {
			assigns.CheckImgSelected = true
		}
		assigns.ListSelectThumb = append(assigns.ListSelectThumb, nameThumbDefault)
	}
	assigns.Row = row
	for _, adBreak := range row.AdBreaks {
		switch adBreak.Type {
		case "preroll":
			assigns.AdBreakPreRoll = model.ContentAdBreakRecord{TableContentAdBreak: adBreak}
			break
		case "midroll":
			assigns.AdBreakMidRoll = append(assigns.AdBreakMidRoll, model.ContentAdBreakRecord{TableContentAdBreak: adBreak})
			break
		case "postroll":
			assigns.AdBreakPostRoll = model.ContentAdBreakRecord{TableContentAdBreak: adBreak}
			break
		}
	}
	assigns.ListAdBreak = new(model.ContentAdBreak).GetByContent(row.Id)
	assigns.Title = config.TitleWithPrefix("Edit Content: " + row.Title)
	return ctx.Render("content/edit", assigns, view.LAYOUTMain)
}

func (t *Content) EditPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	// }}
	// Get Post Data
	inputs := payload.ContentCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Content).Update(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

func (t *Content) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	content := new(model.Content)
	notify := content.Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}

func (t *Content) AddQuiz(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentAddQuiz)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// tạm thời bỏ quiz vì chưa có js quiz
	assigns := AssignContent{Schema: assign.Get(ctx)}
	assigns.Categories = new(model.Category).GetAll()
	assigns.Title = config.TitleWithPrefix("Add Content")
	return ctx.Render("content/add-quiz", assigns, view.LAYOUTMain)
}

func (t *Content) AddQuizPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentAddQuiz)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// tạm thời bỏ quiz vì chưa có js quiz
	return ctx.SendStatus(fiber.StatusNotFound)

	// Get Post Data
	inputs := payload.ContentQuizCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Content).CreateQuiz(inputs, userLogin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

func (t *Content) EditQuiz(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentEditQuiz)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// assigns := AssignContent{Schema: assign.Get(ctx)}
	// assigns.Categories = new(model.Category).GetAll()
	// assigns.Title = config.TitleWithPrefix("Edit Content")
	// return ctx.Render("content/edit-quiz", assigns, view.LAYOUTMain)

	assigns := AssignContentEditQuiz{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	row, err := new(model.Content).GetById(idSearch, userLogin.Id)
	if err != nil || row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.ContentError.NotFound.ToString())
	}
	assigns.Categories = new(model.Category).GetAll()
	assigns.Tags = new(model.ContentTag).GetByContent(row.Id)
	var domainUpload string
	if utility.IsWindow() {
		domainUpload = "http://127.0.0.1:8543"
	} else {
		domainUpload = "https://ul.pubpowerplatform.io"
	}
	assigns.DomainUpload = domainUpload
	assigns.Row = row
	Questions := new(model.ContentQuestion).GetByContent(row.Id)
	assigns.Questions = new(model.ContentQuestion).HandleQuestion(Questions)
	assigns.Title = config.TitleWithPrefix("Edit Content: " + row.Title)
	return ctx.Render("content/edit-quiz", assigns, view.LAYOUTMain)
}

func (t *Content) EditQuizPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIContentEditQuiz)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.ContentQuizCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}

	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Content).UpdateQuiz(inputs, userLogin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}
