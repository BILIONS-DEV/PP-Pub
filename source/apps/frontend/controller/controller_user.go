package controller

import (
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/pkg/ajax"
	"source/pkg/encryption"
	"strings"
	"time"
)

type User struct{}

type AssignUserBilling struct {
	assign.Schema
	Row model.UserBillingRecord
}

type AssignUserAccount struct {
	assign.Schema
	Row     model.UserRecord
	Billing model.UserBillingRecord
}

type AssignUserForgetPassword struct {
	assign.Schema
}

type AssignUserNewPass struct {
	assign.Schema
	Uuid  string
	Email string
}

// Login login page
//
// param: ctx
func (t *User) Login(ctx *fiber.Ctx) error {
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Login")
	assigns.Theme = "muze-login"
	assigns.Logo = ""

	// Lấy domain từ URL
	domain := ctx.Hostname()
	fmt.Println("domain: ", domain)
	if domain != "" {
		publisherAdmin := new(model.User).GetInFoPublisherAdminBySubDomain(domain)
		if publisherAdmin.Id > 0 {
			assigns.Logo = publisherAdmin.Logo
			assigns.Brand = publisherAdmin.Brand
			config.TitlePrefix = publisherAdmin.Brand
		}
	}

	if assigns.UserLogin.IsFound() {
		return ctx.Redirect(assigns.BackURL)
	}
	return ctx.Render("user/login", assigns, view.LAYOUTLogin)
}

func (t *User) LoginPost(ctx *fiber.Ctx) error {
	postData := new(payload.Login)
	if err := ctx.BodyParser(postData); err != nil {
		return err
	}
	response := ajax.Responses{}
	userLogin, errs := new(model.User).Login(postData, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		remember := false
		if postData.Remember == 1 {
			remember = true
		}
		userLogin.SetLogin(ctx, remember)
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

// Register register page
//
// param: ctx
func (t *User) Register(ctx *fiber.Ctx) error {

	return ctx.SendStatus(fiber.StatusNotFound)

	reddit := ctx.Query("utm_campain")
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Register")
	assigns.Theme = "muze-login"
	if reddit == "reddit" {
		assigns.RedditPixel = true
	}
	if assigns.UserLogin.IsFound() {
		return ctx.Redirect(assigns.BackURL)
	}
	return ctx.Render("user/register", assigns, view.LAYOUTLogin)
}

func (t *User) RegisterPost(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotFound)
	postData := new(payload.Register)
	if err := ctx.BodyParser(postData); err != nil {
		return err
	}
	response := ajax.Responses{}
	referer := new(model.User).GetReferer(ctx)
	newMember, errs := new(model.User).Register(postData, referer, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		newMember.SetLogin(ctx, false)
		response.Status = ajax.SUCCESS
		response.DataObject = newMember

		// _, errs = new(model.System).AutoCreateApacdex(newMember)
		errs = new(model.User).AfterRegister(newMember, postData.Password)
		if len(errs) > 0 {
			response.Status = ajax.ERROR
			response.Errors = errs
		}
	}

	return ctx.JSON(response)
}

func (t *User) BillingSettingGet(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBillingSetting)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignUserBilling{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Update Billing")
	assigns.Row = new(model.UserBilling).GetByUserId(userLogin.Id)
	return ctx.Render("user/billing", assigns, view.LAYOUTMain)
}

func (t *User) BillingSettingPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBillingSetting)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	postData := new(payload.UpdateBilling)
	if err := ctx.BodyParser(postData); err != nil {
		return err
	}
	response := ajax.Responses{}
	postData.UserId = userLogin.Id
	errs := new(model.User).UpdateBilling(postData, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		// response.DataObject = obj
	}
	return ctx.JSON(response)
}

func (t *User) AccountSettingGet(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAccountSetting)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignUserAccount{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Update Account")
	assigns.Row = new(model.User).GetById(userLogin.Id)
	assigns.Billing = new(model.UserBilling).GetByUserId(userLogin.Id)
	return ctx.Render("user/account", assigns, view.LAYOUTMain)
}

func (t *User) AccountSettingPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAccountSetting)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	postData := new(payload.UpdateAccount)
	if err := ctx.BodyParser(postData); err != nil {
		return err
	}
	response := ajax.Responses{}
	record, errs := new(model.User).UpdateAccount(postData, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		record.SetLogin(ctx, false)
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *User) PassWordSettingPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChangePassword)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignUserAccount{Schema: assign.Get(ctx)}
	postData := new(payload.NewPassWord)
	if err := ctx.BodyParser(postData); err != nil {
		return err
	}
	response := ajax.Responses{}
	record, errs := new(model.User).UpdatePassWord(postData, userLogin, userAdmin, assigns.RootDomain, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		record.SetLogin(ctx, false)
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *User) Logout(ctx *fiber.Ctx) error {
	ctx.Request().Header.VisitAllCookie(func(key, value []byte) {
		ctx.Cookie(&fiber.Cookie{
			Name: string(key),
			// Set expiry date to the past
			Expires: time.Now().Add(-(time.Hour * 365 * 24)),
			// HTTPOnly: true,
			// SameSite: "lax",
		})
		return
	})
	ctx.Cookie(&fiber.Cookie{
		Name: model.CookieLogin,
		// Set expiry date to the past
		Expires: time.Now().Add(-(time.Hour * 15 * 24)),
		// HTTPOnly: true,
		// SameSite: "lax",
	})
	ctx.Cookie(&fiber.Cookie{
		Name: model.CookieLoginAdmin,
		// Set expiry date to the past
		Expires: time.Now().Add(-(time.Hour * 15 * 24)),
		// HTTPOnly: true,
		// SameSite: "lax",
	})
	return ctx.Redirect(config.URILogin)
}

func (t *User) ForgotPassWordGet(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotFound)

	assigns := AssignUserAccount{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Forget PassWord")
	assigns.Logo = ""

	// Lấy đường dẫn URL chính của yêu cầu
	urlPath := ctx.OriginalURL()

	// Parse URL từ đường dẫn
	parsedURL, err := url.Parse(urlPath)
	if err != nil {
		// Xử lý lỗi nếu cần thiết
		return err
	}

	// Lấy domain từ URL
	domain := parsedURL.Hostname()
	if domain != "" {
		publisherAdmin := new(model.User).GetInFoPublisherAdminBySubDomain(domain)
		if publisherAdmin.Id > 0 {
			assigns.Logo = publisherAdmin.Logo
			assigns.Brand = publisherAdmin.Brand
			config.TitlePrefix = publisherAdmin.Brand
		}
	}

	return ctx.Render("user/forget-password", assigns, view.LAYOUTLogin)
}

func (t *User) ForgotPassWordPost(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotFound)

	assigns := AssignUserAccount{Schema: assign.Get(ctx)}
	postData := new(payload.UpdateAccount)
	if err := ctx.BodyParser(postData); err != nil {
		return err
	}
	response := ajax.Responses{}
	errs := new(model.User).SendLinkToEmail(postData.Email, assigns.RootDomain, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *User) NewPassWord(ctx *fiber.Ctx) error {
	assigns := AssignUserNewPass{Schema: assign.Get(ctx)}
	assigns.Uuid = ctx.Query("uuid")
	assigns.Email = ctx.Query("email")
	err := new(model.UserForgetPassword).GetByUuidEmail(assigns.Uuid, assigns.Email, GetLang(ctx))
	if err != nil {
		return ctx.SendString(err.Error())
	}
	assigns.Title = config.TitleWithPrefix("New PassWord")
	return ctx.Render("user/new-password", assigns, view.LAYOUTLogin)
}

func (t *User) NewPassWordPost(ctx *fiber.Ctx) error {
	postData := new(payload.NewPassWord)
	if err := ctx.BodyParser(postData); err != nil {
		return err
	}
	response := ajax.Responses{}
	errs := new(model.User).HandleNewPass(postData, postData.Uuid, postData.Email, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *User) ChangePasswordGet(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChangePassword)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignUserForgetPassword{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Change PassWord")
	return ctx.Render("user/password", assigns, view.LAYOUTMain)
}

func (t *User) ChangePasswordPost(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIChangePassword)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	postData := new(payload.NewPassWord)
	if err := ctx.BodyParser(postData); err != nil {
		return err
	}
	response := ajax.Responses{}
	errs := new(model.User).HandleChangePassword(postData, userLogin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		ctx.Cookie(&fiber.Cookie{
			Name: model.CookieLogin,
			// Set expiry date to the past
			Expires: time.Now().Add(-(time.Hour * 15 * 24)),
			// HTTPOnly: true,
			// SameSite: "lax",
		})
	}
	return ctx.JSON(response)
}

func (t *User) AutoLogin(ctx *fiber.Ctx) error {
	ctx.Request().Header.VisitAllCookie(func(key, value []byte) {
		ctx.Cookie(&fiber.Cookie{
			Name: string(key),
			// Set expiry date to the past
			Expires: time.Now().Add(-(time.Hour * 365 * 24)),
			// HTTPOnly: true,
			// SameSite: "lax",
		})
		return
	})
	token := ctx.Query("token")
	if token == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	ciphertextByte, err := hex.DecodeString(token)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	plaintext := encryption.Decrypt(ciphertextByte, "autologin")

	arrayString := strings.Split(string(plaintext), "|")
	if len(arrayString) == 0 || len(arrayString) < 3 {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	tie, err := time.Parse("2006-01-02 15:04:05", arrayString[1])

	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	flag := t.CheckLifeTime(tie)
	if !flag {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	user := new(model.User).GetByLoginToken(arrayString[0])
	if !user.IsFound() {
		return ctx.Redirect(config.URILogin, fiber.StatusFound)
	}
	user.SetLogin(ctx, true)
	userAdmin := new(model.User).GetByLoginToken(arrayString[2])
	userAdmin.SetLoginAdmin(ctx)
	// cookie := new(fiber.Cookie)
	// cookie.Name = model.CookieLogin
	// cookie.Value = arrayString[0]
	// cookie.Expires = time.Now().Add(999999 * time.Hour)
	// // Set cookie
	// ctx.Cookie(cookie)
	return ctx.Redirect(config.URIHome)
}

func (t *User) QuickLogin(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	tokenSale := ctx.Query("sa")
	if tokenSale == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	email := ctx.Query("email")
	if email == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	user := new(model.User).GetByLoginToken(token)
	if !user.IsFound() || user.Email != email {
		return ctx.Redirect(config.URILogin, fiber.StatusFound)
	}
	user.SetLogin(ctx, true)
	userAdmin := new(model.User).GetByLoginToken(tokenSale)
	userAdmin.SetLoginAdmin(ctx)
	// cookie := new(fiber.Cookie)
	// cookie.Name = model.CookieLogin
	// cookie.Value = arrayString[0]
	// cookie.Expires = time.Now().Add(999999 * time.Hour)
	// // Set cookie
	// ctx.Cookie(cookie)
	return ctx.Redirect(config.URIHome)
}

func (t User) CheckLifeTime(tie time.Time) bool {
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	nowTime, _ := time.Parse("2006-01-02 15:04:05", nowStr)
	count := nowTime.Sub(tie)
	if count.Seconds() > 10 {
		return false
	}
	return true
}
