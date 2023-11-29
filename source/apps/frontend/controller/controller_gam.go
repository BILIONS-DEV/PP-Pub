package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/ggapi"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"strconv"
	"time"
)

type Gam struct{}

type AssignGam struct {
	assign.Schema
	Gam           model.GamRecord
	Networks      []model.GamNetworkRecord
	Params        payload.GamIndex
	FirstSetup    bool
	CheckSelected bool
	StepSelected  StepSelected
}

func init() {
	//new(Gam).GetNetworkFromGG()
}

func (t *Gam) GetNetworkFromGG() {
	refreshToken := "1//01veGbO4JQ_DCCgYIARAAGAESNwF-L9IrWc3bs9EFkSp_Co21lwSQRzkfKHJzZD4tdS1w9Yy1CmgEKoUF1XChcU5JpuMZ0uBIF88"
	out, err := ggapi.GetNetworks(refreshToken)
	fmt.Printf("\n out: %+v \n", string(out))
	if err != nil {
		log.Println(err)
		return
	}
	resp := ggapi.Response{}
	err = json.Unmarshal(out, &resp)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("\n resp: %+v \n", resp)
}

func (t *Gam) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIGam)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignGam{Schema: assign.Get(ctx)}
	params := payload.GamIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("GAM")
	return ctx.Render("gam/index", assigns, view.LAYOUTMain)
}

func (t *Gam) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIGam)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.GamFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.GamNetwork).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		return err
	}
	return ctx.JSON(dataTable)
}

func (t *Gam) Add(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIGamAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignGam{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Create GAM")
	return ctx.Render("gam/add", assigns, view.LAYOUTMain)
}

type StepSelected struct {
	Step1 bool
	Step2 bool
	Step3 bool
	Step4 bool
}

func (t *Gam) Edit(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIGamEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignGam{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Edit GAM")
	gamId := ctx.Query("id")
	if gamId == "" {
		return ctx.Status(500).SendString("id is required")
	}
	gamIdInt, _ := strconv.ParseInt(gamId, 10, 64)
	gam := new(model.Gam).GetByIdOfUser(gamIdInt, userLogin.Id)
	if gam.Id < 1 {
		return ctx.Status(500).SendString("gam not found")
	}

	StepSelected := StepSelected{}
	StepSelected.Step1 = true

	networksOfGam := new(model.GamNetwork).GetByGamId(gam.Id, userLogin.Id)
	assigns.Networks = networksOfGam
	assigns.Gam = gam
	assigns.FirstSetup = true
	assigns.CheckSelected = false
	for _, network := range assigns.Networks {
		if network.ApiAccess != 0 {
			assigns.FirstSetup = false
		}
		if network.Status == 2 {
			assigns.CheckSelected = true
			StepSelected.Step2 = true
		}
		if network.ApiAccess == 1 {
			StepSelected.Step3 = true
		}
		if network.PushLineItem == 1 {
			StepSelected.Step4 = true
		}

	}
	//if assigns.FirstSetup {
	//	assigns.EnableCheckApi = true
	//}
	assigns.StepSelected = StepSelected
	return ctx.Render("gam/edit", assigns, view.LAYOUTMain)
}

func (t *Gam) SelectNetwork(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIGamSelectNetwork)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	var responses ajax.Responses

	// Get Payloads
	var inputs payload.GamSelectGam
	if err := ctx.BodyParser(&inputs); err != nil {
		responses.Status = "error"
		responses.Message = err.Error()
		return ctx.JSON(responses)
	}

	// Get recordOld
	recordOld := new(model.GamNetwork).GetById(inputs.NetworkId, userLogin.Id)

	// Handle Post Data
	if err := new(model.GamNetwork).SelectByUser(inputs, userLogin, GetLang(ctx)); err != nil {
		responses.Status = "error"
		responses.Message = err.Error()
		return ctx.JSON(responses)
	}

	responses.Status = "success"
	responses.Message = "Successful selection"

	// Get recordNew
	recordNew := new(model.GamNetwork).GetById(inputs.NetworkId, userLogin.Id)
	// Push History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userLogin.Id
	}
	_ = history.PushHistory(&history.GAM{
		Detail:    history.DetailGAMSetupFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableGamNetwork,
		RecordNew: recordNew.TableGamNetwork,
	})
	return ctx.JSON(responses)
}

func (t *Gam) GetNetworks(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIGamGetNetworks)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	var inputs struct {
		GamId int64 `json:"gam_id"`
	}
	if err := ctx.BodyParser(&inputs); err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	records := new(model.GamNetwork).GetByGam(userLogin.Id, inputs.GamId)
	//fmt.Printf("\n records: %+v \n", records)
	return ctx.JSON(records)
}

type QueryPushLine struct {
	GamNetworkIds []int64 `query:"gam_network_ids" json:"gam_network_ids"`
}

func (t *Gam) PushLine(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIGamPushLine)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := QueryPushLine{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	err := new(model.GamNetwork).PushLineItem(inputs.GamNetworkIds, userLogin.Id, userAdmin, GetLang(ctx))
	if err != nil {
		response.Status = ajax.ERROR
		response.Errors = []ajax.Error{ajax.Error{
			Id:      "",
			Message: "Error",
		}}
	} else {
		response.Status = ajax.SUCCESS
	}

	return ctx.JSON(response)
}

// Your credentials should be obtained from the Google
// Developer Console (https://console.developers.google.com).
var conf = &oauth2.Config{
	//ClientID:     "697051439959-ola1rso2vsjf3cqhvmnoh3lnmkd5mn1t.apps.googleusercontent.com",
	//ClientSecret: "HcItQm1nUvIIWgp1BCbWI07d",
	ClientID:     "1033050029778-nv67vpjb627a68e7kur22sn1s0so2v2s.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "https://apps.valueimpression.com/gam/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/dfp", oidc.ScopeOpenID, "profile", "email",
	},
	Endpoint: google.Endpoint,
}

func (t *Gam) Connect(ctx *fiber.Ctx) error {
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return ctx.Redirect(url)
}

type GamInputCallback struct {
	State string `json:"state"`
	Code  string `json:"code"`
	Scope string `json:"scope"`
}

func (t *Gam) Callback(ctx *fiber.Ctx) error {

	//https://myaccount.google.com/u/0/permissions?pli=1
	if errorString := ctx.Query("error"); errorString == "access_denied" {
		return renderCallbackError(ctx, "Access denied")
	}

	code := ctx.Query("code")
	if code == "" {
		return renderCallbackError(ctx, "Access token invalid")
	}

	ctxContext := context.Background()
	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctxContext = context.WithValue(ctxContext, oauth2.HTTPClient, httpClient)
	token, err := conf.Exchange(ctxContext, code)
	if err != nil {
		//return renderCallbackError(ctx, err.Error())
		return renderCallbackError(ctx, `Token has expired`)
	}

	// Xử lý lấy email
	provider, err := oidc.NewProvider(ctxContext, "https://accounts.google.com")
	if err != nil {
		fmt.Println(err)
	}
	var verifier = provider.Verifier(&oidc.Config{ClientID: "1033050029778-nv67vpjb627a68e7kur22sn1s0so2v2s.apps.googleusercontent.com"})

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		fmt.Println(err)
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctxContext, rawIDToken)
	if err != nil {
		fmt.Println(err)
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Log GAM Token: %+v \n", token)
	fmt.Printf("Log GAM Mail: %+v \n", claims)
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	var resp ggapi.Response

	if token.RefreshToken != "" {
		// GET NETWORK DFP
		out, err := ggapi.GetNetworks(token.RefreshToken)
		if err != nil {
			return renderCallbackError(ctx, err.Error())
		}
		resp = ggapi.Response{}
		err = json.Unmarshal(out, &resp)
		if err != nil {
			return renderCallbackError(ctx, err.Error())
		}
		if !resp.Status {
			return renderCallbackError(ctx, resp.Message)
		}

		//gam, err := new(model.Gam).Add(userLogin.Id, token, resp.User, GetLang(ctx))
		//if err != nil {
		//	return renderCallbackError(ctx, err.Error())
		//}

		//if len(resp.Networks) > 0 {
		//	err = new(model.GamNetwork).Push(userLogin.Id, resp, token)
		//	if err != nil {
		//		return renderCallbackError(ctx, err.Error())
		//	}
		//}
		//
		////return ctx.Redirect(config.URIGamEdit + "?id=" + strconv.FormatInt(gam.Id, 10))
		//return ctx.Redirect(config.URIGam)

		//return ctx.JSON(resp)

	} else if token.AccessToken != "" {
		// GET NETWORK DFP
		out, err := ggapi.GetNetworks(token.AccessToken)
		if err != nil {
			return renderCallbackError(ctx, err.Error())
		}

		resp = ggapi.Response{}
		err = json.Unmarshal(out, &resp)
		if err != nil {
			return renderCallbackError(ctx, err.Error())
		}
		if !resp.Status {
			return renderCallbackError(ctx, resp.Message)
		}

		//if len(resp.Networks) > 0 {
		//	firstNetwork := new(model.GamNetwork).GetByNetworkId(resp.Networks[0].Id, userLogin.Id)
		//	if firstNetwork.Id > 0 {
		//		return ctx.Redirect(config.URIGamEdit + "?id=" + strconv.FormatInt(firstNetwork.GamId, 10))
		//	}
		//}
		//return renderCallbackError(ctx, `Sorry, we were unable to sync your GAM data. Please go to <a href="https://myaccount.google.com/security-checkup/3?hl=en" rel="nofollow" target="_blank">google Third-party access</a>, then delete access "<b>GAM Service - PubPower</b>" and try again`)

		//gam := new(model.Gam).GetByUser(userLogin.Id, resp.User.Email)
		//if gam.Id <= 0 {
		//	return renderCallbackError(ctx, `Your token has expired, please login again. If you still don't see this message, please <a href="https://myaccount.google.com/security-checkup/3" rel="nofollow">remove our app on google</a> and try to sign in again. `)
		//}
		//if len(resp.Networks) > 0 {
		//	firstNetwork := new(model.GamNetwork).GetByNetworkId(resp.Networks[0].Id, userLogin.Id, gam.Id)
		//	if firstNetwork.Id > 0 {
		//		return ctx.Redirect(config.URIGamEdit + "?id=" + strconv.FormatInt(firstNetwork.GamId, 10))
		//	} else {
		//		err = new(model.GamNetwork).Push(userLogin.Id, gam, resp.Networks)
		//		if err != nil {
		//			return renderCallbackError(ctx, err.Error())
		//		}
		//		return ctx.Redirect(config.URIGamEdit + "?id=" + strconv.FormatInt(gam.Id, 10))
		//	}
		//}
	}
	resp.User.Email = claims.Email
	if len(resp.Networks) > 0 {
		err = new(model.GamNetwork).Push(userLogin.Id, userAdmin, resp, token)
		if err != nil {
			return renderCallbackError(ctx, err.Error())
		}
	}
	return ctx.Redirect(config.URIGam)
}

type AssignGamCallback struct {
	assign.Schema
	Message string
}

func renderCallbackError(ctx *fiber.Ctx, message string) error {
	assigns := AssignGamCallback{
		Schema:  assign.Get(ctx),
		Message: message,
	}
	assigns.Title = config.TitleWithPrefix("Connect GAM Error")
	return ctx.Render("gam/callback", assigns, view.LAYOUTMain)
}

func (t *Gam) GetToken(ctx *fiber.Ctx) error {
	gamID := ctx.Query("gamId")
	if gamID == "" {
		return ctx.JSON(fiber.Map{
			"error": "gamId not have",
		})
	}
	gam := mysql.TableGam{}
	mysql.Client.Table("gam").Last(&gam, gamID)
	if gam.Id == 0 {
		return ctx.JSON(fiber.Map{
			"error": "can not find gam in mysql",
		})
	}
	token := &oauth2.Token{RefreshToken: gam.RefreshToken}
	tokenSource := conf.TokenSource(oauth2.NoContext, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(newToken)
}

func (t *Gam) CheckApiAccess(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIGamCheckApiAccess)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	var inputs payload.GamCheckApiAcess
	err := json.Unmarshal(ctx.Body(), &inputs)

	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	type Response struct {
		Status   string                   `json:"status"`
		Message  string                   `json:"message"`
		Networks []model.GamNetworkRecord `json:"networks"`
	}
	var networks []model.GamNetworkRecord
	gam := new(model.Gam).GetById(inputs.GamId)

	for _, id := range inputs.ListNetwork {
		recNetwork := new(model.GamNetwork).GetById(id, userLogin.Id)
		if recNetwork.Status != mysql.StatusGamSelected {
			continue
		}
		// Lưu lại recordOld
		recordOld := recNetwork
		isEnable, err := ggapi.CheckAccessApi(gam.RefreshToken, recNetwork.NetworkId, recNetwork.NetworkName)
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}
		if !isEnable { // Nếu enable = false
			if err = mysql.Client.Model(&recNetwork).Update("api_access", mysql.ApiAccessDisable).Error; err != nil {
				return ctx.JSON(fiber.Map{
					"status":  "error",
					"message": err.Error(),
				})
			}
		} else { // Nếu enable = true
			//Nếu không phải lỗi do api access vẫn tính là api được bật
			if err = mysql.Client.Model(&recNetwork).Update("api_access", mysql.ApiAccessEnable).Error; err != nil {
				return ctx.JSON(fiber.Map{
					"status":  "error",
					"message": err.Error(),
				})
			}
		}
		networks = append(networks, recNetwork)

		// Get recordNew
		recordNew := new(model.GamNetwork).GetById(recNetwork.Id, userLogin.Id)
		// Push History
		_ = history.PushHistory(&history.GAM{
			Detail:    history.DetailGAMSetupFE,
			CreatorId: userLogin.Id,
			RecordOld: recordOld.TableGamNetwork,
			RecordNew: recordNew.TableGamNetwork,
		})
	}
	return ctx.JSON(fiber.Map{
		"status":   "success",
		"networks": networks,
	})
}
