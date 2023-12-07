package model

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"source/pkg/utility"
	"time"
)

func (User) TableName() string {
	return "user"
}

type User struct {
	ID                  int64               `gorm:"column:id" json:"id"`
	Email               string              `gorm:"column:email" json:"email"`
	Presenter           int64               `gorm:"column:presenter" json:"presenter"`
	Permission          TYPEUserPermission  `gorm:"column:permission" json:"permission"`
	Status              TYPEStatus          `gorm:"column:status" json:"status"`
	Password            string              `gorm:"column:password" json:"password"`
	LoginToken          string              `gorm:"column:login_token" json:"login_token"`
	AdsTxtCustomByAdmin TYPEAdsTxtCustom    `gorm:"column:ads_txt_custom_by_admin" json:"ads_txt_custom_by_admin"`
	SellerType          int64               `gorm:"column:seller_type" json:"seller_type"`
	SellerName          string              `gorm:"column:seller_name" json:"seller_name"`
	SellerDomain        string              `gorm:"column:seller_domain" json:"seller_domain"`
	FirstName           string              `gorm:"column:first_name" json:"first_name"`
	LastName            string              `gorm:"column:last_name" json:"last_name"`
	Address             string              `gorm:"column:address" json:"address"`
	City                string              `gorm:"column:city" json:"city"`
	State               string              `gorm:"column:state" json:"state"`
	ZipCode             string              `gorm:"column:zip_code" json:"zip_code"`
	Country             string              `gorm:"column:country" json:"country"`
	PhoneNumber         string              `gorm:"column:phone_number" json:"phone_number"`
	PaymentTerm         int                 `gorm:"column:payment_term" json:"payment_term"`
	PaymentNote         string              `gorm:"column:payment_note" json:"payment_note"`
	PaymentNet          int                 `gorm:"column:payment_net" json:"payment_net"`
	AccountType         TYPEAccountTypeUser `gorm:"column:account_type" json:"account_type"`
	Revenue             float64             `gorm:"column:revenue" json:"revenue"`
	Referer             string              `gorm:"column:referer" json:"referer"`
	Demo                string              `gorm:"column:demo" json:"demo"`
	Logged              int64               `gorm:"column:logged" json:"logged"`
	ParentSub           string              `gorm:"column:parent_sub" json:"parent_sub"`
	Logo                string              `gorm:"column:logo" json:"logo"`
	RootDomain          string              `gorm:"column:root_domain" json:"root_domain"`
	Brand               string              `gorm:"column:brand" json:"brand"`
	SystemSync          int64               `gorm:"column:system_sync" json:"system_sync"`
	CreatedBy           int64               `gorm:"column:created_by" json:"created_by"`
	UpdatedBy           int64               `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt           time.Time           `gorm:"column:created_at" json:"created_at"`
	DeletedAt           gorm.DeletedAt      `gorm:"column:deleted_at" json:"deleted_at"`
}

type UserLoginModel struct {
	User
	HaveAccountManager bool
	AccountManager     *User
}

const (
	//CookieLoginFE    = "mcflgi"
	CookieLoginFE    = "mctehj"
	CookieLoginBE    = "mcfIgi"
	CookieLoginAdmin = "mcfIgiAd"
	CookieMaster     = "mcflgim"
	CookieReferer    = "referer"
)

// TYPEUserPermission :
type TYPEUserPermission int64

const (
	UserPermissionMember TYPEUserPermission = iota + 1
	UserPermissionAdmin
	UserPermissionNetwork
	UserPermissionCreator
	UserPermissionSale
	UserPermissionManagedService
	UserPermissionPublisher
)

func (p TYPEUserPermission) String() string {
	switch p {
	case UserPermissionMember:
		return "self-serve"

	case UserPermissionAdmin:
		return "admin"

	case UserPermissionNetwork:
		return "network"

	case UserPermissionCreator:
		return "creator"

	case UserPermissionSale:
		return "sale"

	case UserPermissionManagedService:
		return "managed service"

	case UserPermissionPublisher:
		return "publisher"

	default:
		return ""
	}
}

// TYPEAccountTypeUser :
type TYPEAccountTypeUser int64

const (
	TYPEAccountTypeUserFree TYPEAccountTypeUser = iota + 1
	TYPEAccountTypeUserSubcription
)

func (p TYPEAccountTypeUser) String() string {
	switch p {
	case TYPEAccountTypeUserFree:
		return "free"

	case TYPEAccountTypeUserSubcription:
		return "subcription"
	default:
		return ""
	}
}

func (u *User) IsFound() bool {
	if u.ID > 0 {
		return true
	}
	return false
}

func (u *User) IsAdmin() bool {
	if u.Permission == UserPermissionAdmin {
		return true
	}
	return false
}

func (u *User) IsActive() bool {
	if u.Status == StatusApproved {
		return true
	}
	return false
}

func (u *User) SetLoginFE(ctx *fiber.Ctx, remember bool) {
	fmt.Println(*u)
	// Set cookie user
	cookie := &fiber.Cookie{
		Name:  CookieLoginFE,
		Value: u.LoginToken,
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
		Value:   u.LoginToken,
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

func (u *User) SetLoginAff(ctx *fiber.Ctx, remember bool) {
	// Set cookie user
	cookie := &fiber.Cookie{
		Name:  CookieLoginBE,
		Value: u.LoginToken,
	}
	if !remember {
		cookie.Expires = time.Now().Add(999999 * time.Hour)
	} else {
		cookie.Expires = time.Now().Add(999999 * time.Hour)
	}
	ctx.Cookie(cookie)
	return
}

func HashPassword(password string) (hashPassword string) {
	key := "as3df!"
	return utility.GetMD5Hash(utility.GetMD5Hash(password + key))
}

func (u *User) HashPassword() (password string) {
	return HashPassword(u.Password)
	// if u.Password == "" {
	//	return ""
	// }
	// key := "as3df!"
	// return utility.GetMD5Hash(utility.GetMD5Hash(u.Password + key))
}

func (u *User) RenderLoginToken() (loginToken string) {
	if utility.ValidateString(u.Email) == "" {
		return ""
	}
	if utility.ValidateString(u.Password) == "" {
		return ""
	}
	key := "as3df!2312@"
	return utility.GetMD5Hash(utility.GetMD5Hash(u.Email + u.Password + key))
}

func (u User) MakePassword(password string) string {
	key := "as3df!"
	return utility.GetMD5Hash(utility.GetMD5Hash(password + key))
}

// MakeLoginToken render login token for user
//
// return:
func (rec User) MakeLoginToken() string {
	if utility.ValidateString(rec.Email) == "" {
		return "email empty"
	}
	if utility.ValidateString(rec.Password) == "" {
		return "password empty"
	}
	key := "as3df!2312@"
	return utility.GetMD5Hash(utility.GetMD5Hash(rec.Email + rec.Password + key))
}
