package controller

import (
	// "encoding/json"
	// "fmt"
	"github.com/gofiber/fiber/v2"
	// "source/apps/worker/push_line_item/model"
	"source/apps/frontend/model"
	// "source/apps/frontend/config"
	// "source/apps/frontend/config/assign"
	// "source/apps/frontend/payload"
	// "source/apps/frontend/view"
	// "source/pkg/ajax"
	// "source/pkg/utility"
	"strconv"
)

type Notification struct{}

func (t *Notification) Index(ctx *fiber.Ctx) error {
	user := GetUserLogin(ctx)
	if !user.IsFound() {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Get data from model
	data := new(model.Notification).GetNotifications(user.Id)
	// if err != nil {
	// 	return err
	// }
	return ctx.JSON(data)
}

func (t *Notification) ReadAll(ctx *fiber.Ctx) error {
	user := GetUserLogin(ctx)
	if !user.IsFound() {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Get data from model
	result := new(model.Notification).ReadNotificationsForUser(user.Id, 0)
	data := make(map[string]string)
	if result == false {
		data["error"] = "Read all notification fail!"
	} else {
		data["success"] = "Success!"
	}
	return ctx.JSON(data)
}
func (t *Notification) Read(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(err)
	}
	user := GetUserLogin(ctx)
	if !user.IsFound() {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Get data from model
	result := new(model.Notification).ReadNotificationsForUser(user.Id, idSearch)
	data := make(map[string]string)
	if result == false {
		data["error"] = "Read notification fail!"
	} else {
		data["success"] = "Success!"
	}
	return ctx.JSON(data)
}
