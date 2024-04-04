package controllerv2

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html"
	"html/template"
	"log"
	"net/url"
	"reflect"
	"source/apps/frontend/config"
	"source/apps/frontend/helper"
	"source/apps/frontend/helper/router"
	"source/apps/frontend/view"
	"source/infrastructure/theme"
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/usecase"
	"source/pkg/block"
	"source/pkg/utility"
)

type handler struct {
	useCases    *usecase.UseCases
	translation *lang.Translation
	block       *block.Block
}
type Deps struct {
	UseCases    *usecase.UseCases
	Translation *lang.Translation
	Block       *block.Block
}

// NewHandler : khởi tạo service server với các tiêm phụ thuộc là userCases & translation dùng cho multiple language
//
// param: useCases *usecase.UseCases
// param: translation *lang.Translation
// return: handler *handler
func NewHandler(deps Deps) *handler {
	return &handler{
		useCases:    deps.UseCases,
		translation: deps.Translation,
		block:       deps.Block,
	}
}

// StartServer : start server với package fiber.
// Đây sẽ là nơi khai báo & cấu hình fiber
//
// param: port
func (h *handler) StartServer(port int) {
	app := fiber.New(fiber.Config{
		BodyLimit:      20 * 1024 * 1024,
		ReadBufferSize: 128 * 1024,
		Views:          engine(),
	})
	app.Static("/assets", "../../www/themes/muze/assets", fiber.Static{MaxAge: 9999999})
	app.Static("/assets", "./view/.assets", fiber.Static{MaxAge: 9999999, CacheDuration: 9999999})
	app.Use(favicon.New()) //=> remove favicon.ico request
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	//=> setting router dành cho bản cũ, sau này không dùng nữa thì xóa đi
	app.Use(helper.Bootstrap)
	router.Register(app)

	//=> bản mới
	//Auth(app, authConfig{
	//	Secret:      config.JwtSecret,
	//	TokenLookup: config.JwtTokenLookup,
	//	ContextKey:  config.JwtContextKey,
	//	CookieName:  config.JwtCookieName,
	//}) //=> kiểm tra jwt + set userID vào fiberLocals + check permission
	app.Use(func(ctx *fiber.Ctx) error {
		return bootstrap(ctx, h)
	})
	h.setupRouters(app)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
func (h *handler) setupRouters(app fiber.Router) {
	handlerType := reflect.TypeOf(h)
	handlerValue := reflect.ValueOf(h)
	for i := 0; i < handlerType.NumMethod(); i++ {
		method := handlerType.Method(i)
		//=> tất cả các hàm initRoutes của controller đều phải bắt đầu bằng `InitRoutes`. Ví dụ: controller_user thì là InitRoutesUser
		methodPrefixRequired := string([]rune(method.Name)[0:10])
		if methodPrefixRequired == "InitRoutes" {
			method.Func.Call([]reflect.Value{handlerValue, reflect.ValueOf(app)})
		}
	}
}
func (h *handler) getUserLogin(ctx *fiber.Ctx) model.UserLoginModel {
	userLoginID, _ := ctx.Locals("userLoginID").(int64)
	accountManagerID, _ := ctx.Locals("accountManagerID").(int64)
	return h.useCases.User.GetUserLogin(userLoginID, accountManagerID)
}

/**
SUPPORT FUNCTION FOR CONTROLLER

Các hàm private hỗ trợ xử lý trong controller,
và khi khởi tạo các router
*/

const (
	assignKEY   = "AssignGlobal"
	titlePrefix = "PubPower.io Digital Media"
)

// Hàm preload trước khi đi vào từng controller method.
// Nghĩa là hàm này sẽ được load ra trước khi đi vào từng controller/method với toàn bộ các page
//
// param: ctx
func bootstrap(ctx *fiber.Ctx, h *handler) error {
	uri := ctx.Path()
	hostName := ctx.Hostname()
	rootDomain := ctx.Protocol() + "://" + ctx.Hostname()
	currentURL := ctx.BaseURL() + ctx.OriginalURL()
	backURL, _ := url.QueryUnescape(ctx.Query("backurl"))
	if utility.ValidateString(backURL) == "" {
		backURL = rootDomain
	}

	_, userLogin := h.useCases.User.GetUserByCookie(ctx, model.CookieLoginFE)
	ctx.Locals("UserLogin", userLogin)

	_, userAdmin := h.useCases.User.GetUserByCookie(ctx, model.CookieLoginAdmin)
	ctx.Locals("UserAdmin", userAdmin)

	//if !isLogin {
	//	if !utility.InArray(uri, []string{"/user/login", "/user/logout", "/user/atl"}, true) {
	//		linkLogin := "/user/login?backurl=" + url.QueryEscape(currentURL)
	//		return ctx.Redirect(linkLogin)
	//	}
	//} else {
	//	if utility.InArray(uri, []string{"/user/login"}, true) {
	//		return ctx.Redirect(backURL)
	//	}
	//}

	// check uri
	//if !checkPermissionURI(userLogin, userAdmin, ctx.OriginalURL()) {
	//	return ctx.Status(fiber.StatusNotFound).SendString("permission denied")
	//}

	parentSub := h.useCases.User.GetInfoByUserID(userLogin.Presenter)
	ctx.Locals(assignKEY, Assign{
		Uri:          uri,
		RootDomain:   rootDomain,
		HostName:     hostName,
		BackURL:      backURL,
		CurrentURL:   currentURL,
		Version:      "dev",
		Title:        "Self-service advertising system",
		Logo:         parentSub.Logo,
		Brand:        parentSub.Brand,
		Theme:        "muze",
		TemplatePath: uri,
		LANG:         h.translation,
		UserLogin:    userLogin,
		UserAdmin:    userAdmin,
		SidebarSetup: config.SidebarSetup,
		DeviceUA:     utility.GetDeviceFromUA(string(ctx.Context().UserAgent())),
	})
	return ctx.Next()
}

// getUserLoginGetUserLogin get user login for controller
//
// param: ctx
// return:
func getUserLogin(ctx *fiber.Ctx) model.User {
	user := ctx.Locals("UserLogin").(model.User)
	return user
}
func getUserAdmin(ctx *fiber.Ctx) model.User {
	user := ctx.Locals("UserAdmin").(model.User)
	return user
}
func checkPermissionURI(userLogin model.User, userAdmin model.User, uri string) (isAccept bool) {
	// Mặc định là true
	isAccept = true

	// return true luôn nếu đang được login thông qua user admin hoặc sale
	if userAdmin.Permission == model.UserPermissionSale || userAdmin.Permission == model.UserPermissionAdmin {
		return
	}

	switch userLogin.Permission {
	case model.UserPermissionManagedService:
		isAccept = checkPermissionManagedService(uri)
		return
	case model.UserPermissionPublisher:
		isAccept = checkPermissionPublisher(uri)
		return
	}
	return
}
func checkPermissionManagedService(uri string) (isAccept bool) {
	// Mặc định là true
	isAccept = true

	if utility.InArray(uri, config.SidebarSetup.VideoGroup, true) {
		return false
	}

	if utility.InArray(uri, config.SidebarSetup.SystemGroup, true) {
		return false
	}

	return
}
func checkPermissionPublisher(uri string) (isAccept bool) {
	// Mặc định là true
	isAccept = true

	if utility.InArray(uri, config.SidebarSetup.Demand, true) {
		return false
	}

	if utility.InArray(uri, config.SidebarSetup.Bidder, true) {
		return false
	}

	return
}

type Assign struct {
	Uri          string      `json:"uri"`
	RootDomain   string      `json:"root_domain"`
	HostName     string      `json:"host_name"`
	BackURL      string      `json:"back_url"`
	CurrentURL   string      `json:"current_url"`
	Version      string      `json:"version"`
	Title        string      `json:"title"`
	Logo         string      `json:"logo"`
	Brand        string      `json:"brand"`
	Theme        string      `json:"theme"`
	TemplatePath string      `json:"template_path"`
	ThemeSetting interface{} `json:"theme_setting"`
	UserLogin    model.User  `json:"user_login"`
	UserAdmin    model.User  `json:"user_admin"`
	SidebarSetup config.SidebarSetupUri
	LANG         *lang.Translation
	DeviceUA     string `json:"is_set_currency"`
}

func newAssign(ctx *fiber.Ctx, title string) Assign {
	assign := ctx.Locals(assignKEY).(Assign)
	assign.Title = titleWithPrefix(title)
	return assign
}
func titleWithPrefix(title string) (titleWithPrefix string) {
	titleWithPrefix = title
	return
}

/*
*
VIEW ENGINE
Khởi tạo view engine và truyền các hàm bổ trợ vào view của ứng dụng
*/
func engine() *html.Engine {
	engine := html.New("./view/pages", ".gohtml")
	engine.Layout("embed")    // Optional. Default: "embed"
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters
	if utility.IsWindow() {
		engine.Reload(true) // Optional. Default: false ... Bật cái này lên là bị race condition
		engine.Debug(false) // Optional. Default: false
	}
	addFunc(engine)
	return engine
}
func addFunc(engine *html.Engine) {
	engine.AddFunc("IsActiveSidebar", view.IsActiveSidebar)
	engine.AddFunc("IsActiveSidebarWithGroup", view.IsActiveSidebarWithGroup)

	engine.AddFunc("EmbedCssInline", theme.EmbedCssInline)
	engine.AddFunc("EmbedCSS", theme.EmbedCSS)
	engine.AddFunc("EmbedJS", theme.EmbedJS)
	engine.AddFunc("safeHTML", theme.SafeHTML)
	engine.AddFunc("FirstCharacter", utility.FirstCharacter)
	engine.AddFunc("FormatFloat", utility.FormatFloat)
	engine.AddFunc("StrToTime", theme.StrToTime)
	engine.AddFunc("IncIndex", theme.IncIndex)
	engine.AddFunc("FormatHtml", theme.FormatHtml)
	engine.AddFunc("Nl2br", theme.Nl2br)
	engine.AddFunc("InArray", utility.InArray)
	engine.AddFunc("MkArray", theme.MkSlice)
	engine.AddFunc("MkSlice", theme.MkSlice)
	engine.AddFunc("InSlice", func(search interface{}, array interface{}) bool {
		return utility.InArray(search, array, true)
	})
	// AddFunc adds a function to the template's global function map.
	engine.AddFunc("greet", func(name string) string {
		return "Hello, greetFunc: " + name + "!"
	})
	engine.AddFunc("IsStringEmpty", func(name string) bool {
		if name == "" {
			return false
		} else {
			return true
		}
	})
	engine.AddFunc("unescape", func(s string) template.HTML {
		return template.HTML(s)
	})
}
