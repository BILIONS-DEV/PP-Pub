package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/htmlblock"
	"strconv"
)

type Bidder struct{}

type AssignSystemIndex struct {
	assign.Schema
	Params payload.SystemIndex
}

func (t *Bidder) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidder)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.SystemIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignSystemIndex{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Bidder")
	return ctx.Render("bidder/index", assigns, view.LAYOUTMain)
}

func (t *Bidder) Filter(ctx *fiber.Ctx) error {
	// Get data from model
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidder)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.SystemFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	dataTable, err := new(model.System).GetByFilters(&inputs, userLogin, GetLang(ctx))
	if err != nil {
		return err
	}
	return ctx.JSON(dataTable)
}

type AssignSystemAdd struct {
	assign.Schema
	Params                  map[string]string
	MediaType               []model.MediaTypeRecord
	GamNetworks             []model.GamNetworkRecord
	BidderTemplatesUnSelect []model.BidderTemplateRecord
	BidderTemplatesSelected []model.BidderTemplateRecord
}

func (t *Bidder) Add(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidderAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignSystemAdd{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Add Bidder")
	assigns.MediaType = new(model.MediaType).GetAll()
	bidderTemplates := new(model.BidderTemplate).GetAll()
	bidderSelected := new(model.Bidder).GetAllUser(userLogin.Id)
	mapBidderSelected := make(map[int64]int)
	for _, bidder := range bidderSelected {
		//Loại trừ bidder = 1 là google được phép add nhiều
		if bidder.BidderTemplateId != 1 {
			mapBidderSelected[bidder.BidderTemplateId] = 1
		}
	}
	for _, bidder := range bidderTemplates {
		if _, ok := mapBidderSelected[bidder.TableBidderTemplate.Id]; ok {
			assigns.BidderTemplatesSelected = append(assigns.BidderTemplatesSelected, bidder)
		} else {
			assigns.BidderTemplatesUnSelect = append(assigns.BidderTemplatesUnSelect, bidder)
		}
	}

	assigns.GamNetworks = new(model.GamNetwork).GetByUser(userLogin.Id)
	return ctx.Render("bidder/add", assigns, view.LAYOUTMain)
}

func (t *Bidder) AddPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidderAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}
	// Get Post Data
	inputs := payload.SystemCreate{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.System).Create(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

type AssignSystemEdit struct {
	assign.Schema
	Bidder                  model.BidderRecord
	MediaType               []model.MediaTypeRecord
	ListMediaTypeSelected   []int64
	Params                  []ParamBidder
	GamNetworks             []model.GamNetworkRecord
	BidderTemplatesUnSelect []model.BidderTemplateRecord
	BidderTemplatesSelected []model.BidderTemplateRecord
	ParamsForBidder         []model.PbBidderParamRecord
}

type ParamBidder struct {
	Name       string
	Type       string
	IsTemplate bool
}

func (t *Bidder) Edit(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidderEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignSystemEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idBidder, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	var row model.BidderRecord
	// Xử lý cho trang view và edit
	if assigns.Uri == config.URISystemBidderEdit {
		// Check bidder default
		isDefault := new(model.Bidder).CheckDefault(idBidder)
		if isDefault {
			// Nếu bidder là default thì get bidder k check user
			row = new(model.Bidder).GetByIdNoCheckUser(idBidder)
		} else {
			// Nếu như bidder k phải default check bidder đó theo user
			row = new(model.Bidder).GetById(idBidder, userLogin.Id)
		}
		if row.Id == 0 || row.IsLock == mysql.TYPEIsLockTypeLock {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
	} else if assigns.Uri == config.URISystemBidderView {
		// Check bidder default
		isDefault := new(model.Bidder).CheckDefault(idBidder)
		if isDefault {
			// Nếu bidder là default thì get bidder k check user
			row = new(model.Bidder).GetByIdNoCheckUser(idBidder)
		} else {
			// Nếu như bidder k phải default get bidder đó theo user
			row = new(model.Bidder).GetById(idBidder, userLogin.Id)
		}
		if row.Id == 0 || row.IsLock != mysql.TYPEIsLockTypeLock {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
	}

	assigns.Bidder = row
	if assigns.Bidder.Id == 0 {
		return ctx.Redirect(config.URISystemBidder)
		// return fmt.Errorf(GetLang(ctx).Errors.BidderError.NotFound.ToString())
	}
	// Get list bidder template phân loại bidder nào đã dược chọn và chưa được chọn
	bidderTemplates := new(model.BidderTemplate).GetAll()
	bidderSelected := new(model.Bidder).GetAllUser(userLogin.Id)
	mapBidderSelected := make(map[int64]int)
	for _, bidder := range bidderSelected {
		if bidder.BidderTemplateId != 1 && bidder.BidderTemplateId != 2 && bidder.BidderTemplateId != assigns.Bidder.BidderTemplateId {
			mapBidderSelected[bidder.BidderTemplateId] = 1
		}
	}
	for _, bidder := range bidderTemplates {
		if _, ok := mapBidderSelected[bidder.TableBidderTemplate.Id]; ok {
			assigns.BidderTemplatesSelected = append(assigns.BidderTemplatesSelected, bidder)
		} else {
			assigns.BidderTemplatesUnSelect = append(assigns.BidderTemplatesUnSelect, bidder)
		}
	}

	assigns.MediaType = new(model.MediaType).GetAll()
	assigns.GamNetworks = new(model.GamNetwork).GetByUser(userLogin.Id)
	//var mapParams map[string]string
	//var mapParamTemplate map[string]string
	var params []ParamBidder
	// Get param của bidder
	bidderParams := new(model.BidderParams).GetByBidderId(idBidder)
	for _, bidderParam := range bidderParams {
		param := ParamBidder{}
		param.Name = bidderParam.Name
		param.Type = bidderParam.Type
		if bidderParam.BidderTemplateId != 0 {
			param.IsTemplate = true
		} else {
			param.IsTemplate = false
		}
		params = append(params, param)
	}
	////Lấy ra bidder template sau đó parse để được các param mẫu
	//recordBidderTemplate := new(model.BidderTemplate).GetById(assigns.Bidder.BidderTemplateId)
	//err = json.Unmarshal([]byte(recordBidderTemplate.Params), &mapParamTemplate)
	////Lấy ra các param trong bidder của pub trong đó bao gồm cả các param template
	//err = json.Unmarshal([]byte(assigns.Bidder.Params), &mapParams)
	//if err == nil {
	//	//Xử lý add param hiển thị của param template trước và đánh dấu đây là param template
	//	for k, v := range mapParamTemplate {
	//		params = append(params, ParamBidder{
	//			Name:       k,
	//			Type:       v,
	//			IsTemplate: true,
	//		})
	//	}
	//	//Add tiếp param của bidder và loại bỏ các param mẫu đã có sẵn đồng thời đánh dâu các param được add k nằm trong param template
	//	for k, v := range mapParams {
	//		if _, ok := mapParamTemplate[k]; !ok {
	//			params = append(params, ParamBidder{
	//				Name:       k,
	//				Type:       v,
	//				IsTemplate: false,
	//			})
	//		}
	//	}
	//}

	recordsMediaType := new(model.RlBidderMediaType).GetByBidderId(idBidder)
	for _, v := range recordsMediaType {
		assigns.ListMediaTypeSelected = append(assigns.ListMediaTypeSelected, v.MediaTypeId)
	}

	pbBidder := new(model.PbBidder).GetPbBidderByBidderCode(row.BidderCode)
	if pbBidder.Id != 0 {
		assigns.ParamsForBidder = new(model.PbBidder).GetPbBidderParamsByBidder(pbBidder.Id)
	}
	for i, param := range assigns.ParamsForBidder {
		switch param.Type {
		case "float", "currency", "decimal", "number":
			assigns.ParamsForBidder[i].Type = "float"
			break
		case "int", "integer":
			assigns.ParamsForBidder[i].Type = "int"
			break
		case "bool", "boolean":
			assigns.ParamsForBidder[i].Type = "boolean"
			break
		}
	}

	assigns.Params = params
	if row.IsLock == 1 {
		assigns.Title = config.TitleWithPrefix("Edit Bidder")
		return ctx.Render("bidder/edit", assigns, view.LAYOUTMain)
	} else {
		assigns.Title = config.TitleWithPrefix("View Bidder")
		return ctx.Render("bidder/view", assigns, view.LAYOUTMain)
	}
}

func (t *Bidder) EditPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidderEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.SystemCreate{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	// Handle
	response := ajax.Responses{}
	newLineItem, errs := new(model.System).UpdateBidder(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = newLineItem
	}

	return ctx.JSON(response)
}

func (t *Bidder) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidderDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	system := new(model.System)
	notify := system.Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}

func (t *Bidder) AddTemplate(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidderAddTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id := ctx.Query("id")
	idBidderTemplate, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	bidderTemplateParams := new(model.BidderTemplateParams).GetByBidderTemplateId(idBidderTemplate)
	var params []ParamBidder
	for _, bidderTemplateParam := range bidderTemplateParams {
		params = append(params, ParamBidder{
			Name: bidderTemplateParam.Name,
			Type: bidderTemplateParam.Type,
		})
	}
	type JsonTemplate struct {
		Params        string  `json:"params"`
		ListMediaType []int64 `json:"list_media_type"`
	}
	data := htmlblock.Render("bidder/block_html/param.block.gohtml", params).String()

	var listMediaType []int64
	recordMediaTypes := new(model.RlBidderMediaType).GetByBidderTemplateId(idBidderTemplate)
	for _, v := range recordMediaTypes {
		listMediaType = append(listMediaType, v.MediaTypeId)
	}
	//fmt.Println(listMediaType)
	return ctx.JSON(JsonTemplate{
		Params:        data,
		ListMediaType: listMediaType,
	})
}

func (t *Bidder) View(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidderView)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignSystemEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idBidder, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	assigns.Bidder = new(model.Bidder).GetById(idBidder, userLogin.Id)
	if assigns.Bidder.Id == 0 {
		return ctx.Redirect(config.URISystemBidder)
		// return fmt.Errorf(GetLang(ctx).Errors.BidderError.NotFound.ToString())
	}
	// Get list bidder template phân loại bidder nào đã dược chọn và chưa được chọn
	bidderTemplates := new(model.BidderTemplate).GetAll()
	bidderSelected := new(model.Bidder).GetAllUser(userLogin.Id)
	mapBidderSelected := make(map[int64]int)
	for _, bidder := range bidderSelected {
		if bidder.BidderTemplateId != 1 && bidder.BidderTemplateId != 2 && bidder.BidderTemplateId != assigns.Bidder.BidderTemplateId {
			mapBidderSelected[bidder.BidderTemplateId] = 1
		}
	}
	for _, bidder := range bidderTemplates {
		if _, ok := mapBidderSelected[bidder.TableBidderTemplate.Id]; ok {
			assigns.BidderTemplatesSelected = append(assigns.BidderTemplatesSelected, bidder)
		} else {
			assigns.BidderTemplatesUnSelect = append(assigns.BidderTemplatesUnSelect, bidder)
		}
	}

	assigns.MediaType = new(model.MediaType).GetAll()
	assigns.GamNetworks = new(model.GamNetwork).GetByUser(userLogin.Id)
	//var mapParams map[string]string
	//var mapParamTemplate map[string]string
	var params []ParamBidder
	// Get param của bidder
	bidderParams := new(model.BidderParams).GetByBidderId(idBidder)
	for _, bidderParam := range bidderParams {
		param := ParamBidder{}
		param.Name = bidderParam.Name
		param.Type = bidderParam.Type
		if bidderParam.BidderTemplateId != 0 {
			param.IsTemplate = true
		} else {
			param.IsTemplate = false
		}
		params = append(params, param)
	}
	////Lấy ra bidder template sau đó parse để được các param mẫu

	recordsMediaType := new(model.RlBidderMediaType).GetByBidderId(idBidder)
	for _, v := range recordsMediaType {
		assigns.ListMediaTypeSelected = append(assigns.ListMediaTypeSelected, v.MediaTypeId)
	}

	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Bidder View")
	return ctx.Render("bidder/view", assigns, view.LAYOUTMain)
}

func (t *Bidder) UploadHandleXlsx(ctx *fiber.Ctx) error {
	//Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemBidderUploadXlsx)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	response := ajax.Responses{}
	errs, data := new(model.Bidder).SaveFileXlsx(ctx)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

func (t *Bidder) AddParamBidder(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISystemAddParam)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	demand := ctx.Query("demand")

	var paramBidderPB ParamBidderPB
	paramBidderPB.Bidder = new(model.PbBidder).GetPbBidderByBidderCode(demand)
	if paramBidderPB.Bidder.Id != 0 {
		paramBidderPB.Params = new(model.PbBidder).GetPbBidderParamsByBidder(paramBidderPB.Bidder.Id)
	}
	for i, param := range paramBidderPB.Params {
		switch param.Type {
		case "float", "currency", "decimal", "number":
			paramBidderPB.Params[i].Type = "float"
			break
		case "int", "integer":
			paramBidderPB.Params[i].Type = "int"
			break
		case "bool", "boolean":
			paramBidderPB.Params[i].Type = "boolean"
			break
		}
	}
	// Bidder := new(model.PbBidder).GetPbBidderByBidderCode(demand)
	// BidderParamsPB := new(model.PbBidder).GetPbBidderParamsByBidder(Bidder.Id)
	data := htmlblock.Render("bidder/block_html/add_param_bidder.gohtml", paramBidderPB).String()
	return ctx.JSON(data)
}
