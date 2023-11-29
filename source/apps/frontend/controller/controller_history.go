package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"time"

	// "source/core/technology/mysql"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/pkg/ajax"
	"source/pkg/utility"
	"strings"
)

type History struct{}

type AssignHistoryIndex struct {
	assign.Schema
	Params    payload.HistoryIndex
	Pages     []string
	StartDate string
	EndDate   string
	Label     string
	IsWindow  bool
	// Histories []model.HistoryRecord
}

func (t *History) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	if userLogin.Id == 0 {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignHistoryIndex{Schema: assign.Get(ctx)}

	params := payload.HistoryIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns.Params = params

	if utility.IsWindow() {
		assigns.IsWindow = true
	}

	objectPages := new(model.History).GetAllHistoryPage()
	assigns.Pages = objectPages

	assigns.StartDate = time.Now().AddDate(0, 0, -100).Format("2006-01-02")
	assigns.EndDate = time.Now().Format("2006-01-02")
	assigns.Label = time.Now().Format("02 Jan, 2006") + " - " + time.Now().AddDate(0, 0, -100).Format("02 Jan, 2006")

	assigns.Title = config.TitleWithPrefix("Activity")
	return ctx.Render("history/index", assigns, view.LAYOUTMain)
}

func (t *History) Filter(ctx *fiber.Ctx) error {
	// Get data from model
	userLogin := GetUserLogin(ctx)
	if userLogin.Id == 0 {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	var inputs payload.HistoryFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// dataTable, err := new(model.History).GetByFilters(&inputs, userLogin)
	// if err != nil {
	// 	return err
	// }
	// return ctx.JSON(dataTable)
	return ctx.JSON("")
}

type AssignLoadHistory struct {
	Title     string
	Host      string
	Histories []model.ResponseHistory
	Total     int
}

func (t *History) LoadHistory(ctx *fiber.Ctx) error {
	assigns := AssignLoadHistory{}
	userLogin := GetUserLogin(ctx)
	if !userLogin.IsFound() {
		// return ctx.SendStatus(fiber.StatusNotFound)
	}
	// var responses ajax.Responses
	inputs := payload.LoadHistory{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	histories, _ := new(model.History).LoadHistories(inputs, userLogin.Id)
	assigns.Histories = histories
	assigns.Title = t.GetTitle(inputs.Object)
	assigns.Host = new(model.History).GetHost(histories)
	assigns.Total = len(histories)
	return ctx.Render("history/load_history", assigns, view.LAYOUTEmpty)
}

func (t *History) ObjectByPage(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	if !userLogin.IsFound() {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// var responses ajax.Responses
	inputs := payload.LoadObjectByPage{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	Objects, err := new(model.History).LoadObjectsByPage(inputs, userLogin.Id)
	if err != nil {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return ctx.JSON(ajax.Responses{
		Status:     "success",
		Message:    "",
		DataObject: Objects,
	})
}

func (t *History) GetTitle(object string) (title string) {
	title = strings.Title(strings.Replace(strings.Replace(object, "_", " ", -1), "fe", "", -1))
	// title = strings.Title(strings.Replace(title, "_", " ", -1))
	switch object {
	// case "line_item":
	// 	title = "Line item"
	// 	break
	case "inventory_config_fe":
		title = "Config"
		break
	case "inventory_consent_fe":
		title = "Consent"
		break
	case "inventory_connection_fe":
		title = "Connection"
		break
	case "inventory_adstxt_fe":
		title = "Ads.txt"
		break
	// case "floor_fe":
	// 	title = "Floor"
	// 	break
	// case "identity_fe":
	// 	title = "Identity"
	// 	break
	default:
		break
	}
	return
}
