package controller

import (
	"encoding/json"
	"fmt"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/helpers"
	"source/pkg/htmlblock"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LineItem struct{}

type AssignLineItemIndex struct {
	assign.Schema
	Params        payload.LineItemIndex
	Domain        []model.SearchLineItem
	Format        []model.SearchLineItem
	Size          []model.SearchLineItem
	AdTag         []model.SearchLineItem
	Country       []model.SearchLineItem
	Device        []model.SearchLineItem
	DomainSearch  []model.SearchLineItem
	FormatSearch  []model.SearchLineItem
	SizeSearch    []model.SearchLineItem
	AdTagSearch   []model.SearchLineItem
	CountrySearch []model.SearchLineItem
	DeviceSearch  []model.SearchLineItem
}

func (t *LineItem) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItem)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.LineItemIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignLineItemIndex{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Domain = new(model.LineItem).GetFilter("domain", userLogin.Id, params.Domain)
	assigns.Format = new(model.LineItem).GetFilter("format", userLogin.Id, params.AdFormat)
	assigns.Size = new(model.LineItem).GetFilter("size", userLogin.Id, params.AdSize)
	assigns.AdTag = new(model.LineItem).GetFilter("adtag", userLogin.Id, params.AdTag)
	assigns.Country = new(model.LineItem).GetFilter("country", userLogin.Id, params.Country)
	assigns.Device = new(model.LineItem).GetFilter("device", userLogin.Id, params.Device)
	assigns.Title = config.TitleWithPrefix("Line Item")
	return ctx.Render("line-item/index", assigns, view.LAYOUTMain)
}

func (t *LineItem) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItem)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.LineItemFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.LineItem).GetByFilters(&inputs, userLogin, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *LineItem) Choose(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemChoose)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Choose a bidder")
	return ctx.Render("line-item/choose", assigns, view.LAYOUTMain)
}

type AssignBidderAdd struct {
	assign.Schema
	ListBoxCollapse   []string
	Bidders           []model.BidderRecord
	ListAccountGoogle []model.BidderRecord
	Sizes             []model.AdSizeRecord
	GamNetworks       []model.GamNetworkRecord

	Tabs   []helpers.TabInfoStruct
	TabsJS string
}

func (t *LineItem) Add(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	collapse := new(model.LineItem).GetListBoxCollapse(userLogin.Id, 0, "line-item", "add")
	assigns := AssignBidderAdd{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Add Line Item")
	recBidders := new(model.Bidder).GetAllUser(userLogin.Id)
	for _, bidder := range recBidders {
		if bidder.BidderTemplateId == 1 {
			// Bỏ qua bidder adx system default nếu chưa được approve ít nhất 1 networkId
			if bidder.AccountType.IsAdx() && bidder.IsDefault == mysql.TypeOn && bidder.UserId == 0 {
				if !new(model.RlsBidderSystemInventory).CheckApproveBidderAdxByUser(bidder.Id, userLogin.Id) {
					continue
				}
			}
			assigns.ListAccountGoogle = append(assigns.ListAccountGoogle, bidder)
		} else {
			assigns.Bidders = append(assigns.Bidders, bidder)
		}
	}
	assigns.Sizes = new(model.AdSize).GetAllSizeForLineGoogle()
	assigns.GamNetworks = new(model.GamNetwork).GetForSelect(userLogin.Id)
	assigns.ListBoxCollapse = collapse

	assigns.Tabs = helpers.TabInfo

	bytesJson, err := json.Marshal(helpers.TabInfo)
	if err != nil {
		fmt.Println(err)
	}
	assigns.TabsJS = string(bytesJson)

	return ctx.Render("line-item/add", assigns, view.LAYOUTMain)
}

type AssignBidderEdit struct {
	assign.Schema
	Row          model.LineItemRecord
	BidderParams []payload.BidderInfo
	Tabs         []helpers.TabInfoStruct
	TabsJS       string

	StartDate              string
	EndDate                string
	Bidders                []model.BidderRecord
	ListBoxCollapse        []string
	ListAccountGoogle      []model.BidderRecord
	AccountGoogleSelected  model.BidderRecord
	Sizes                  []model.AdSizeRecord
	ListSizeSelected       []string
	LineItemAdsenseAdSlots []model.LineItemAdSenseAdSlotRecord
	GamNetworks            []model.GamNetworkRecord
	CheckShowAdsenseSlot   bool

	//Target
	AdSizes             []model.AdSizeRecord
	AdFormats           []model.AdTypeRecord
	AdTags              []model.InventoryAdTagRecord
	Inventories         []model.InventoryRecord
	Countries           []model.CountryRecord
	Devices             []model.DeviceRecord
	InventoryOfBidder   []int64
	AdTagOfBidder       []int64
	AdFormatOfBidder    []int64
	AdSizeOfBidder      []int64
	CountryOfBidder     []int64
	DeviceOfBidder      []int64
	AdSizesIncluded     []model.AdSizeRecord
	AdFormatsIncluded   []model.AdTypeRecord
	AdTagsIncluded      []model.InventoryAdTagRecord
	InventoriesIncluded []model.InventoryRecord
	CountriesIncluded   []model.CountryRecord
	DevicesIncluded     []model.DeviceRecord

	IsMoreInventory   bool
	IsMoreTag         bool
	IsMoreAdSize      bool
	IsMoreAdFormat    bool
	IsMoreDevice      bool
	IsMoreGeography   bool
	InventoryLastPage bool
	TagLastPage       bool
	SizeLastPage      bool
	FormatLastPage    bool
	DeviceLastPage    bool
	GeoLastPage       bool
}

func (t *LineItem) Edit(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignBidderEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idLineItem, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	// Get record line item
	var row model.LineItemRecord
	// Xử lý cho trang view và edit
	if assigns.Uri == config.URILineItemEdit {
		row, _ = new(model.LineItem).GetById(idLineItem, userLogin.Id)
		if row.Id == 0 || row.IsLock != mysql.TYPEIsLockTypeUnlock {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
	} else if assigns.Uri == config.URILineItemView {
		row, _ = new(model.LineItem).GetById(idLineItem, userLogin.Id)
		if row.Id == 0 || row.IsLock != mysql.TYPEIsLockTypeLock {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
	}

	assigns.Row = row

	//Get list bidder của user trong đó có cả bidder của google
	recBidders := new(model.Bidder).GetAllUser(userLogin.Id)
	for _, bidder := range recBidders {
		if bidder.BidderTemplateId == 1 {
			assigns.ListAccountGoogle = append(assigns.ListAccountGoogle, bidder)
		} else {
			assigns.Bidders = append(assigns.Bidders, bidder)
		}
	}
	//Get account google selected từ bảng line_item_account
	recLineItemAccount, _ := new(model.LineItemAccount).GetByLineItem(row.Id)
	assigns.AccountGoogleSelected = new(model.Bidder).GetByIdNoCheckUser(recLineItemAccount.BidderId)

	collapse := new(model.LineItem).GetListBoxCollapse(userLogin.Id, idLineItem, "line-item", "edit")
	assigns.ListBoxCollapse = collapse
	//Parse json param from record
	var jsonBidderInfos []payload.BidderInfo
	var bidderInfos []model.LineItemBidderInfoRecord
	mapBidderIndex := make(map[int64]int)
	bidderInfos = new(model.LineItemBidderInfo).GetBidderInfoByLineItem(row.Id)
	for _, bidderInfo := range bidderInfos {
		var bidderParams []payload.BidderParam
		params := new(model.LineItemBidderParams).GetParams(bidderInfo.Id)
		for _, param := range params {
			// Get bidder để check param required
			bidder := new(model.Bidder).GetById(param.BidderId, userLogin.Id)
			isrequired := new(model.PbBidder).IsParamRequiredByBidder(bidder.BidderCode, param.Name)
			bidderParams = append(bidderParams, payload.BidderParam{
				Param:      param,
				IsRequired: isrequired,
			})
		}
		bidder := new(model.Bidder).GetById(bidderInfo.BidderId, userLogin.Id)
		mapBidderIndex[bidderInfo.BidderId] = mapBidderIndex[bidderInfo.BidderId] + 1
		jsonBidderInfos = append(jsonBidderInfos, payload.BidderInfo{
			BidderId:       bidderInfo.BidderId,
			BidderName:     bidderInfo.Name,
			ConfigType:     bidderInfo.ConfigType,
			BidderType:     bidderInfo.BidderType,
			BidderIndex:    mapBidderIndex[bidderInfo.BidderId],
			BidderTemplate: bidder.BidderTemplateId,
			BidderParams:   bidderParams,
			Link:           new(model.PbBidder).GetLinkByBidderCode(bidderInfo.Name),
		})
	}
	assigns.BidderParams = jsonBidderInfos
	layoutISO := "01/02/2006"
	assigns.StartDate = row.StartDate.Time.Format(layoutISO)
	assigns.EndDate = row.EndDate.Time.Format(layoutISO)
	targets := new(model.Target).GetTargetLineItem(idLineItem)
	mapInventory := make(map[int64]int)
	mapAdFormat := make(map[int64]int)
	mapAdSize := make(map[int64]int)
	mapAdTag := make(map[int64]int)
	mapGeo := make(map[int64]int)
	mapDevice := make(map[int64]int)
	for _, target := range targets {
		if target.InventoryId != 0 {
			mapInventory[target.InventoryId] = 1
		}
		if target.AdFormatId != 0 {
			mapAdFormat[target.AdFormatId] = 1
		}
		if target.AdSizeId != 0 {
			mapAdSize[target.AdSizeId] = 1
		}
		if target.TagId != 0 {
			mapAdTag[target.TagId] = 1
		}
		if target.GeoId != 0 {
			mapGeo[target.GeoId] = 1
		}
		if target.DeviceId != 0 {
			mapDevice[target.DeviceId] = 1
		}
	}

	// Append list sau khi đã lọc bỏ những id trùng nhau với map
	for inventoryId, _ := range mapInventory {
		assigns.InventoryOfBidder = append(assigns.InventoryOfBidder, inventoryId)
	}
	for adFormatId, _ := range mapAdFormat {
		assigns.AdFormatOfBidder = append(assigns.AdFormatOfBidder, adFormatId)
	}
	for adSizeId, _ := range mapAdSize {
		assigns.AdSizeOfBidder = append(assigns.AdSizeOfBidder, adSizeId)
	}
	for adTagId, _ := range mapAdTag {
		assigns.AdTagOfBidder = append(assigns.AdTagOfBidder, adTagId)
	}
	for geoId, _ := range mapGeo {
		assigns.CountryOfBidder = append(assigns.CountryOfBidder, geoId)
	}
	for deviceId, _ := range mapDevice {
		assigns.DeviceOfBidder = append(assigns.DeviceOfBidder, deviceId)
	}
	var listIdInventory []int64
	for _, v := range assigns.InventoryOfBidder {
		var recInventory model.InventoryRecord
		if v != -1 {
			recInventory, _ = new(model.Inventory).GetByIdForFilter(v, userLogin.Id)
		}
		if recInventory.Id == 0 {
			continue
		}
		assigns.InventoriesIncluded = append(assigns.InventoriesIncluded, recInventory)
		listIdInventory = append(listIdInventory, recInventory.Id)
	}

	var listIdTag []int64
	for _, v := range assigns.AdTagOfBidder {
		var recAdTags model.InventoryAdTagRecord
		if v != -1 {
			recAdTags = new(model.InventoryAdTag).GetByIdForFilter(v)
		}
		if recAdTags.Id == 0 {
			continue
		}
		assigns.AdTagsIncluded = append(assigns.AdTagsIncluded, recAdTags)
		listIdTag = append(listIdTag, recAdTags.Id)
	}

	var listIdFormat []int64
	for _, v := range assigns.AdFormatOfBidder {
		var recAdFormat model.AdTypeRecord
		if v != -1 {
			recAdFormat = new(model.AdType).GetById(v)
		}
		if recAdFormat.Id == 0 {
			continue
		}
		assigns.AdFormatsIncluded = append(assigns.AdFormatsIncluded, recAdFormat)
		listIdFormat = append(listIdFormat, recAdFormat.Id)
	}

	var listIdSize []int64
	for _, v := range assigns.AdSizeOfBidder {
		var recAdSize model.AdSizeRecord
		if v != -1 {
			recAdSize = new(model.AdSize).GetById(v)
		}
		if recAdSize.Id == 0 {
			continue
		}
		assigns.AdSizesIncluded = append(assigns.AdSizesIncluded, recAdSize)
		listIdSize = append(listIdSize, recAdSize.Id)
	}

	var listIdGeo []int64
	for _, v := range assigns.CountryOfBidder {
		var recCountry model.CountryRecord
		if v != -1 {
			recCountry = new(model.Country).GetById(v)
		}
		if recCountry.Id == 0 {
			continue
		}
		assigns.CountriesIncluded = append(assigns.CountriesIncluded, recCountry)
		listIdGeo = append(listIdGeo, recCountry.Id)
	}

	var listIdDevice []int64
	for _, v := range assigns.DeviceOfBidder {
		var recDevice model.DeviceRecord
		if v != -1 {
			recDevice = new(model.Device).GetById(v)
		}
		if recDevice.Id == 0 {
			continue
		}
		assigns.DevicesIncluded = append(assigns.DevicesIncluded, recDevice)
		listIdDevice = append(listIdDevice, recDevice.Id)
	}

	assigns.Inventories, assigns.IsMoreInventory, assigns.InventoryLastPage = new(model.Inventory).LoadMoreDataPageEdit(userLogin.Id, listIdInventory)
	assigns.AdSizes, assigns.IsMoreAdSize, assigns.SizeLastPage = new(model.AdSize).LoadMoreDataPageEdit(listIdSize)
	assigns.AdFormats, assigns.IsMoreAdFormat, assigns.FormatLastPage = new(model.AdType).LoadMoreDataPageEdit(listIdFormat)
	assigns.AdTags, assigns.IsMoreTag, _, assigns.TagLastPage = new(model.InventoryAdTag).LoadMoreDataPageEdit(payload.FilterTarget{
		Inventory: assigns.InventoryOfBidder,
		Format:    assigns.AdFormatOfBidder,
		Size:      assigns.AdSizeOfBidder,
	}, userLogin.Id, listIdTag)
	assigns.Countries, assigns.IsMoreGeography, assigns.GeoLastPage = new(model.Country).LoadMoreDataPageEdit(listIdGeo)
	assigns.Devices, assigns.IsMoreDevice, assigns.DeviceLastPage = new(model.Device).LoadMoreDataPageEdit(listIdDevice)
	assigns.Sizes = new(model.AdSize).GetAllSizeForLineGoogle()
	lineItemAdsenseAdslots, _ := new(model.LineItemAdsenseAdSlot).GetByLineItem(row.Id)
	for _, lineItemAdsenseAdslot := range lineItemAdsenseAdslots {
		assigns.ListSizeSelected = append(assigns.ListSizeSelected, lineItemAdsenseAdslot.Size)
	}
	assigns.LineItemAdsenseAdSlots = lineItemAdsenseAdslots
	assigns.GamNetworks = new(model.GamNetwork).GetForSelect(userLogin.Id)
	if assigns.Row.ServerType == mysql.TYPEServerTypeGoogle && assigns.AccountGoogleSelected.AccountType == mysql.TYPEAccountTypeAdsense && assigns.Row.ConnectionType == mysql.TYPEConnectionTypeLineItems && assigns.Row.GamLineItemType == mysql.TYPEGamLineItemTypeDisplay {
		assigns.CheckShowAdsenseSlot = true
	}

	assigns.Tabs = helpers.TabInfo

	bytesJson, err := json.Marshal(helpers.TabInfo)
	if err != nil {
		fmt.Println(err)
	}
	assigns.TabsJS = string(bytesJson)

	if row.AutoCreate == 1 { // 1 - là trường hợp auto create dùng layout view riêng
		assigns.Title = config.TitleWithPrefix("View Line Item " + row.Name)
		return ctx.Render("line-item/view", assigns, view.LAYOUTMain)
	} else {
		assigns.Title = config.TitleWithPrefix("Edit Line Item " + row.Name)
		return ctx.Render("line-item/edit", assigns, view.LAYOUTMain)
	}
}

type AssignBidderCreate struct {
	assign.Schema
	Countries []model.CountryRecord
	AdSizes   []model.AdSizeRecord
	AdFormats []model.AdTypeRecord
}

func (t *LineItem) AddPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.LineItemAdd{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	newLineItem, errs := new(model.LineItem).AddLineItem(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = newLineItem
	}

	return ctx.JSON(response)
}

func (t *LineItem) EditPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.LineItemAdd{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		fmt.Println(err)
		return err
	}
	// Handle
	response := ajax.Responses{}
	newLineItem, errs := new(model.LineItem).UpdateLineItem(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = newLineItem
	}

	return ctx.JSON(response)
}

func (t *LineItem) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemDelete)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	notify := new(model.LineItem).Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}

type DataLoadMoreBidder struct {
	Key             int                `json:"key"`
	Search          string             `json:"search"`
	BiddersSelected []payload.LineItem `json:"bidders_selected"`
}

type TemplateLoadMoreBidder struct {
	Id        int64  `gorm:"column:id" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	IsChecked bool
	Placement string `json:"placement"`
	Publisher string `json:"publisher"`
}

type BidderParam struct {
	Bidders          []model.BidderRecord
	BidderId         int64
	BidderName       string
	BidderIndex      int
	BidderType       string
	BidderTemplateId int64
	Params           []model.BidderParamsRecord
	ParamValue       []payload.BidderInfo
	Link             string
}

func (t *LineItem) BidderParam(ctx *fiber.Ctx) (err error) {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemLoadParam)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	bidderId := ctx.Query("id")
	bidderName := ctx.Query("name")
	bidderIndex := ctx.Query("index")
	bidderType := ctx.Query("type")

	var bidderParam BidderParam
	bidderParam.BidderId, err = strconv.ParseInt(bidderId, 10, 64)
	if err != nil {
		return err
	}
	recBidder := new(model.Bidder).GetById(bidderParam.BidderId, userLogin.Id)
	bidderParam.BidderTemplateId = recBidder.BidderTemplateId
	bidderParam.BidderName = bidderName
	bidderParam.BidderIndex, _ = strconv.Atoi(bidderIndex)
	bidderParam.Params = new(model.BidderParams).GetByBidderId(bidderParam.BidderId)
	bidderParam.BidderType = bidderType
	bidderParam.Link = new(model.PbBidder).GetLinkByBidderCode(recBidder.BidderCode)
	data := htmlblock.Render("line-item/block_html/bidder_param_block.gohtml", bidderParam).String()
	var moreData payload.LoadMoreData
	moreData.Data = data
	return ctx.JSON(moreData)
}

func (t *LineItem) SearchDomain(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISearchDomain)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	search := ctx.Query("q")
	data := new(model.LineItem).GetInventoryByName(search, userLogin.Id)
	return ctx.JSON(data)
}

func (t *LineItem) SearchAdFormat(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISearchAdFormat)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	search := ctx.Query("q")
	data := new(model.LineItem).GetAdFormatByName(search)
	return ctx.JSON(data)
}

func (t *LineItem) SearchAdSize(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISearchAdSize)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	search := ctx.Query("q")
	data := new(model.LineItem).GetAdSizeByName(search)
	return ctx.JSON(data)
}

func (t *LineItem) SearchAdTag(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISearchAdTag)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	search := ctx.Query("q")
	data := new(model.LineItem).GetAdTagByName(search, userLogin.Id)
	return ctx.JSON(data)
}

func (t *LineItem) SearchDevice(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISearchDevice)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	search := ctx.Query("q")
	data := new(model.LineItem).GetDeviceByName(search)
	return ctx.JSON(data)
}

func (t *LineItem) SearchCountry(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URISearchCountry)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	search := ctx.Query("q")
	data := new(model.LineItem).GetCountryByName(search)
	return ctx.JSON(data)
}

func (t *LineItem) Collapse(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemCollapse)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := model.PageCollapseRecord{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	inputs.UserId = userLogin.Id
	inputs.PageCollapse = "line-item"
	errs := new(model.PageCollapse).HandleCollapse(inputs)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *LineItem) CheckParam(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemCheckParam)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	name := ctx.Query("name")
	param := ctx.Query("param")
	check := new(model.LineItem).CheckParamRequired(name, param)
	return ctx.JSON(fiber.Map{
		"required": check,
	})
}

func (t *LineItem) ListLinkedGam(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URILineItemListLinkedGam)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	ppAdx := ctx.Query("ppAdx")
	bidderAdxId, _ := strconv.ParseInt(ctx.Query("bidderAdxId"), 10, 64)
	var listIdNetworkAccept []int64

	var networks []model.GamNetworkRecord
	if ppAdx == "true" {
		// Get list rlsConnectionMCM approve
		listRls := new(model.RlsConnectionMCM).GetByBidderId(bidderAdxId)
		for _, rls := range listRls {
			if rls.Status == mysql.TYPEConnectionMCMTypeAccept {
				gamNetwork := new(model.GamNetwork).GetByNetworkId(rls.NetworkId, userLogin.Id)
				listIdNetworkAccept = append(listIdNetworkAccept, gamNetwork.Id)
			}
		}
		mysql.Client.Where("id in ? AND user_id = ?", listIdNetworkAccept, userLogin.Id).Find(&networks)
	} else {
		networks = new(model.GamNetwork).GetByUser(userLogin.Id)
	}

	return ctx.JSON(fiber.Map{
		"list_gam": networks,
	})
}

type ParamBidderPB struct {
	Bidder model.PbBidderRecord
	Params []model.PbBidderParamRecord
}

func (t *LineItem) AddParamBidder(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAddParamBidder)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	demand := ctx.Query("demand")

	var paramBidderPB ParamBidderPB
	paramBidderPB.Bidder = new(model.PbBidder).GetPbBidderByBidderCode(demand)
	if paramBidderPB.Bidder.Id != 0 {
		paramBidderPB.Params = new(model.PbBidder).GetPbBidderParamsByBidder(paramBidderPB.Bidder.Id)
	}
	// Bidder := new(model.PbBidder).GetPbBidderByBidderCode(demand)
	// BidderParamsPB := new(model.PbBidder).GetPbBidderParamsByBidder(Bidder.Id)
	data := htmlblock.Render("line-item/block_html/add_param_bidder.gohtml", paramBidderPB).String()
	return ctx.JSON(data)
}
