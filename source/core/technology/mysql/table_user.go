package mysql

import (
	"github.com/syyongx/php2go"
	"gorm.io/gorm"
	"source/pkg/utility"
	"strings"
	"time"
)

type TableUser struct {
	Id                  int64               `gorm:"column:id" json:"id"`
	Email               string              `gorm:"column:email" json:"email"`
	Presenter           int64               `gorm:"column:presenter" json:"presenter"`
	Permission          TYPEUserPermission  `gorm:"column:permission" json:"permission"`
	Status              TYPEStatus          `gorm:"column:status" json:"status"`
	Password            string              `gorm:"column:password" json:"password"`
	LoginToken          string              `gorm:"column:login_token" json:"login_token"`
	AdsTxtCustomByAdmin TYPEAdsTxtCustom    `gorm:"column:ads_txt_custom_by_admin" json:"ads_txt_custom_by_admin"`
	SellerID            int64               `gorm:"column:seller_id" json:"seller_id"`
	SellerType          int64               `gorm:"column:seller_type" json:"seller_type"`
	SellerName          string              `gorm:"column:seller_name" json:"seller_name"`
	SellerDomain        string              `gorm:"column:seller_domain" json:"seller_domain"`
	DisabledSeller      string              `gorm:"column:disabled_seller" json:"disabled_seller"`
	FirstName           string              `gorm:"column:first_name" json:"first_name"`
	LastName            string              `gorm:"column:last_name" json:"last_name"`
	Address             string              `gorm:"column:address" json:"address"`
	City                string              `gorm:"column:city" json:"city"`
	State               string              `gorm:"column:state" json:"state"`
	ZipCode             string              `gorm:"column:zip_code" json:"zip_code"`
	Country             string              `gorm:"column:country" json:"country"`
	PhoneNumber         string              `gorm:"column:phone_number" json:"phone_number"`
	PaymentTerm         int                 `gorm:"column:payment_term" json:"payment_term"`
	PaymentNet          int                 `gorm:"column:payment_net" json:"payment_net"`
	AccountType         TYPEAccountTypeUser `gorm:"column:account_type" json:"account_type"`
	Revenue             float64             `gorm:"column:revenue" json:"revenue"`
	Referer             string              `gorm:"column:referer" json:"referer"`
	ApdUid              int64               `gorm:"column:apd_uid" json:"apd_uid"`
	ApdPlacmentId       int64               `gorm:"column:apd_placment_id" json:"apd_placment_id"`
	ApdAdsTxt           TYPEAdsTxtCustom    `gorm:"column:apd_ads_txt" json:"apd_ads_txt"`
	Demo                string              `gorm:"column:demo" json:"demo"`
	Guarantee           int                 `gorm:"column:guarantee" json:"guarantee"`
	GuaranteeCeiling    string              `gorm:"column:guarantee_ceiling" json:"guarantee_ceiling"`
	GuaranteeFrom       string              `gorm:"column:guarantee_from" json:"guarantee_from"`
	GuaranteePeriods    string              `gorm:"column:guarantee_periods" json:"guarantee_periods"`
	PaymentNote         string              `gorm:"column:payment_note" json:"payment_note"`
	Logged              int64               `gorm:"column:logged" json:"logged"`
	SystemSync          int64               `gorm:"column:system_sync" json:"system_sync"`
	SyncPocPoc          string              `gorm:"default:pending;column:sync_pocpoc" json:"sync_pocpoc"`
	ParentSub           string              `gorm:"column:parent_sub" json:"parent_sub"`
	Logo                string              `gorm:"column:logo" json:"logo"`
	RootDomain          string              `gorm:"column:root_domain" json:"root_domain"`
	CreatedBy           int64               `gorm:"column:created_by" json:"created_by"`
	UpdatedBy           int64               `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt           time.Time           `gorm:"column:created_at" json:"created_at"`
	DeletedAt           gorm.DeletedAt      `gorm:"column:deleted_at" json:"deleted_at"`
	UserBilling         TableUserBilling    `gorm:"-"`
}

func (TableUser) TableName() string {
	return Tables.User
}

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

type TYPEUserPermission int64

const (
	UserPermissionMember TYPEUserPermission = iota + 1
	UserPermissionAdmin
	UserPermissionNetwork
	UserPermissionCreator
	UserPermissionSale
	UserPermissionManagedService
	UserPermissionPublisher
	UserPermissionAffiliate
	UserPermissionSubPublisher
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

	case UserPermissionAffiliate:
		return "affiliate"

	case UserPermissionSubPublisher:
		return "sub-publisher"

	default:
		return ""
	}
}

// IsFound Check User Exists
//
// return:
func (rec TableUser) IsFound() (flag bool) {
	if rec.Id > 0 {
		flag = true
	}
	return
}
func (rec TableUser) IsEditPayment() (flag bool) {
	if rec.Id == 156 || rec.Id == 7 || rec.Id == 38 || rec.Id == 90 {
		flag = true
	}
	return
}

// IsAdmin Check permission admin for user
//
// return:
func (rec TableUser) IsCreator() (flag bool) {
	if rec.Permission == UserPermissionCreator {
		flag = true
	}
	return
}

// IsAdmin Check permission admin for user
//
// return:
func (rec TableUser) IsSubPublisher() (flag bool) {
	if rec.Permission == UserPermissionSubPublisher {
		flag = true
	}
	return
}

// IsAdmin Check permission admin for user
//
// return:
func (rec TableUser) IsAdmin() (flag bool) {
	if rec.Permission == UserPermissionAdmin {
		flag = true
	}
	return
}

// IsAdmin Check permission sale for user
//
// return:
func (rec TableUser) IsSale() (flag bool) {
	if rec.Permission == UserPermissionSale {
		flag = true
	}
	return
}

// IsMember Check permission member for user
//
// return:
func (rec TableUser) IsMember() (flag bool) {
	if rec.Permission == UserPermissionMember {
		flag = true
	}
	return
}

// IsActive Check status approved of user
//
// return:
func (rec TableUser) IsActive() (flag bool) {
	if rec.Status == StatusApproved {
		flag = true
	}
	return
}

// IsApproved Check status approved of user
//
// return:
func (rec TableUser) IsApproved() (flag bool) {
	return rec.IsActive()
}

// MakePassword render password with hash from input string
//
// param: password
// return:
func (rec TableUser) MakePassword(password string) string {
	key := "as3df!"
	return utility.GetMD5Hash(utility.GetMD5Hash(password + key))
}

func (rec TableUser) MakePasswordVLI(password string) string {
	return php2go.Md5(password + php2go.Sha1(php2go.Md5(password)))
}

// MakeLoginToken render login token for user
//
// return:
func (rec TableUser) MakeLoginToken() string {
	if utility.ValidateString(rec.Email) == "" {
		return "email empty"
	}
	if utility.ValidateString(rec.Password) == "" {
		return "password empty"
	}
	key := "as3df!2312@"
	return utility.GetMD5Hash(utility.GetMD5Hash(rec.Email + rec.Password + key))
}

// MakeUuid render uuid for recovery password for user
//
// return:
func (rec TableUser) MakeUuid(email, firstname, lastname, time string) string {
	key := "$5$^&*#!"
	return utility.GetMD5Hash(utility.GetMD5Hash(email + key + firstname + lastname + time))
}

func (p TYPEUserPermission) ToPermission() string {
	switch p {
	case UserPermissionMember:
		return "MEMBER"
	case UserPermissionAdmin:
		return "ADMIN"
	default:
		return ""
	}
}

func (rec TableUser) SplitNameUser() string {
	if len(rec.FirstName) > 0 {
		newArr := strings.Split(rec.FirstName, "")
		return newArr[0]
	}
	return ""
}
func (rec *TableUser) GetRls() {
	if rec.Id == 0 {
		return
	}
	var userBilling TableUserBilling
	Client.Where("id = ?", rec.Id).Find(&userBilling)
	rec.UserBilling = userBilling
}
