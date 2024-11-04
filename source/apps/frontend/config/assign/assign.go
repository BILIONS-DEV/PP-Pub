package assign

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/lang"
	"source/apps/frontend/model"
)

const (
	KEY = "AssignGlobal"
)

type Schema struct {
	Uri             string           `json:"uri"`
	RootDomain      string           `json:"root_domain"`
	HostName        string           `json:"host_name"`
	BackURL         string           `json:"back_url"`
	CurrentURL      string           `json:"current_url"`
	Version         string           `json:"version"`
	Title           string           `json:"title"`
	Logo            string           `json:"logo"`
	LogoWidth       int              `json:"logo_width"`
	Brand           string           `json:"brand"`
	ServiceHostName string           `json:"service_host_name`
	Theme           string           `json:"theme"`
	TemplatePath    string           `json:"template_path"`
	RedditPixel     bool             `json:"reddit_pixel"`
	ThemeSetting    interface{}      `json:"theme_setting"`
	UserLogin       model.UserRecord `json:"user_login"`
	UserAdmin       model.UserRecord `json:"user_admin"`
	SidebarSetup    config.SidebarSetupUri
	LANG            lang.Translation
	DeviceUA        string `json:"is_set_currency"`
}

func Get(ctx *fiber.Ctx) Schema {
	assign := ctx.Locals(KEY).(Schema)
	return assign
}
