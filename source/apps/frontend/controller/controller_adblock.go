package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/pkg/ajax"
	"time"
)

type AdBlock struct{}

type AssignAdBlockIndex struct {
	assign.Schema
	Params      payload.AdBlockFilterPayload
	Inventories []model.InventoryRecord
}

func (t *AdBlock) Analytics(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	params := payload.AdBlockFilterPayload{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignAdBlockIndex{Schema: assign.Get(ctx)}
	if params.StartDate == "" || params.EndDate == "" {
		params.StartDate = time.Now().AddDate(0, 0, -10).Format("2006-01-02")
		params.EndDate = time.Now().Format("2006-01-02")
	}
	assigns.Params = params
	assigns.Inventories = new(model.Inventory).GetByUser(userLogin.Id)
	assigns.Title = config.TitleWithPrefix("Adblock Analytics")
	assigns.LANG.Title = "Adblock Analytics"
	return ctx.Render("adblock/analytics", assigns, view.LAYOUTMain)
}

type ResponseAdBlockAnalytics struct {
	Status           string `json:"status,omitempty"`
	Message          string `json:"message,omitempty"`
	Params           payload.AdBlockFilterPayload
	AdblockAnalytics AdBlockAnalytic
}

type AdBlockAnalytic struct {
	Found []int `gorm:"column:found" json:"found"`
	//NotFound []int    `gorm:"column:notfound" json:"notfound"`
	Dates []string `gorm:"column:date" json:"date"`
}

func (t *AdBlock) AnalyticsFilter(ctx *fiber.Ctx) error {
	var response = ResponseAdBlockAnalytics{
		Status: ajax.SUCCESS,
	}
	userLogin := GetUserLogin(ctx)
	// Get payload post
	var inputs payload.AdBlockFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// s, _ := json.MarshalIndent(inputs, "", "\t")
	// fmt.Printf("%+v\n", string(s))
	if inputs.StartDate == "" || inputs.EndDate == "" {
		inputs.StartDate = time.Now().AddDate(0, 0, -10).Format("2006-01-02")
		inputs.EndDate = time.Now().Format("2006-01-02")
	}
	preriod := new(model.AdBlock).GetPreriod(inputs.StartDate, inputs.EndDate)
	// Get data from model
	adBlocks, err := new(model.AdBlock).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		response.Status = ajax.ERROR
		response.Message = err.Error()
		return ctx.JSON(response)
	}
	// Print Data to JSON
	var adblockAnalytic AdBlockAnalytic
	for _, value := range preriod {
		flag := false
		date, _ := time.Parse("2006-01-02", value)
		adblockAnalytic.Dates = append(adblockAnalytic.Dates, date.Format("02 Jan 2006"))
		for _, adBlock := range adBlocks {
			if value == adBlock.Date {
				flag = true
				adblockAnalytic.Found = append(adblockAnalytic.Found, adBlock.Found)
				//adblockAnalytic.NotFound = append(adblockAnalytic.NotFound, adBlock.NotFound)
			}
		}
		if !flag {
			adblockAnalytic.Found = append(adblockAnalytic.Found, 0)
			//adblockAnalytic.NotFound = append(adblockAnalytic.NotFound, 0)
		}
	}
	response.AdblockAnalytics = adblockAnalytic

	response.Params = inputs
	return ctx.JSON(response)
}

type AdBlockGenerator struct {
	assign.Schema
	Inventories []model.InventoryRecord
}

func (t *AdBlock) AdblockAlertGenerator(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	assigns := AdBlockGenerator{Schema: assign.Get(ctx)}
	assigns.Inventories = new(model.Inventory).GetByUser(userLogin.Id)
	assigns.Title = config.TitleWithPrefix("Adblock Alert Generator")
	// assigns.LANG.Title = "Adblock Alert Generator"
	return ctx.Render("adblock/generator", assigns, view.LAYOUTMain)
}
