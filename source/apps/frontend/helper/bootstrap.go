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
	hostName := ctx.Hostname()
	is, _ := utility.InStringArray(hostName, []string{
		"self-serve.interdogmedia.com",
	})
	if is {
		ctx.Context().Response.Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		return ctx.SendString(`Go to page: <a href="https://apps.valueimpression.com/">https://apps.valueimpression.com</a>`)
	}
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
		config.URIForgotPassWord, config.URIResetPassWord,
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

	ctx.Locals(assign.KEY, assign.Schema{
		Uri:          uri,
		RootDomain:   rootDomain,
		HostName:     ctx.Hostname(),
		BackURL:      backURL,
		CurrentURL:   currentURL,
		Version:      "dev",
		Title:        "Self-service advertising system - Valueimpression.com",
		Theme:        "muze",
		TemplatePath: uri,
		ThemeSetting: view.Setting.ReleVersion,
		UserLogin:    UserLogin,
		UserAdmin:    userAdmin,
		SidebarSetup: config.SidebarSetup,
		LANG:         lang.Translate,
		DeviceUA:     utility.GetDeviceFromUA(string(ctx.Context().UserAgent())),
	})
	ctx.Locals("UserLogin", UserLogin)
	ctx.Locals("UserAdmin", userAdmin)
	return ctx.Next()
}
