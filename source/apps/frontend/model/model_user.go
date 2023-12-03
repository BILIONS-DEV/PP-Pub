package model

import (
	"encoding/json"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/mail"
	"source/apps/frontend/config"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/block"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/mailjet"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type User struct{}

type UserRecord struct {
	mysql.TableUser
}

func (UserRecord) TableName() string {
	return mysql.Tables.User
}

type RecoverPassword struct {
	RootDomain string
	UserName   string
	Email      string
	Uuid       string
}

const (
	//CookieLogin      = "mcflgi"
	CookieLogin      = "mctehj"
	CookieLoginAdmin = "mcfIgiAd"
	CookieMaster     = "mcflgim"
	CookieReferer    = "referer"

	URlQuantumdex        = "https://be.quantumdex.io/"
	URLReValueimpression = "https://re.valueimpression.com/"
)

func (t *User) GetByFiltersForAdmin(inputs *payload.UserFilterPayload) (response datatable.Response, err error) {
	var users []UserRecord
	var total int64
	err = mysql.Client.
		Scopes(
			t.SetFilterStatus(inputs),
		).
		Model(&users).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&users).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Translate.Errors.InventoryError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(users)
	return
}

type UserRecordDatatable struct {
	UserRecord
	RowId    string `json:"DT_RowId"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`
	Action   string `json:"action"`
}

func (t *User) MakeResponseDatatable(users []UserRecord) (records []UserRecordDatatable) {
	for _, user := range users {
		rec := UserRecordDatatable{
			UserRecord: user,
			RowId:      strconv.FormatInt(user.Id, 10),
			Email:      block.RenderToString("supply/index/block.email.gohtml", user),
			Status:     block.RenderToString("supply/index/block.status.gohtml", user),
			FullName:   block.RenderToString("supply/index/block.live.gohtml", user),
			Action:     block.RenderToString("supply/index/block.action.gohtml", user),
		}
		records = append(records, rec)
	}
	return
}

func (t *User) SetFilterStatus(inputs *payload.UserFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Status != nil {
			switch inputs.PostData.Status.(type) {
			case string, int:
				if inputs.PostData.Status != "" {
					return db.Where("status = ?", inputs.PostData.Status)
				}
			case []string, []interface{}:
				return db.Where("status IN ?", inputs.PostData.Status)
			}
		}
		return db
	}
}

func (t *User) setFilterSearch(inputs *payload.UserFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool
		// Search from form of datatable <- not use
		if inputs.Search != nil && inputs.Search.Value != "" {
			flag = true
		}
		// Search from form filter
		if inputs.PostData.QuerySearch != "" {
			flag = true
		}
		if !flag {
			return db
		}
		return db.Where("name  LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *User) setOrder(inputs *payload.UserFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db.Order("id desc")
	}
}

func (t *User) GetAll() (users []UserRecord) {
	mysql.Client.Find(&users)
	return
}

func (t *User) GetEmailById(userId int64) string {
	var user UserRecord
	mysql.Client.Select("email").Last(&user, userId)
	return user.Email
}

func (t *User) UserAdminLogin(ctx *fiber.Ctx) (user UserRecord) {
	loginToken := ctx.Cookies(CookieLoginAdmin)
	if utility.ValidateString(loginToken) == "" {
		return UserRecord{}
	}
	user = t.GetByLoginToken(loginToken)
	if !user.IsFound() {
		return UserRecord{}
	}
	if !user.IsActive() {
		return UserRecord{}
	}
	if user.IsAdmin() || user.IsSale() {
		return user
	}
	return
}

// UserLogin check cookie and get user login info
//
// param: ctx
// return:
func (t *User) UserLogin(ctx *fiber.Ctx) (user UserRecord) {
	loginToken := ctx.Cookies(CookieLogin)
	if utility.ValidateString(loginToken) == "" {
		return
	}
	user = t.GetByLoginToken(loginToken)
	if !user.IsFound() {
		return
	}
	if !user.IsActive() {
		return
	}
	if !user.IsAdmin() {
		return
	}
	return
}

func (t *User) GetReferer(ctx *fiber.Ctx) (referer string) {
	referer = ctx.Cookies(CookieReferer)
	return
}

// GetByLoginToken get user by login token
//
// param: loginToken
// return:
func (t *User) GetByLoginToken(loginToken string) (user UserRecord) {
	mysql.Client.Where(UserRecord{mysql.TableUser{LoginToken: loginToken}}).Last(&user)
	return
}

// Login login for controller
//
// param: postData
// return:
func (t *User) Login(postData *payload.Login, lang lang.Translation) (user UserRecord, errs []ajax.Error) {
	// Validate input
	validate := t.ValidateLogin(postData)
	if len(validate) > 0 {
		errs = append(errs, validate...)
		return
	}
	// Select in mysql
	err := mysql.Client.Select("id", "login_token", "permission").
		Where(UserRecord{mysql.TableUser{
			Email:    postData.Email,
			Password: user.MakePassword(postData.Password),
		}}).First(&user).Error
	if !user.IsFound() {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.UserError.Login.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	} else if user.IsSubPublisher() != true {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "Permission denied",
		})
	}

	return
}

// ValidateLogin validate input for login
//
// param: postData
// return:
func (t *User) ValidateLogin(postData *payload.Login) (ajaxErrors []ajax.Error) {
	if utility.ValidateString(postData.Email) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "email",
			Message: "email is required",
		})
	} else {
		err := checkmail.ValidateFormat(postData.Email)
		if err != nil {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "email",
				Message: err.Error(),
			})
		}
	}
	if utility.ValidateString(postData.Password) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "password",
			Message: "password is required",
		})
	}
	return
}

// Register register account for controller
//
// param: postData
// return:
func (t *User) Register(postData *payload.Register, referer string, lang lang.Translation) (user UserRecord, errs []ajax.Error) {
	// validate before create member
	validate := t.ValidateRegister(postData)
	if len(validate) > 0 {
		errs = append(errs, validate...)
		return
	}
	// check if user already in use
	err := mysql.Client.Select("id").Where("email = ?", postData.Email).First(&user).Error
	if user.IsFound() {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.UserError.Email.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	// Create new member
	// newMember, err := th.CreateMemberPending(postData)
	newMember, err := t.CreateMemberApproved(postData, referer)
	if err != nil {
		errs = append(errs, ajax.Error{
			Message: err.Error(),
		})
		return
	} else {
		var info = mailjet.InfoMail{
			To: mailjet.Email{
				Email: postData.Email,
				Name:  postData.FirstName + postData.LastName,
			},
			CC:          nil,
			BCC:         nil,
			Subject:     "Registration successful.Welcome" + postData.FirstName + postData.LastName,
			ContentText: "",
			ContentHtml: htmlblock.Render("user/block/register_success.gohtml", newMember).String(),
		}
		err = info.SendMail()
		// fmt.Println(err)
	}

	// set revenueShare default
	err = new(RevenueShare).CreateRevenueShareForUser(newMember, mysql.RevenueShareDefault)
	if err != nil {
		errs = append(errs, ajax.Error{
			Message: err.Error(),
		})
		return
	}
	return newMember, errs
}

// ValidateRegister validate input from form post
//
// param: postData
// return:
func (t *User) ValidateRegister(postData *payload.Register) (ajaxErrors []ajax.Error) {
	if utility.ValidateString(postData.FirstName) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "first_name",
			Message: "first name is required",
		})
	}
	if utility.ValidateString(postData.LastName) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "last_name",
			Message: "last name is required",
		})
	}
	if utility.ValidateString(postData.Email) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "email",
			Message: "email is required",
		})
	} else {
		err := checkmail.ValidateFormat(postData.Email)
		if err != nil {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "email",
				Message: err.Error(),
			})
		}
		flag := t.CheckEmail(postData.Email)
		if flag {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "email",
				Message: "Email already exists",
			})
		}
	}
	if utility.ValidateString(postData.Password) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "password",
			Message: "password is required",
		})
	} else {
		if len(postData.Password) < 8 {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "password",
				Message: "password must be longer than 8 characters",
			})
		}
	}
	if utility.ValidateString(postData.ConfirmPassWord) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "confirm_password",
			Message: "confirm password is required",
		})
	} else {
		if len(postData.ConfirmPassWord) < 8 {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "confirm_password",
				Message: "confirm password must be longer than 8 characters",
			})
		}
	}
	if utility.ValidateString(postData.Password) != utility.ValidateString(postData.ConfirmPassWord) {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "",
			Message: "Password and Confirm Password don't match",
		})
	}
	if postData.Agree != 1 {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "agree",
			Message: "please read carefully and accept our terms of use",
		})
	}
	return
}

// CreateMemberPending create member with status default = pending
//
// param: postData
// return:
func (t *User) CreateMemberPending(postData *payload.Register, referer string) (user UserRecord, err error) {
	user.makeMemberInfo(postData, referer)
	user.Status = mysql.StatusPending
	err = mysql.Client.Create(&user).Error
	return
}

// CreateMemberApproved create member with status default = approved
//
// param: postData
// return:
func (t *User) CreateMemberApproved(postData *payload.Register, referer string) (user UserRecord, err error) {
	user.makeMemberInfo(postData, referer)
	user.Status = mysql.StatusApproved
	err = mysql.Client.Create(&user).Error
	return
}

// make member default info
//
// param: postData
func (rec *UserRecord) makeMemberInfo(postData *payload.Register, referer string) {
	if referer != "" {
		rec.Referer = referer
	}
	rec.Permission = mysql.UserPermissionSubPublisher
	rec.Email = strings.TrimSpace(postData.Email)
	rec.FirstName = strings.TrimSpace(postData.FirstName)
	rec.LastName = strings.TrimSpace(postData.LastName)
	rec.Password = rec.MakePassword(strings.TrimSpace(postData.Password))
	rec.AccountType = mysql.TYPEAccountTypeUserFree
	rec.LoginToken = rec.MakeLoginToken()
	rec.PaymentNet = 30
	rec.PaymentTerm = 1
}

// SetLogin set cookie for login
// param: ctx
// param: remember
func (rec *UserRecord) SetLogin(ctx *fiber.Ctx, remember bool) {
	// Set cookie user
	cookie := &fiber.Cookie{
		Name:  CookieLogin,
		Value: rec.LoginToken,
	}
	if !remember {
		cookie.Expires = time.Now().Add(999999 * time.Hour)
	} else {
		cookie.Expires = time.Now().Add(999999 * time.Hour)
	}
	ctx.Cookie(cookie)

	// set cookie DomainMaster pubpower.io => login cho help.pubpower.io
	cookieMaster := &fiber.Cookie{
		Name:    CookieMaster,
		Value:   rec.LoginToken,
		Expires: time.Now().Add(999999 * time.Hour),
		Domain:  "pubpower.io",
		Path:    "/",
	}
	ctx.Cookie(cookieMaster)

	// remove cookie referer
	ctx.Cookie(&fiber.Cookie{
		Name:    CookieReferer,
		Expires: time.Now().Add(-(time.Hour * 15 * 24)),
	})

	return
}

func (rec *UserRecord) SetLoginAdmin(ctx *fiber.Ctx) {
	// Set cookie admin
	if rec.IsAdmin() || rec.IsSale() {
		cookieAdmin := &fiber.Cookie{
			Name:    CookieLoginAdmin,
			Value:   rec.LoginToken,
			Expires: time.Now().Add(999999 * time.Hour),
		}
		ctx.Cookie(cookieAdmin)
	}
	return
}

func (t *User) GetById(id int64) (user UserRecord) {
	mysql.Client.Where("id = ?", id).Find(&user)
	user.GetRls()
	return
}

func (t *User) UpdateBilling(inputs *payload.UpdateBilling, userLogin UserRecord, userAdmin UserRecord) (errs []ajax.Error) {
	lang := lang.Translate
	validate := new(UserBilling).validateBilling(inputs)
	if len(validate) > 0 {
		errs = append(errs, validate...)
		return
	}
	recordOld := new(UserBilling).GetByUserId(inputs.UserId)
	record := new(UserBilling).makeRowBilling(inputs)
	if recordOld.Id != 0 {
		record.Id = recordOld.Id
		if record.Method == "bank" {
			record.NewUpdate = "new"
		}
		err := new(UserBilling).Update(record)
		if err != nil {
			if !utility.IsWindow() {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: lang.Errors.UserError.UpdateBilling.ToString(),
				})
			} else {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: err.Error(),
				})
			}
			return
		}
		// Push History
		recordNew := new(UserBilling).GetByUserId(inputs.UserId)
		var creatorId int64
		if userAdmin.Id != 0 {
			creatorId = userAdmin.Id
		} else {
			creatorId = userLogin.Id
		}
		err = history.PushHistory(&history.User{
			CreatorId: creatorId,
			Detail:    history.DetailUserBillingFE,
			RecordOld: mysql.TableUser{UserBilling: recordOld.TableUserBilling},
			RecordNew: mysql.TableUser{UserBilling: recordNew.TableUserBilling},
		})
	} else {
		err := new(UserBilling).Create(record)
		if err != nil {
			if !utility.IsWindow() {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: lang.Errors.UserError.CreateBilling.ToString(),
				})
			} else {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: err.Error(),
				})
			}
			return
		}

		// Push History
		recordNew := new(UserBilling).GetByUserId(inputs.UserId)
		var creatorId int64
		if userAdmin.Id != 0 {
			creatorId = userAdmin.Id
		} else {
			creatorId = userLogin.Id
		}
		_ = history.PushHistory(&history.User{
			CreatorId: creatorId,
			Detail:    history.DetailUserBillingFE,
			RecordNew: mysql.TableUser{UserBilling: recordNew.TableUserBilling},
		})
	}
	return
}

func (t *User) UpdateAccount(inputs *payload.UpdateAccount, user UserRecord, userAdmin UserRecord) (recordNew UserRecord, errs []ajax.Error) {
	lang := lang.Translate
	validate := t.validateAccount(inputs, user)
	if len(validate) > 0 {
		errs = append(errs, validate...)
		return
	}
	recordOld := t.GetById(user.Id)
	recordNew = t.makeRowAccount(inputs, user.Id)
	if recordNew.Id != 0 {
		err := t.Update(recordNew)
		if err != nil {
			if !utility.IsWindow() {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: lang.Errors.UserError.UpdateAccount.ToString(),
				})
			} else {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: err.Error(),
				})
			}
			return
		}
	}
	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.User{
		CreatorId: creatorId,
		Detail:    history.DetailUserProfileFE,
		RecordOld: recordOld.TableUser,
		RecordNew: recordNew.TableUser,
	})
	return
}

func (t *User) validateAccount(row *payload.UpdateAccount, user UserRecord) (ajaxErrors []ajax.Error) {
	if utility.ValidateString(row.Email) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "email",
			Message: "Email is required",
		})
	} else {
		flag := t.valid(row.Email)
		if !flag {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "email",
				Message: "Invalid Email",
			})
		}
		if user.Email != row.Email {
			flag = t.CheckEmail(row.Email)
			if flag {
				ajaxErrors = append(ajaxErrors, ajax.Error{
					Id:      "email",
					Message: "Email already exists",
				})
			}
		}
	}
	err := checkmail.ValidateFormat(row.Email)
	if err != nil {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "email",
			Message: err.Error(),
		})
	}
	if utility.ValidateString(row.FirstName) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "first_name",
			Message: "First Name is required",
		})
	}
	if utility.ValidateString(row.LastName) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "last_name",
			Message: "Last Name is required",
		})
	}
	if utility.ValidateString(row.Address) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "address",
			Message: "Address is required",
		})
	}
	if utility.ValidateString(row.State) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "state",
			Message: "State is required",
		})
	}
	if utility.ValidateString(row.City) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "city",
			Message: "City is required",
		})
	}
	if utility.ValidateString(row.Country) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "country",
			Message: "Country is required",
		})
	}
	if utility.ValidateString(row.ZipCode) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "zip_code",
			Message: "ZipCode is required",
		})
	}
	if utility.ValidateString(row.PhoneNumber) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "phone_number",
			Message: "Phone Number is required",
		})
	}
	// if row.PassWord == "" {
	//	ajaxErrors = append(ajaxErrors, ajax.Error{
	//		Id:      "password",
	//		Message: "PassWord is required",
	//	})
	// }
	// if row.ConfirmPassWord == "" {
	//	ajaxErrors = append(ajaxErrors, ajax.Error{
	//		Id:      "confirm_password",
	//		Message: "ConfirmPassWord is required",
	//	})
	// }
	// if row.PassWord != row.ConfirmPassWord {
	//	ajaxErrors = append(ajaxErrors, ajax.Error{
	//		Id:      "match",
	//		Message: "Password and Confirm Password don't match",
	//	})
	// }
	return
}

func (t *User) makeRowAccount(row *payload.UpdateAccount, userId int64) (user UserRecord) {
	user = t.GetById(userId)
	user.FirstName = row.FirstName
	user.LastName = row.LastName
	user.Email = strings.TrimSpace(row.Email)
	user.FirstName = strings.TrimSpace(row.FirstName)
	user.LastName = strings.TrimSpace(row.LastName)
	user.City = strings.TrimSpace(row.City)
	user.State = strings.TrimSpace(row.State)
	user.Country = strings.TrimSpace(row.Country)
	user.ZipCode = strings.TrimSpace(row.ZipCode)
	user.Address = strings.TrimSpace(row.Address)
	user.PhoneNumber = strings.TrimSpace(row.PhoneNumber)
	// user.Password = user.MakePassword(strings.TrimSpace(row.PassWord))
	// user.LoginToken = user.MakeLoginToken()
	return
}

func (t *User) Create(row UserRecord) (err error) {
	err = mysql.Client.Create(&row).Error
	return
}

func (t *User) Update(row UserRecord) (err error) {
	err = mysql.Client.Save(&row).Error
	return
}

func (t *User) GetByEmail(email string) (row UserRecord, err error) {
	mysql.Client.Where("email = ?", email).Find(&row)
	return
}

func (t *User) valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (t *User) SendLinkToEmail(email, rootDomain string, lang lang.Translation) (ajaxErrors []ajax.Error) {
	if utility.ValidateString(email) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "email",
			Message: "Email is required",
		})
		return
	}
	account, err := t.GetByEmail(email)
	if err != nil {
		if !utility.IsWindow() {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "email",
				Message: lang.Errors.UserError.GetByEmail.ToString(),
			})
		} else {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "email",
				Message: err.Error(),
			})
		}
		return
	}
	// if !account.IsFound() {
	//	ajaxErrors = append(ajaxErrors, ajax.Error{
	//		Id:      "email",
	//		Message: "The email you entered does not belong to any account",
	//	})
	//	return
	// }
	record := RecoverPassword{
		RootDomain: rootDomain,
		UserName:   account.FirstName + " " + account.LastName,
		Email:      email,
		Uuid:       account.MakeUuid(account.Email, account.FirstName, account.LastName, time.Now().Format("2006/01/02 15:04:05")),
	}
	var info = mailjet.InfoMail{
		To: mailjet.Email{
			Email: email,
			Name:  "self-serve",
		},
		CC:          nil,
		BCC:         nil,
		Subject:     "Change password link",
		ContentText: "",
		ContentHtml: htmlblock.Render("user/block/forget_password.gohtml", record).String(),
	}
	err = info.SendMail()
	if err != nil {
		if !utility.IsWindow() {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: lang.Errors.UserError.SendMail.ToString(),
			})
		} else {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	// Format("2006/01/02 15:04:05")
	expiredTime := time.Now().Add(time.Minute * 5)
	userForget := UserForgetPasswordRecord{mysql.TableUserForgetPassword{
		UserId:      account.Id,
		Email:       email,
		Uuid:        record.Uuid,
		IsUsed:      2,
		CreatedDate: time.Now(),
		ExpiredTime: expiredTime,
	}}
	err = new(UserForgetPassword).Handle(&userForget)
	if err != nil {
		if !utility.IsWindow() {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: lang.Errors.UserError.SendMail.ToString(),
			})
		} else {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	return
}

func (t *User) HandleNewPass(inputs *payload.NewPassWord, uuid, email string, lang lang.Translation) (ajaxErrors []ajax.Error) {
	row := new(UserForgetPassword).GetRecord(uuid, email)
	user := t.GetById(row.UserId)
	err := t.validateNewPass(inputs)
	if len(err) > 0 {
		ajaxErrors = append(ajaxErrors, err...)
		return
	}
	log := new(UserForgetPassword).UpdateLinkUsed(row)
	if log != nil {
		if !utility.IsWindow() {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: lang.Errors.UserError.ChangePassword.ToString(),
			})
		} else {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: log.Error(),
			})
		}
		return
	}
	newRow := UserRecord{mysql.TableUser{
		Id:         user.Id,
		Password:   user.MakePassword(inputs.NewPassWord),
		LoginToken: user.MakeLoginToken(),
	}}
	bug := mysql.Client.Table("user").Updates(&newRow).Where("id = ?", newRow.Id).Error
	if bug != nil {
		if !utility.IsWindow() {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: lang.Errors.UserError.ChangePassword.ToString(),
			})
		} else {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: bug.Error(),
			})
		}
		return
	}
	return
}

func (t *User) HandleChangePassword(inputs *payload.NewPassWord, user UserRecord, lang lang.Translation) (ajaxErrors []ajax.Error) {
	err := t.validateChangePass(inputs, user)
	if len(err) > 0 {
		ajaxErrors = append(ajaxErrors, err...)
		return
	}
	newRow := UserRecord{mysql.TableUser{
		Id:         user.Id,
		Password:   user.MakePassword(inputs.NewPassWord),
		LoginToken: user.MakeLoginToken(),
	}}
	bug := mysql.Client.Table("user").Updates(&newRow).Where("id = ?", newRow.Id).Error
	if bug != nil {
		if !utility.IsWindow() {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: lang.Errors.UserError.ChangePassword.ToString(),
			})
		} else {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "",
				Message: bug.Error(),
			})
		}
		return
	}
	return
}

func (t *User) validateNewPass(input *payload.NewPassWord) (ajaxErrors []ajax.Error) {
	if utility.ValidateString(input.ConfirmPassWord) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "confirm_password",
			Message: "Confirm PassWord is required",
		})
	}

	if utility.ValidateString(input.ConfirmPassWord) != utility.ValidateString(input.NewPassWord) {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "new_password",
			Message: "New PassWord don't match",
		})
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "confirm_password",
			Message: "Confirm PassWord don't match",
		})
	}
	return
}

func (t *User) validateChangePass(input *payload.NewPassWord, user UserRecord) (ajaxErrors []ajax.Error) {
	if utility.ValidateString(input.OldPassWord) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "old_password",
			Message: "Old PassWord is required",
		})
	} else {
		passMd5 := user.MakePassword(input.OldPassWord)
		if user.Password != passMd5 {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "old_password",
				Message: "Old PassWord is incorrect",
			})
		}

	}
	if utility.ValidateString(input.NewPassWord) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "new_password",
			Message: "New PassWord is required",
		})
	} else {
		if len(input.NewPassWord) < 8 {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "new_password",
				Message: "New Password must be longer than 8 characters",
			})
		}
	}

	if utility.ValidateString(input.ConfirmPassWord) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "confirm_password",
			Message: "Confirm PassWord is required",
		})
	} else {
		if len(input.ConfirmPassWord) < 8 {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "confirm_password",
				Message: "Confirm Password must be longer than 8 characters",
			})
		}
	}

	if utility.ValidateString(input.ConfirmPassWord) != utility.ValidateString(input.NewPassWord) {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "",
			Message: "PassWord and ConfirmPassword don't match",
		})
		// ajaxErrors = append(ajaxErrors, ajax.Error{
		//	Id:      "confirm_password",
		//	Message: "Confirm PassWord don't match",
		// })
	}
	return
}

func (t *User) CheckEmail(email string) bool {
	var user UserRecord
	mysql.Client.Model(&UserRecord{}).Where("email = ?", email).Find(&user)
	if user.Id == 0 {
		return false
	}
	return true
}

func (t *User) UpdatePassWord(inputs *payload.NewPassWord, user UserRecord, userAdmin UserRecord, rootDomain string, lang lang.Translation) (record UserRecord, errs []ajax.Error) {
	record = user
	validate := t.validateChangePass(inputs, user)
	if len(validate) > 0 {
		errs = append(errs, validate...)
		return
	}
	record = UserRecord{mysql.TableUser{
		Id:         user.Id,
		Password:   user.MakePassword(inputs.NewPassWord),
		LoginToken: user.MakeLoginToken(),
	}}
	bug := mysql.Client.Model(&UserRecord{}).Where("id = ?", record.Id).Updates(&record).Error
	if bug != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.UserError.ChangePassword.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: bug.Error(),
			})
		}
		return
	}
	t.SendEmailUpdatePassWord(user.Email, rootDomain)

	// Push History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.User{
		CreatorId: creatorId,
		Detail:    history.DetailUserChangePassFE,
		RecordOld: record.TableUser,
		RecordNew: record.TableUser,
	})
	return
}

func (t *User) SendEmailUpdatePassWord(email, rootDomain string) {
	layoutTime := "3:04:05 PM, January 2, 2006 MST"
	timeChangePassword := time.Now().Format(layoutTime)
	var info = mailjet.InfoMail{
		To: mailjet.Email{
			Email: email,
			Name:  "self-serve",
		},
		CC:          nil,
		BCC:         nil,
		Subject:     "Warning",
		ContentText: "",
		ContentHtml: htmlblock.Render("user/block/update_password.gohtml", fiber.Map{
			"time":       timeChangePassword,
			"rootDomain": rootDomain,
		}).String(),
	}
	info.SendMail()
	return
}

func (t *User) UpdateUserAfterSyncExchange(user UserRecord, Data Result) (errs []ajax.Error) {
	record := UserRecord{mysql.TableUser{
		ApdUid:        Data.Uid,
		ApdPlacmentId: Data.Placment,
		ApdAdsTxt:     Data.AdsTxt,
	}}
	bug := mysql.Client.Model(&UserRecord{}).Where("id = ?", user.Id).Updates(&record).Error
	if bug != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: bug.Error(),
		})
	}
	return
}

func (t *User) AfterRegister(user UserRecord, password string) (errs []ajax.Error) {

	// Data, errs := t.SyncToBequantumdex(user, password)
	// if len(errs) > 0 {
	//	return
	// }
	//
	// // update table user (apd_uid, apd_placment_id, apd_ads_txt)
	// errs = t.UpdateUserAfterSyncExchange(user, Data)
	// if len(errs) > 0 {
	//	return
	// }

	// create bidder apacdex
	// bidderApacdex, errs := new(System).AutoCreateApacdex(user, Data.AdsTxt)
	// if len(errs) > 0 {
	// 	return
	// }

	// create line-item apacdex
	// errs = new(LineItem).AutoCreateLineApacdex(user.Id, bidderApacdex, strconv.FormatInt(Data.Placment, 10))
	// if len(errs) > 0 {
	// 	return
	// }

	// Create default profile
	errs = new(Identity).AutoCreateDefaultProfile(user)
	if len(errs) > 0 {
		return
	}

	// Create default floor
	_ = new(Floor).AutoCreateDefaultFloor(user)
	return
}

type Result struct {
	Status    string                 `json:"status"`
	Messenger string                 `json:"messenger"`
	Placment  int64                  `json:"placment"`
	Uid       int64                  `json:"uid"`
	AdsTxt    mysql.TYPEAdsTxtCustom `json:"ads_txt"`
}

func (t *User) SyncToBequantumdex(user UserRecord, password string) (Data Result, errs []ajax.Error) {
	// sync exchange => get ads.txt, placementID
	url := "https://be.quantumdex.io/api/system/SyncUserPubpower"
	if utility.IsWindow() {
		url = "http://be.quantumdex.local/api/system/SyncUserPubpower"
		// return
	}
	payload := strings.NewReader("token=3fd42d7b0a1caaa77d5eddadaaa4f8a8&email=" + user.Email + "&password=" + password)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &Data)
	if Data.Status == "error" {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: Data.Messenger,
		})
		return
	}

	return
}

func (t *User) CheckUserLogin(userLogin UserRecord, userAdmin UserRecord, uri string) (isAccept bool) {
	// Mặc định là true
	isAccept = true

	// return true luôn nếu đang được login thông qua user admin hoặc sale
	if userAdmin.Permission == mysql.UserPermissionSale || userAdmin.Permission == mysql.UserPermissionAdmin {
		return
	}

	if !userLogin.IsFound() {
		return false
	}
	switch userLogin.Permission {
	case mysql.UserPermissionManagedService:
		isAccept = t.checkPermissionManagedService(uri)
		return
	case mysql.UserPermissionPublisher:
		isAccept = t.checkPermissionPublisher(uri)
		return
	}
	return
}

func (t *User) checkPermissionManagedService(uri string) (isAccept bool) {
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

func (t *User) checkPermissionPublisher(uri string) (isAccept bool) {
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

type ResponseGetSellerSystem struct {
	Status       bool     `json:"status"`
	Message      string   `json:"message"`
	Errors       []string `json:"errors"`
	SellerID     int64    `json:"seller_id"`
	SellerName   string   `json:"seller_name"`
	SellerDomain string   `json:"seller_domain"`
}

func (t *User) GetSellerSystem(user UserRecord) (sellerID int64, errs []ajax.Error) {

	url := URLReValueimpression + "api/sellers/getSellerSystem"
	// if utility.IsWindow() {
	// url = "http://re.valueimpression.local/api/sellers/SyncSellerAllSystem"
	// return
	// }

	if user.SellerType == 0 || user.SellerName == "" || user.SellerDomain == "" {
		return
	}
	payload := strings.NewReader("type=" + strconv.FormatInt(user.SellerType, 10) + "&seller_domain=" + user.SellerDomain)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Password", "bli@123")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var Data ResponseGetSellerSystem
	json.Unmarshal(body, &Data)
	if !Data.Status {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: Data.Message,
		})
	}
	sellerID = Data.SellerID

	return
}
