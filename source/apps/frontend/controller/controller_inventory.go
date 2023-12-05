package controller

import (
	"encoding/json"
	"fmt"
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
	"time"
)

type Inventory struct{}

type AssignInventory struct {
	assign.Schema
	Params                     payload.InventoryIndex
	ParamsAdTag                payload.InventoryAdTagIndex
	Row                        model.InventoryRecord
	Config                     model.InventoryConfigRecord
	ModuleUserId               []model.ModuleUserIdRecord
	ModuleDefault              []payload.ModuleUserIdAssign
	Modules                    []payload.ModuleUserIdAssign
	ListIdOfInventory          []int64
	ListModuleIdDefault        []int64
	ListBoxCollapse            []string
	ListBoxCollapseSticky      []string
	Tab                        int64
	GamNetworks                []model.GamNetworkRecord
	AdTypes                    []model.AdTypeRecord
	Bidders                    []model.InventoryConnectionDemands
	CountConnectWaiting        int
	AdsTxtMissingLineSyncError string
	AdsTxtMissingLines         []model.AdsTxtMissingLine
}

func (t *Inventory) Setup(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventorySetup)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id := ctx.Query("id")
	// start := ctx.Query("start")
	// length := ctx.Query("length")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	tab, err := strconv.ParseInt(ctx.Query("tab"), 10, 64)
	if err != nil {
		tab = 1
	}
	fmt.Println("data: ", userLogin.Permission)
	// bỏ tab 2,3,5 với user permission managed service
	if userLogin.Permission == mysql.UserPermissionManagedService || userLogin.Permission == mysql.UserPermissionSubPublisher {
		if tab == 2 || tab == 3 || tab == 5 {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
	}
	assigns := AssignInventory{Schema: assign.Get(ctx)}
	if err := ctx.QueryParser(&assigns.Params); err != nil {
		return err
	}
	params := payload.InventoryAdTagIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	row, err := new(model.Inventory).GetById(idSearch, userLogin.Id)
	if err != nil || row.Id == 0 || row.Status != 1 {
		return fmt.Errorf(GetLang(ctx).Errors.InventoryError.NotFound.ToString())
	}
	cf := new(model.InventoryConfig).GetByInventoryId(idSearch)
	if cf.Id == 0 {
		cf = new(model.InventoryConfig).MakeRowDefault(row.Id)
	}
	_ = row.ScanAdsTxt()

	module, err := new(model.ModuleUserId).GetAll()
	if err != nil {
		return err
	}
	collapse := new(model.Inventory).GetListBoxCollapse(userLogin.Id, idSearch, "inventory")
	assigns.Tab = tab
	assigns.Row = row
	assigns.Config = cf
	assigns.ModuleUserId = module
	assigns.ParamsAdTag = params
	assigns.ListBoxCollapse = collapse
	assigns.GamNetworks = new(model.GamNetwork).GetByUser(userLogin.Id)
	assigns.AdTypes = new(model.AdType).GetAll(userLogin, userAdmin)
	assigns.Bidders = new(model.InventoryConnectionDemand).GetByInventory(row)
	assigns.CountConnectWaiting = 0
	for _, bidderId := range new(model.RlsBidderSystemInventory).GetListIdBidderApprove(row.Name) {
		bidder := new(model.Bidder).GetByIdNoCheckUser(bidderId)
		if bidder.Id == 0 || bidder.BidderTemplateId == 1 { // Bỏ qua nếu k tìm thấy bidder hoặc là bidder google
			continue
		}
		status := new(model.InventoryConnectionDemand).GetStatus(row.Id, bidder.Id)
		if status == 2 {
			assigns.CountConnectWaiting++
		}
	}

	adsMissingLine, syncError := row.GetAllMissingAdsTxt()
	assigns.AdsTxtMissingLines = adsMissingLine
	assigns.AdsTxtMissingLineSyncError = syncError

	assigns.Title = config.TitleWithPrefix("Setup Inventory")
	assigns.LANG.Title = "Setup Inventory"
	return ctx.Render("supply/setup", assigns, view.LAYOUTMain)
}

func (t *Inventory) SetupConfig(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventorySetup)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.GeneralInventory{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	errs := new(model.InventoryConfig).SetupInventoryConfig(inputs, "config", userLogin.Id, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *Inventory) SetupConsent(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryConsent)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	response := ajax.Responses{}
	inputs := payload.GeneralInventory{}
	if err := ctx.BodyParser(&inputs); err != nil {
		response.Status = ajax.ERROR
		response.Errors = append(response.Errors, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	errs := new(model.InventoryConfig).SetupInventoryConfig(inputs, "consent", userLogin.Id, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *Inventory) SetupUserId(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotFound)
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryUserId)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.GeneralInventory{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	errs := new(model.InventoryConfig).SetupInventoryConfig(inputs, "userid", userLogin.Id, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *Inventory) Delete(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotFound)
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryDelete)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	notify := new(model.Inventory).Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}

func (t *Inventory) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventory)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.InventoryIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignInventory{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Supplies")
	assigns.LANG.Title = "Supplies"
	return ctx.Render("supply/index", assigns, view.LAYOUTMain)
}

func (t *Inventory) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventory)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.InventoryFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.Inventory).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *Inventory) FilterAdTag(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryAdTag)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.InventoryAdTagFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.InventoryAdTag).GetByFilters(&inputs, userLogin, userAdmin, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *Inventory) FilterConnection(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryConnection)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.InventoryConnectionDemandFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.InventoryConnectionDemand).GetByFilters(&inputs, userLogin)
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *Inventory) Submit(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotFound)
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventorySubmit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := new(payload.InventorySubmit)
	if err := ctx.BodyParser(inputs); err != nil {
		return err
	}
	response := new(model.Inventory).Submit(inputs, userLogin.Id, GetLang(ctx))
	return ctx.JSON(response)
}

func (t *Inventory) LoadModuleParam(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryLoadParam)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	moduleId := ctx.Query("id")
	// moduleName := ctx.Query("name")
	moduleIndex := ctx.Query("index")
	var moduleParam payload.ModuleParam
	moduleParam.ModuleId, err = strconv.ParseInt(moduleId, 10, 64)
	if err != nil {
		return err
	}
	// moduleParam.ModuleName = moduleName
	moduleParam.ModuleIndex, _ = strconv.Atoi(moduleIndex)
	module := new(model.ModuleUserId).GetById(moduleParam.ModuleId)
	moduleParam.ModuleName = module.Name
	var mapParam []payload.ParamModuleUserId
	err = json.Unmarshal([]byte(module.Params), &mapParam)
	if err != nil {
		return err
	}
	var storage []payload.StorageModuleUserId
	err = json.Unmarshal([]byte(module.Storage), &storage)
	if err != nil {
		return err
	}
	moduleParam.Params = mapParam
	moduleParam.Storage = storage
	data := htmlblock.Render("supply/index/block_module.gohtml", moduleParam).String()
	var moreData payload.LoadMoreData
	moreData.Data = data
	return ctx.JSON(moreData)
}

func (t *Inventory) Collapse(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryCollapse)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := model.PageCollapseRecord{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	inputs.UserId = userLogin.Id
	inputs.PageCollapse = "inventory"
	errs := new(model.PageCollapse).HandleCollapse(inputs)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

type CopyTag struct {
	Inventory                model.InventoryRecord
	ListAdTagDisplay         []Tag
	ListAdTagOutstream       []Tag
	ListAdTagInstream        []Tag
	ListAdTagSticky          []Tag
	ListAdTagPinZone         []Tag
	ListAdTagPlayZoneRelated []Tag
	ListAdTagPlayZoneQuiz    []Tag
	VastAdTagInstream        []VastTag
	VastAdTagOutstream       []VastTag
	ListTagTool              []model.AdTagRecord
	ListQuiz                 []model.QizPostsRecord
	ListAdTagNative          []Tag
}

type Tag struct {
	AdTag   model.InventoryAdTagRecord
	Size    model.AdSizeRecord
	IsVideo bool
}
type VastTag struct {
	AdTag     model.InventoryAdTagRecord
	VastVpaid Vast
	VastS2s   Vast
}

type Vast struct {
	VastUrl  string
	VastName string
}

func (t *Inventory) CopyAdTag(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryCopyAdTag)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id, _ := strconv.ParseInt(ctx.Query("id"), 10, 64)
	var copyTag CopyTag
	copyTag.Inventory, err = new(model.Inventory).GetByIdSystem(id)
	if err != nil {
		return fmt.Errorf(GetLang(ctx).Errors.InventoryError.NotFound.ToString())
	}
	adTags := new(model.InventoryAdTag).GetByInventory(copyTag.Inventory.Id)
	for _, tag := range adTags {
		var size model.AdSizeRecord

		// Tạo macro cho link vast
		pageUrl := "[PAGE_URL]"
		playerSize := "[PLAYER_SIZE]"
		switch tag.Renderer {
		case mysql.TYPERendererJWPlayer:
			pageUrl = "__page-url__"
			playerSize = "__player-width__x__player-height__"
			break
		case mysql.TYPERendererVideoJS:
			pageUrl = "{player.pageUrl}"
			playerSize = "{player.width}x{player.height}"
			break
		case mysql.TYPERendererFlowPlayer:
			pageUrl = "page_url"
			playerSize = "player_widthxplayer_height"
			break
		}

		if tag.Type == mysql.TYPEDisplay {
			size = new(model.AdSize).GetById(tag.PrimaryAdSize)
			copyTag.ListAdTagDisplay = append(copyTag.ListAdTagDisplay, Tag{
				AdTag: tag,
				Size:  size,
			})
		} else if tag.Type == mysql.TYPEInStream {
			if tag.Renderer != mysql.TYPERendererPubPower && tag.Renderer != mysql.TYPERendererOverlayAd && tag.Status == mysql.TypeStatusAdTagLive {
				vastTag := VastTag{
					AdTag: tag,
					VastVpaid: Vast{
						VastUrl: "https://cdn.vlitag.com/vpaid/w/" + copyTag.Inventory.Uuid + "?tagid=" + strconv.FormatInt(tag.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize,
						//VastUrl:  "https://nc.pubpowerplatform.io/vpaid/w/" + copyTag.Inventory.Uuid + "?tagid=" + strconv.FormatInt(tag.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize,
						VastName: "Vast Vpaid",
					},
					VastS2s: Vast{
						VastUrl:  "https://ss-pbs.quantumdex.io/vast?description_url=" + pageUrl + "&sz=" + playerSize + "&tagid=" + strconv.FormatInt(tag.Id, 10) + "&ad_type=video&tfcd=0&vpmute=1&vpos=preroll&vpa=auto",
						VastName: "Vast S2S",
					},
				}
				copyTag.VastAdTagInstream = append(copyTag.VastAdTagInstream, vastTag)
			} else {
				copyTag.ListAdTagInstream = append(copyTag.ListAdTagInstream, Tag{
					AdTag: tag,
					Size:  size,
				})
			}
		} else if tag.Type == mysql.TYPEOutStream {
			if tag.Renderer != mysql.TYPERendererPubPower && tag.Status == mysql.TypeStatusAdTagLive {
				vastTag := VastTag{
					AdTag: tag,
					VastVpaid: Vast{
						VastUrl: "https://cdn.vlitag.com/vpaid/w/" + copyTag.Inventory.Uuid + "?tagid=" + strconv.FormatInt(tag.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize,
						//VastUrl:  "https://nc.pubpowerplatform.io/vpaid/w/" + copyTag.Inventory.Uuid + "?tagid=" + strconv.FormatInt(tag.Id, 10) + "&page_url=" + pageUrl + "&sz=" + playerSize,
						VastName: "Vast Vpaid",
					},
					VastS2s: Vast{
						VastUrl:  "https://ss-pbs.quantumdex.io/vast?description_url=" + pageUrl + "&sz=" + playerSize + "&tagid=" + strconv.FormatInt(tag.Id, 10) + "&ad_type=video&tfcd=0&vpmute=1&vpos=preroll&vpa=auto",
						VastName: "Vast S2S",
					},
				}
				copyTag.VastAdTagOutstream = append(copyTag.VastAdTagOutstream, vastTag)
			} else {
				copyTag.ListAdTagOutstream = append(copyTag.ListAdTagOutstream, Tag{
					AdTag: tag,
					Size:  size,
				})
			}
		} else if tag.Type == mysql.TYPEStickyBanner {
			copyTag.ListAdTagSticky = append(copyTag.ListAdTagSticky, Tag{
				AdTag: tag,
				Size:  size,
			})
		} else if tag.Type == mysql.TYPETopArticles {
			var isVideo bool
			contentTopArticle := new(model.ContentTopArticles).GetContent(copyTag.Inventory.Name, tag.RelatedContent.Code(), tag.FeedUrl)
			if contentTopArticle.Id != 0 {
				isVideo = true
			}
			copyTag.ListAdTagPinZone = append(copyTag.ListAdTagPinZone, Tag{
				AdTag:   tag,
				Size:    size,
				IsVideo: isVideo,
			})
		} else if tag.Type == mysql.TYPEPlayZone {
			if tag.ContentType == mysql.TYPEContentTypeRelated {
				copyTag.ListAdTagPlayZoneRelated = append(copyTag.ListAdTagPlayZoneRelated, Tag{
					AdTag: tag,
					Size:  size,
				})
			} else if tag.ContentType == mysql.TYPEContentTypeQuiz {
				copyTag.ListAdTagPlayZoneQuiz = append(copyTag.ListAdTagPlayZoneQuiz, Tag{
					AdTag: tag,
					Size:  size,
				})
			}
		} else if tag.Type == mysql.TYPEPNative {
			copyTag.ListAdTagNative = append(copyTag.ListAdTagNative, Tag{
				AdTag: tag,
				Size:  size,
			})
		}
	}
	copyTag.ListTagTool = new(model.AdTag).GetAdTagForTool(id)
	copyTag.ListQuiz, _ = new(model.Quiz).GetAllQuizForUserSelect(userLogin.Email)
	// data := htmlblock.Render("supply/box.copy_adtag.gohtml", copyTag).String()
	return ctx.Render("supply/copy_adtag/box.copy_adtag", copyTag)
}

type BuildScript struct {
	TimeStamp       int64
	TagDesktop      ScriptTag
	TagMobile       ScriptTag
	PlaceHolder     mysql.TypeOnOff
	PlaceHolderText string
	TextColor       string
	BorderColor     string
	IsShowStyle     bool
}

type ScriptTag struct {
	Slot  string
	Size  model.AdSizeRecord
	AdTag model.AdTagRecord
}

func (t *Inventory) BuildScript(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryBuildScript)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	var buildScript payload.BuildScript
	_ = ctx.BodyParser(&buildScript)

	adTagDesktop, _ := new(model.AdTag).GetDetail(buildScript.TagDesktop)
	tagDesktop := ScriptTag{
		Slot:  "pw_" + strconv.FormatInt(buildScript.TagDesktop, 10),
		Size:  new(model.AdSize).GetById(adTagDesktop.PrimaryAdSize),
		AdTag: adTagDesktop,
	}
	adTagMobile, _ := new(model.AdTag).GetDetail(buildScript.TagMobile)
	tagMobile := ScriptTag{
		Slot:  "pw_" + strconv.FormatInt(buildScript.TagMobile, 10),
		Size:  new(model.AdSize).GetById(adTagMobile.PrimaryAdSize),
		AdTag: adTagMobile,
	}
	var isShowStyle bool
	if tagDesktop.Size.Width != tagMobile.Size.Width || tagDesktop.Size.Height != tagMobile.Size.Height {
		isShowStyle = true
	}
	// data := htmlblock.Render("supply/box.copy_adtag.gohtml", copyTag).String()
	return ctx.Render("supply/copy_adtag/script", BuildScript{
		TimeStamp:       time.Now().Unix(),
		TagDesktop:      tagDesktop,
		TagMobile:       tagMobile,
		PlaceHolder:     buildScript.PlaceHolder,
		PlaceHolderText: buildScript.PlaceHolderText,
		TextColor:       buildScript.TextColor,
		BorderColor:     buildScript.BorderColor,
		IsShowStyle:     isShowStyle,
	})
}

func (t *Inventory) ChangeStatusConnection(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventoryChangeStatusConnection)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	var inputs payload.PayloadInventoryChangeStatusConnection
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	resp := new(model.InventoryConnectionDemand).ChangeStatus(inputs, userLogin, userAdmin)
	return ctx.JSON(resp)
}
