package controller

import (
	// "encoding/json"
	"github.com/gofiber/fiber/v2"
	// "io"
	// "net/http"
	"source/apps/frontend/model"
	"source/core/technology/mysql"
	// "source/core/technology/mysql"
	// "source/apps/frontend/payload"
	// "source/core/technology/mysql"
	"source/pkg/ajax"
	// "source/pkg/utility"
	// "strconv"
	// "fmt"
	"source/apps/frontend/view"
	"strings"
	// "time"
)

type Api struct {

}
type AssignAccountManagerInfo struct {
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

func (t *Api) GetInfoAccountManager(ctx *fiber.Ctx) error {
	response := ajax.Responses{}
	//locVN, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	requestHeaders := ctx.GetReqHeaders()
	token, exist := requestHeaders["Token"]
	if !exist || token == "" {
		return ctx.JSON(map[string]string{
			"error": "Token is required",
		})
		// token = "e0c5303bafd0a6db0d872be15b1d6d79"
	}

	publisher := new(model.User).GetByLoginToken(token)
	if publisher.Id == 0 {
		response.Status = ajax.ERROR
		response.Message = "Page not found"
		return ctx.JSON(response)
	}

	// var Manager mysql.TableManagerSub
	presenter := new(model.User).GetById(publisher.Presenter)
	if (presenter.Id == 0) {
		return ctx.SendString("")
	}

	manager, err := new(model.ManagerSub).GetByUser(publisher)
	if err != nil {
		return ctx.SendString(err.Error())
	}
	if (manager.Id == 0 || manager.Status != mysql.TYPEStatusOn) {
		return ctx.SendString("")
	}

	assigns := manager

	if assigns.Telegram != "" {
		if !strings.Contains(assigns.Telegram, "http") {
			if !strings.Contains(assigns.Telegram, "t.me/") {
				assigns.Telegram = "https://t.me/" + assigns.Telegram
			} else {
				assigns.Telegram = "https://" + assigns.Telegram
			}
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
	return ctx.Render("user/ajax/info-manager", assigns, view.LAYOUTEmpty)
}