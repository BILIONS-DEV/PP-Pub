package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/lang"
	"source/internal/repo"
	"source/pkg/ajax"
	"time"
)

type UsecaseUser interface {
	GetByID(ID int64, fields ...string) (record model.User)
	GetUserLogin(ID int64, accountManagerID ...int64) (record model.UserLoginModel)
	Login(input *LoginInput) (userLogin model.User, err error)
	GetUserByCookie(ctx *fiber.Ctx, key string) (isLogin bool, record model.User)
	Filter(inputs *dto.UserFilterPayload) (totalRecord int64, records []model.User, err error)
	GetReferer(ctx *fiber.Ctx) (referer string)
	Register(postData *dto.PayloadUser) (newUser model.User, errs []ajax.Error)
	GetByLoginToken(loginToken string) (user model.User)
	UpdateProfile(input *dto.PayloadProfile, userLogin model.User) (errs []ajax.Error)
	GetBillingByUser(userID int64) (billing model.UserBilling, err error)
	UpdateBilling(input *dto.PayloadBilling) (errs []ajax.Error)
	GetByEmail(email string) (record model.User)
}

type user struct {
	Repos *repo.Repositories
	Lang  *lang.Translation
}

func NewUserFeUsecase(repos *repo.Repositories, lang *lang.Translation) *user {
	return &user{Repos: repos, Lang: lang}
}

func (t *user) GetByID(ID int64, fields ...string) (record model.User) {
	return t.Repos.User.FindByID(ID, fields...)
}

func (t *user) GetUserLogin(ID int64, accountManagerID ...int64) (record model.UserLoginModel) {
	userLogin := t.Repos.User.FindByID(ID)
	record.User = userLogin
	if len(accountManagerID) > 0 && accountManagerID[0] > 0 {
		accountManager := t.Repos.User.FindByID(accountManagerID[0])
		record.AccountManager = &accountManager
		record.HaveAccountManager = true
	}
	return
}

type LoginInput struct {
	Email, Password, JwtSecret string
	LoginExpired               time.Time
}

func (t *user) Login(input *LoginInput) (rec model.User, err error) {
	// Get user từ db
	if rec, err = t.Repos.User.FindByLogin(input.Email, input.Password); err != nil {
		err = errors.New(`Invalid email or password`)
		return
	}
	// Kiểm tra xem user có tồn tại không
	if !rec.IsFound() {
		err = errors.New(`Invalid email or password`)
		return
	}
	return
}

func (t *user) buildJwtToken(rec *model.User, JwtSecret string, exp time.Time) (tokenString string, err error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"id":    rec.ID,
		"amid":  1,
		"roles": []string{rec.Permission.String()},
		"exp":   exp.Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	tokenString, err = token.SignedString([]byte(JwtSecret))
	if err != nil {
		return
	}
	return
}

func (t *user) GetUserByCookie(ctx *fiber.Ctx, key string) (isLogin bool, record model.User) {
	// Get loginToken từ cookie
	loginToken := ctx.Cookies(key)
	if loginToken == "" {
		return
	}
	// lấy user login từ loginToken trong cookie
	var err error
	if record, err = t.Repos.User.FindByLoginToken(loginToken); err != nil {
		return
	}
	// Set isLogin = true và trả ra dữ liệu
	if record.Presenter != 0 {
		presenter := t.GetByID(record.Presenter)
		if presenter.ParentSub == "yes" {
			record.Logo = presenter.Logo
			record.RootDomain = presenter.RootDomain
			record.Brand = presenter.Brand
		}
	}
	isLogin = true
	return isLogin, record
}

func (t *user) Filter(inputs *dto.UserFilterPayload) (totalRecord int64, records []model.User, err error) {
	return
}

func (t *user) GetReferer(ctx *fiber.Ctx) (referer string) {
	referer = ctx.Cookies(model.CookieReferer)
	return
}

type Register struct {
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassWord string `form:"confirm_password" json:"confirm_password"`
	FirstName       string `json:"first_name" form:"first_name"`
	LastName        string `json:"last_name" form:"last_name"`
	Agree           int    `json:"agree" form:"agree"`
}

func (t *user) Register(postData *dto.PayloadUser) (newUser model.User, errs []ajax.Error) {
	flag := t.Repos.User.CheckEmail(postData.Email)
	if flag {
		errs = append(errs, ajax.Error{
			Id:      "email",
			Message: "Email already exists",
		})
		return
	}
	// Create new member
	dataPost := postData.ToModel()
	err := t.Repos.User.Save(dataPost)
	if err != nil {
		errs = append(errs, ajax.Error{
			Message: err.Error(),
		})
		return
	}
	newUser = t.Repos.User.GetNewUser()
	// set revenueShare default
	// err = new(RevenueShare).CreateRevenueShareForUser(newMember, mysql.RevenueShareDefault)
	// if err != nil {
	// 	errs = append(errs, ajax.Error{
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }
	return
}

func (t *user) GetByLoginToken(loginToken string) (user model.User) {
	user = t.Repos.User.GetByLoginToken(loginToken)
	return
}

func (t *user) UpdateProfile(input *dto.PayloadProfile, userLogin model.User) (errs []ajax.Error) {
	if input.ID != userLogin.ID {
		errs = append(errs)
	}
	userOld := t.Repos.User.FindByID(input.ID)
	userOld.FirstName = input.FirstName
	userOld.LastName = input.LastName
	userOld.Address = input.Address
	userOld.City = input.City
	userOld.Country = input.Country
	userOld.State = input.State
	userOld.ZipCode = input.ZipCode
	userOld.PhoneNumber = input.PhoneNumber
	err := t.Repos.User.Save(userOld)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	return
}

func (t *user) GetBillingByUser(userID int64) (billing model.UserBilling, err error) {
	billing, err = t.Repos.User.GetBillingByUser(userID)
	return
}

func (t *user) UpdateBilling(input *dto.PayloadBilling) (errs []ajax.Error) {
	errs = input.Validate()
	if len(errs) > 0 {
		return
	}
	record := input.ToModel()
	err := t.Repos.User.SaveBilling(&record)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	return
}

func (t *user) GetByEmail(email string) (record model.User) {
	return t.Repos.User.GetByEmail(email)
}
