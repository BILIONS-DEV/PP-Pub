package controllerv2

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/view"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/pkg/utility"
	"strconv"
)

type quiz struct {
	*handler
}

func (h *handler) InitRoutesQuiz(app fiber.Router) {
	t := quiz{handler: h}
	app.Get(config.URIContentQuiz, t.Index)
	app.Post(config.URIContentQuiz, t.Filter)
	app.Get(config.URIContentAddQuiz, t.Add)
	app.Post(config.URIContentAddQuiz, t.AddPost)
	app.Get(config.URIContentEditQuiz, t.Edit)
	app.Post(config.URIContentEditQuiz, t.EditPost)
	app.Post(config.URIContentDelQuiz, t.Delete)
}

type AssignQuiz struct {
	Assign
	Params dto.QuizIndex
}

func (t *quiz) Index(ctx *fiber.Ctx) (err error) {

	params := dto.QuizIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignQuiz{Assign: newAssign(ctx, "Quiz")}
	assigns.Params = params
	return ctx.Render("quiz/index", assigns, view.LAYOUTMain)
}

func (t *quiz) Filter(ctx *fiber.Ctx) (err error) {
	// get user login
	userLogin := getUserLogin(ctx)

	// Get payload post
	var inputs dto.QuizFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}

	// Get data from model
	dataTable, err := t.useCases.Quiz.GetByFilters(inputs, userLogin)
	if err != nil {
		return err
	}
	return ctx.JSON(dataTable)
}

type AssignQuizAdd struct {
	Assign
	Categories []model.Category
}

func (t *quiz) Add(ctx *fiber.Ctx) (err error) {

	// assigns
	assigns := AssignQuizAdd{Assign: newAssign(ctx, "Add Quiz")}
	assigns.Categories, _ = t.useCases.Category.GetAll()
	return ctx.Render("content/add-quiz", assigns, view.LAYOUTMain)
}

func (t *quiz) AddPost(ctx *fiber.Ctx) (err error) {
	// get user login
	userLogin := getUserLogin(ctx)
	userAdmin := getUserAdmin(ctx)

	// get payload from http request
	var payload dto.PayloadQuizSubmit
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.JSON(dto.MakeResponseError(err))
	}
	fmt.Printf("%+v\n", payload)
	// run usecase
	response := t.useCases.Quiz.Create(payload, userLogin, userAdmin)

	return ctx.JSON(response)
}

type AssignQuizEdit struct {
	Assign
	Row          model.Quiz
	Categories   []model.Category
	DomainUpload string
}

func (t *quiz) Edit(ctx *fiber.Ctx) (err error) {
	id, _ := strconv.ParseInt(ctx.Query("id"), 10, 64)

	// assigns
	assigns := AssignQuizEdit{Assign: newAssign(ctx, "Edit Quiz")}
	assigns.DomainUpload = "https://ul.vlitag.com"
	if utility.IsWindow() {
		assigns.DomainUpload = "http://127.0.0.1:8543"
	}
	assigns.Row, _ = t.useCases.Quiz.GetById(id)
	assigns.Categories, _ = t.useCases.Category.GetAll()
	return ctx.Render("content/edit-quiz", assigns, view.LAYOUTMain)
}

func (t *quiz) EditPost(ctx *fiber.Ctx) (err error) {
	// get user login
	userLogin := getUserLogin(ctx)
	userAdmin := getUserAdmin(ctx)

	// get payload from http request
	var payload dto.PayloadQuizSubmit
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.JSON(dto.MakeResponseError(err))
	}

	// run usecase
	response := t.useCases.Quiz.Update(payload, userLogin, userAdmin)

	return ctx.JSON(response)
}

func (t *quiz) Delete(ctx *fiber.Ctx) (err error) {
	// get user login
	userLogin := getUserLogin(ctx)
	userAdmin := getUserAdmin(ctx)

	inputs := dto.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	notify, _ := t.useCases.Quiz.Delete(inputs, userLogin, userAdmin)
	return ctx.JSON(notify)
}
