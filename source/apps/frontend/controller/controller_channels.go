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
	"strconv"
)

type Channels struct{}

type AssignChannels struct {
	assign.Schema
	Rows       []model.ChannelsIndex
	Categories []model.CategoryRecord
	Languages  []model.LanguageRecord
	Params     payload.ChannelsIndex
}

type AssignChannelsEdit struct {
	assign.Schema
	Row        model.ChannelsRecord
	Categories []model.CategoryRecord
	Languages  []model.LanguageRecord
	Keyword    []model.ChannelsKeywordRecord
}

func (t *Channels) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChannels)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.ChannelsIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignChannels{Schema: assign.Get(ctx)}
	assigns.Categories = new(model.Category).GetAll()
	assigns.Params = params
	// Get payload post
	var inputs payload.ChannelsFilterPayload
	if err := ctx.QueryParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	channels, err := new(model.Channels).GetByFilters(&inputs, userLogin.Id)
	if err != nil {
		return err
	}
	assigns.Rows = channels
	assigns.Title = config.TitleWithPrefix("List Channels")
	assigns.LANG.Title = "Channels"
	return ctx.Render("channels/index", assigns, view.LAYOUTMain)
}

func (t *Channels) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChannels)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.ChannelsFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.Channels).GetByFilters(&inputs, userLogin.Id)
	if err != nil {
		return err
	}
	return ctx.JSON(dataTable)
}

func (t *Channels) Add(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChannelsAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignChannels{Schema: assign.Get(ctx)}
	assigns.Categories = new(model.Category).GetAll()
	assigns.Languages = new(model.Language).GetAll()
	assigns.Title = config.TitleWithPrefix("Add Channels")
	return ctx.Render("channels/add", assigns, view.LAYOUTMain)
}

func (t *Channels) AddPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChannelsAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.ChannelsCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Channels).Create(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

func (t *Channels) Edit(ctx *fiber.Ctx) error {
	assigns := AssignChannelsEdit{Schema: assign.Get(ctx)}
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChannelsEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id := ctx.Query("id")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	row, err := new(model.Channels).GetById(idSearch, userLogin.Id)
	if err != nil || row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.ChannelsError.NotFound.ToString())
	}
	assigns.Categories = new(model.Category).GetAll()
	assigns.Languages = new(model.Language).GetAll()
	assigns.Keyword = new(model.ChannelsKeyword).GetByChannels(row.Id)
	assigns.Row = row
	// assigns.ListAdBreak = new(model.ContentAdBreak).GetByContent(row.Id)
	assigns.Title = config.TitleWithPrefix("Edit Channels: " + row.Name)
	return ctx.Render("channels/edit", assigns, view.LAYOUTMain)
}

func (t *Channels) EditPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChannelsEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.ChannelsCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Channels).Update(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

func (t *Channels) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChannelsDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	channels := new(model.Channels)
	notify := channels.Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}
