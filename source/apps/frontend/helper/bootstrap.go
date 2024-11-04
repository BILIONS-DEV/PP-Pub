package helper

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/lang"
	"source/apps/frontend/model"
	"source/apps/frontend/view"
	"source/pkg/utility"
	"strings"
	"time"
)

func Bootstrap(ctx *fiber.Ctx) error {

	isLogin, _ := utility.InStringArray(string(ctx.Request().URI().Path()), []string{
		"/user/login",
	})
	if isLogin {
		return ctx.Redirect(config.URILogin)
	}

	uri := ctx.Path()
	rootDomain := ctx.Protocol() + "://" + ctx.Hostname()
	currentURL := ctx.BaseURL() + ctx.OriginalURL()
	backURL, _ := url.QueryUnescape(ctx.Query("backurl"))
	if utility.ValidateString(backURL) == "" {
		backURL = rootDomain
	}
	UserModel := new(model.User)
	UserLogin := UserModel.UserLogin(ctx)
	userAdmin := UserModel.UserAdminLogin(ctx)
	// set cookie referer => save referer khi register
	if string(ctx.Request().Header.Referer()) != "" && strings.Contains(string(ctx.Request().Header.Referer()), ctx.Hostname()) == false && !UserLogin.IsFound() {
		cookieMaster := &fiber.Cookie{
			Name:    "referer",
			Value:   string(ctx.Request().Header.Referer()),
			Expires: time.Now().Add(999999 * time.Hour),
		}
		ctx.Cookie(cookieMaster)
	}

	if !utility.InArray(uri, []string{
		config.URILogin, config.URIRegister, config.URILogout,
		//config.URIForgotPassWord, config.URIResetPassWord,
		config.URILinkVideo, config.URIAtl, config.URIAtlQuick}, false) {
		// Nếu không phải các page Login / Register / Logout,... & chưa đăng nhập
		if !UserLogin.IsFound() {
			linkLogin := config.URILogin + "?backurl=" + url.QueryEscape(currentURL)
			return ctx.Redirect(linkLogin)
		}
	} else {
		if utility.InArray(uri, []string{config.URILogin, config.URIRegister}, true) {
			if UserLogin.IsFound() {
				return ctx.Redirect(backURL)
			}
		}
	}
	if uri == "/" || uri == "/user/login" {
		return ctx.Redirect(config.URIDashboards)
	}
	var serviceHost = "ads-txt.bilsyndication.com"
	if UserLogin.UserInfo.ServiceHostName != "" {
		serviceHost = UserLogin.UserInfo.ServiceHostName
	}
	if UserLogin.UserInfo.LogoWidth == 0 {
		UserLogin.UserInfo.LogoWidth = 100
	}
	ctx.Locals(assign.KEY, assign.Schema{
		Uri:             uri,
		RootDomain:      rootDomain,
		HostName:        ctx.Hostname(),
		BackURL:         backURL,
		CurrentURL:      currentURL,
		Version:         "dev",
		Title:           "Self-service advertising system - Valueimpression.com",
		Logo:            UserLogin.UserInfo.Logo,
		LogoWidth:       UserLogin.UserInfo.LogoWidth,
		Brand:           UserLogin.UserInfo.Brand,
		ServiceHostName: serviceHost,
		Theme:           "muze",
		TemplatePath:    uri,
		ThemeSetting:    view.Setting.ReleVersion,
		UserLogin:       UserLogin,
		UserAdmin:       userAdmin,
		SidebarSetup:    config.SidebarSetup,
		LANG:            lang.Translate,
		DeviceUA:        utility.GetDeviceFromUA(string(ctx.Context().UserAgent())),
	})
	//if UserLogin.Logo != "" && UserLogin.RootDomain != "" {
	config.TitlePrefix = UserLogin.UserInfo.Brand
	//}
	//fmt.Printf("%+v\n", UserLogin)
	ctx.Locals("UserLogin", UserLogin)
	ctx.Locals("UserAdmin", userAdmin)
	return ctx.Next()
}
