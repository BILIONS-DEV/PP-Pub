package controllerv2

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/internal/entity/dto"
	"time"
)

type user struct {
	*handler
}

func (h *handler) InitRoutesUser(app fiber.Router) {
	this := user{h}
	//app.Post("/login", this.Login)
	app.Get("/testlogin", this.TestLogin)

}

func (t *user) TestLogin(ctx *fiber.Ctx) (err error) {
	userLogin := t.getUserLogin(ctx)
	if userLogin.HaveAccountManager {
		return ctx.JSON(fiber.Map{
			"userLoginID":      userLogin.ID,
			"userLogin":        userLogin.Email,
			"accountManagerID": userLogin.AccountManager.ID,
			"accountManager":   userLogin.AccountManager.Email,
		})
	}
	return ctx.JSON(fiber.Map{
		"userLoginID": userLogin.ID,
		"userLogin":   userLogin.Email,
	})
}

var (
	loginExpired = time.Now().Add(time.Hour * 72)
)

func (t *user) Login(ctx *fiber.Ctx) (err error) {
	// get input từ post request
	var input dto.Login
	if err = ctx.BodyParser(&input); err != nil {
		return ctx.JSON(dto.Fail(err))
	}
	// validate input
	//if errs := input.Validate(); len(errs) > 0 {
	//	return ctx.JSON(dto.Fail(errs...))
	//}
	// run login từ usecase
	var jwtToken string
	//jwtToken, userLogin, err := t.useCases.User.Login(&userUC.LoginInput{
	//	Email:        input.Email,
	//	Password:     input.Password,
	//	JwtSecret:    config.JwtSecret,
	//	LoginExpired: loginExpired,
	//})
	if err != nil {
		return ctx.JSON(dto.Fail(err))
	}
	// set cookie
	ctx.Cookie(&fiber.Cookie{
		Name:    config.JwtCookieName,
		Value:   jwtToken,
		Expires: loginExpired,
	})

	// Set login token cho các page kiểu check cũ
	//userLogin.SetLoginFE(ctx, input.Remember)

	return ctx.JSON(dto.OK(fiber.Map{
		"token": jwtToken,
	}))
}

func (t *user) Index(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(fiber.Map{
		"action": "Index",
	})
}
