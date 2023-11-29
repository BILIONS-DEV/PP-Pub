package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/lang"
	"source/apps/frontend/model"
)

const (
	LayoutMain  = "pages/_components/layouts/main"
	LayoutLogin = "pages/_components/layouts/login"
	LayoutEmpty = "pages/_components/layouts/empty"
)

// GetUserLogin get user login for controller
//
// param: ctx
// return:
func GetUserLogin(ctx *fiber.Ctx) model.UserRecord {
	user := ctx.Locals("UserLogin").(model.UserRecord)
	return user
}

// GetUserAdmin get user admin for backend
//
// param: ctx
// return:
func GetUserAdmin(ctx *fiber.Ctx) model.UserRecord {
	user := ctx.Locals("UserAdmin").(model.UserRecord)
	return user
}

func GetLang(ctx *fiber.Ctx) (lang lang.Translation) {
	lang = assign.Get(ctx).LANG
	return
}
