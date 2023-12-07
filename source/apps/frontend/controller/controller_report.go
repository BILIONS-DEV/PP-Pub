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

type Report struct{}

type AssignReport struct {
	assign.Schema
}

func (t *Report) Index(ctx *fiber.Ctx) error {
	assigns := AssignReport{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Report")

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

	return ctx.Render("report/index", assigns, view.LAYOUTMain)
}

func (t *Report) Dimension(ctx *fiber.Ctx) error {
	assigns := AssignReport{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Reports By Dimension")

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
	return ctx.Render("report/dimension", assigns, view.LAYOUTMain)
}

func (t *Report) Saved(ctx *fiber.Ctx) error {
	assigns := AssignReport{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Reports Saved Queries")

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
	return ctx.Render("report/saved", assigns, view.LAYOUTMain)
}
