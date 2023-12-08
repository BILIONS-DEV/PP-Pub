package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func User(app *fiber.App) {
	User := new(controller.User)
	app.Get(config.URILogin, User.Login)
	app.Post(config.URILogin, User.LoginPost)
	app.Get(config.URILogout, User.Logout)
	app.Get(config.URIRegister, User.Register)
	app.Post(config.URIRegister, User.RegisterPost)
	app.Get(config.URIBillingSetting, User.BillingSettingGet)
	app.Post(config.URIBillingSetting, User.BillingSettingPost)
	app.Get(config.URIAccountSetting, User.AccountSettingGet)
	app.Post(config.URIAccountSetting, User.AccountSettingPost)
	app.Get(config.URIForgotPassWord, User.ForgotPassWordGet)
	app.Post(config.URIForgotPassWord, User.ForgotPassWordPost)
	app.Get(config.URIResetPassWord, User.NewPassWord)
	app.Post(config.URIResetPassWord, User.NewPassWordPost)
	//app.Get(config.URIChangePassWord, User.ChangePasswordGet)
	// app.Post(config.URIChangePassWord, User.ChangePasswordPost)
	app.Post(config.URIChangePassword, User.PassWordSettingPost)
	app.Get(config.URIAtl, User.AutoLogin)
	app.Get(config.URIAtlQuick, User.QuickLogin)
}
