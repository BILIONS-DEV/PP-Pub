package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func User(app *fiber.App) {
	user := new(controller.User)
	app.Get(config.URILogin, user.Login)
	app.Post(config.URILogin, user.LoginPost)
	app.Get(config.URILogout, user.Logout)
	app.Get(config.URIRegister, user.Register)
	app.Post(config.URIRegister, user.RegisterPost)
	//app.Get(config.URIBillingSetting, user.BillingSettingGet)
	app.Post(config.URIBillingSetting, user.BillingSettingPost)
	app.Get(config.URIAccountSetting, user.AccountSettingGet)
	app.Get(config.URIAccountBilling, user.AccountSettingGet)
	app.Get(config.URIAccountPassword, user.AccountSettingGet)
	app.Get(config.URIAccountProfile, user.AccountSettingGet)
	app.Post(config.URIAccountProfile, user.AccountSettingPost)
	app.Post(config.URIAccountSetting, user.AccountSettingPost)
	//app.Get(config.URIForgotPassWord, User.ForgotPassWordGet)
	//app.Post(config.URIForgotPassWord, User.ForgotPassWordPost)
	//app.Get(config.URIResetPassWord, User.NewPassWord)
	//app.Post(config.URIResetPassWord, User.NewPassWordPost)
	//app.Get(config.URIChangePassWord, User.ChangePasswordGet)
	// app.Post(config.URIChangePassWord, User.ChangePasswordPost)
	app.Post(config.URIChangePassword, user.PassWordSettingPost)
	app.Post(config.URIChangeTemplate, user.ChangeTemplatePost)
	app.Get(config.URIAtl, user.AutoLogin)
	app.Get(config.URIAtlQuick, user.QuickLogin)
}
