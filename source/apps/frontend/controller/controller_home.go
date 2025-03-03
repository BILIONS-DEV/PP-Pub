package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/view"
	"time"
)

type Home struct{}

type AssignHome struct {
	assign.Schema
}

func (th *Home) Dashboard(ctx *fiber.Ctx) error {
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Dashboards")

	envFilename := os.Getenv("MODE") + "hai.env"
	_, err := os.Stat(envFilename)
	if err == nil {
		envs, erro := godotenv.Read(envFilename)
		if erro == nil {
			assigns.Version = envs["JS_VERSION"]
		} else {
			assigns.Version = time.Now().Format("20060102")
		}
	}
	return ctx.Render("home/dashboard", assigns, view.LAYOUTMain)
}

func (th *Home) Test(ctx *fiber.Ctx) error {
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Test")
	return ctx.Render("home/test", assigns, view.LAYOUTEmpty)
}

func (t *Home) TestNpm(ctx *fiber.Ctx) error {
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Test NPM")
	return ctx.Render("home/test", assigns, view.LAYOUTMain)
}
