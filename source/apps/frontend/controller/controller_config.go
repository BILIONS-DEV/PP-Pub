package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/core/technology/mysql"
	"source/pkg/ajax"
)

type Config struct {
}

//Config
type AssignSystemConfigIndex struct {
	assign.Schema
	Config   model.ConfigRecord
	Currency []model.CurrencyRecord
}

func (t *Config) IndexConfig(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIConfig)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	recConfig := new(model.Config).GetByUserId(userLogin.Id)
	if recConfig.Id < 1 {
		recConfig = model.ConfigRecord{TableConfig: mysql.TableConfig{
			UserId:        userLogin.Id,
			PrebidTimeOut: 1000,
			AdRefreshTime: 30,
		}}
		err := mysql.Client.Create(&recConfig).Error
		if err != nil {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
	}
	assigns := AssignSystemConfigIndex{Schema: assign.Get(ctx)}
	assigns.Config = recConfig
	assigns.Currency = new(model.Currency).GetAll()
	assigns.Title = config.TitleWithPrefix("Config")
	return ctx.Render("config/index", assigns, view.LAYOUTMain)
}

func (t *Config) ConfigSave(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIConfigSave)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.SystemConfig{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}
	//fmt.Printf("%+v \n", inputs)
	//test
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}

	// Handle
	response := ajax.Responses{}
	newConfig, errs := new(model.System).SaveConfig(inputs, userLogin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = newConfig
	}

	return ctx.JSON(response)
}
