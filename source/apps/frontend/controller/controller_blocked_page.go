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
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"strconv"
)

type BlockedPage struct {
}

type AssignBlockedPageAdd struct {
	assign.Schema
}

func (t *BlockedPage) Add(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockedPageAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignBlockedPageAdd{
		Schema: assign.Get(ctx),
	}
	assigns.Title = config.TitleWithPrefix("Add Blocked Page")
	return ctx.Render("blocked-page/add", assigns, view.LAYOUTMain)
}

func (t *BlockedPage) AddPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockedPageAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	resp := ajax.Responses{
		Status: "error",
	}

	// Get inputs
	inputs := payload.RuleSubmit{}

	if err := ctx.BodyParser(&inputs); err != nil {
		resp.Errors = append(resp.Errors, ajax.Error{
			Id:      "",
			Message: "error",
		})
		return ctx.JSON(resp)
	}

	// Create
	inputs.Type = mysql.TYPERuleTypeBlockedPage
	_, errs := new(model.BlockedPage).AddPost(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		resp.Errors = errs
		return ctx.JSON(resp)
	}

	// Trả về response thành công
	resp.Status = "success"

	return ctx.JSON(resp)
}

type AssignBlockedPageEdit struct {
	assign.Schema
	Row model.RuleRecord
}

func (t *BlockedPage) Edit(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockedPageEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignBlockedPageEdit{
		Schema: assign.Get(ctx),
	}
	// Get id
	id, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		return ctx.SendString("not found")
	}

	// Get user
	user := assigns.UserLogin

	// Get record
	assigns.Row, _ = new(model.BlockedPage).GetById(id, user.Id)
	if assigns.Row.ID < 1 {
		return ctx.SendString("not found")
	}

	assigns.Title = config.TitleWithPrefix("Edit Blocked Page")
	return ctx.Render("blocked-page/edit", assigns, view.LAYOUTMain)
}

func (t *BlockedPage) EditPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockedPageEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	resp := ajax.Responses{
		Status: "error",
	}

	// Get inputs
	inputs := payload.RuleSubmit{}

	if err := ctx.BodyParser(&inputs); err != nil {
		fmt.Println(err)
		resp.Errors = append(resp.Errors, ajax.Error{
			Id:      "",
			Message: "error",
		})
		return ctx.JSON(resp)
	}

	// Create
	inputs.Type = mysql.TYPERuleTypeBlockedPage
	_, errs := new(model.BlockedPage).EditPost(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		resp.Errors = errs
		return ctx.JSON(resp)
	}

	// Trả về response thành công
	resp.Status = "success"
	return ctx.JSON(resp)
}

func (t *BlockedPage) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockedPageDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	blockedPage := new(model.BlockedPage)
	notify := blockedPage.Delete(inputs.Id, userLogin.Id, userAdmin)
	return ctx.JSON(notify)
}

func (t *BlockedPage) ImportCSV(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockedPageImportCSV)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	resp := ajax.Responses{
		Status: "error",
	}

	// Get file từ query
	file, err := ctx.FormFile("file")
	if err != nil {
		resp.Errors = append(resp.Errors, ajax.Error{
			Id:      "",
			Message: "file error",
		})
		return ctx.JSON(resp)
	}

	// Save file
	path, errs := new(model.BlockedPage).SaveFileCSV(ctx, file)
	if len(errs) > 0 {
		resp.Errors = errs
		return ctx.JSON(resp)
	}

	// Parser file
	data, errs := new(model.BlockedPage).ParserCSV(path)
	if len(errs) > 0 {
		resp.Errors = errs
		return ctx.JSON(resp)
	}

	// Trả về response thành công
	resp.Status = "success"
	resp.DataObject = data
	return ctx.JSON(resp)
}
