package controller

import (
	// "encoding/json"
	"github.com/gofiber/fiber/v2"
	// "io"
	// "net/http"
	"source/apps/frontend/model"
	// "source/apps/frontend/payload"
	// "source/core/technology/mysql"
	"source/pkg/ajax"
	// "source/pkg/utility"
	// "strconv"
	"strings"
	"time"
	"source/apps/frontend/view"
)

type Api struct{

}

type AccountInfo struct {
	Email     string        `json:"email"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Agency    string        `json:"agency"`
	ZipCode   string        `json:"zipcode"`
	State     string        `json:"state"`
	Address   string        `json:"address"`
	City      string        `json:"city"`
	Country   string        `json:"country"`
	Phone     string        `json:"phone"`
	Company   string        `json:"company"`
	Domains   []DataDomains `json:"domains"`
}
type DataDomains struct {
	ID     int64  `json:"id"`
	Domain string `json:"domain"`
	Status string `json:"status"`
}

type AssignAccountManagerInfo struct {
	Agency     model.UserRecord
	Avatar     string
	Name       string
	Email      string
	Skype      string
	Telegram   string
	Linkedin   string
	Whatsapp   string
	Wechat     string
	AgencyTime string
}

func (t *Api) GetInfoAccount(ctx *fiber.Ctx) error {
	response := ajax.Responses{}
	//locVN, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	requestHeaders := ctx.GetReqHeaders()
	token, exist := requestHeaders["Token"]
	if !exist || token == "" {
		return ctx.JSON(map[string]string{
			"error": "Token is required",
		})
	}
	publisher := new(model.User).GetByLoginToken(token)
	if publisher.Id == 0 {
		response.Status = ajax.ERROR
		response.Message = "Page not found"
		return ctx.JSON(response)
	}

	agency := new(model.User).GetById(publisher.Id)
	// if err != nil {
	// 	response.Status = ajax.ERROR
	// 	response.Message = "Agency not found"
	// 	return ctx.JSON(response)
	// }
	data := AccountInfo{
		Email:     publisher.Email,
		FirstName: publisher.FirstName,
		LastName:  publisher.LastName,
		Agency:    agency.Email,
		ZipCode:   publisher.ZipCode,
		State:     publisher.State,
		Address:   publisher.Address,
		City:      publisher.City,
		Country:   publisher.Country,
		Phone:     publisher.PhoneNumber,
		Domains:   make([]DataDomains, 0),
	}
	// var domains []model.InventoryRecord
	// domains, err = new(model.Inventory).GetByUserId(publisher.Id)
	// for _, value := range domains {
	// 	if value.DeletedAt.Valid {
	// 		continue
	// 	}
	// 	data.Domains = append(data.Domains, DataDomains{
	// 		ID:     value.Id,
	// 		Domain: value.Domain,
	// 		Status: value.Status.String(),
	// 	})

	// }
	//data.Agency = agency
	response.Status = ajax.SUCCESS
	response.DataObject = data
	return ctx.JSON(response)
}

func (t *User) GetInfoAccountManager1(ctx *fiber.Ctx) error {
	locVN, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	requestHeaders := ctx.GetReqHeaders()
	token, exist := requestHeaders["Token"]
	if !exist || len(token) == 0 {
		return ctx.JSON(map[string]string{
			"error": "Token is required",
		})
	}
	userLogin := new(model.User).GetByLoginToken(token)
	assigns := AssignAccountManagerInfo{
		AgencyTime: time.Now().In(locVN).Format("15:04"),
	}
	agency := new(model.User).GetById(userLogin.Presenter)
	// if err != nil {
	// 	return ctx.SendStatus(fiber.StatusNotFound)
	// }
	if agency.Id == 0 {
		return ctx.JSON("")
	}
	q := ctx.Query("q")
	if q != "vli" && q != "pp" {
		return ctx.JSON("")
	}
	if q == "vli" {
		assigns.Avatar = agency.UserInfo.AvatarVLI
		assigns.Name = agency.UserInfo.NameVLI
		assigns.Email = agency.UserInfo.EmailVLI
		assigns.Telegram = agency.UserInfo.TelegramVLI
		assigns.Skype = agency.UserInfo.SkypeVLI
		assigns.Linkedin = agency.UserInfo.LinkedinVLI
		// assigns.Whatsapp = agency.UserInfo.Whatsapp
		// assigns.Wechat = agency.UserInfo.Wechat
	}
	if q == "pp" {
		assigns.Avatar = agency.UserInfo.Avatar
		assigns.Name = agency.UserInfo.Name
		assigns.Email = agency.UserInfo.Email
		assigns.Telegram = agency.UserInfo.Telegram
		assigns.Skype = agency.UserInfo.Skype
		assigns.Linkedin = agency.UserInfo.Linkedin
		// assigns.Whatsapp = agency.UserInfo.Whatsapp
		// assigns.Wechat = agency.UserInfo.Wechat
	}
	if assigns.Email == "" {
		return ctx.JSON("")
	}
	if assigns.Telegram != "" {
		if !strings.Contains(assigns.Telegram, "t.me/") {
			assigns.Telegram = "https://t.me/" + assigns.Telegram
		} else if !strings.Contains(assigns.Telegram, "http") {
			assigns.Telegram = "https://" + assigns.Telegram
		}
	}
	if assigns.Skype != "" {
		if strings.Contains(assigns.Skype, "live:") {
			assigns.Skype = strings.ReplaceAll(assigns.Skype, "live:", "")
		} else if strings.Contains(assigns.Skype, "skype:") {
			assigns.Skype = strings.ReplaceAll(assigns.Skype, "skype:", "")
		}
	}
	if assigns.Linkedin != "" {
		if !strings.Contains(assigns.Linkedin, "http") {
			assigns.Linkedin = "https://" + assigns.Linkedin
		}
	}
	assigns.Agency = agency
	return ctx.Render("user/ajax/info-agency", assigns, view.LAYOUTEmpty)
}
