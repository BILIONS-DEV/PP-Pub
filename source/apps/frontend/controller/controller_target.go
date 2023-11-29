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
	"source/pkg/htmlblock"
	"strconv"
)

type Target struct{}

type AssignRule struct {
	assign.Schema
	Params          payload.RuleIndex
	Countries       []model.CountryRecord
	Devices         []model.DeviceRecord
	Data            []model.InventoryRecord
	IsMoreData      bool
	IsMoreDevice    bool
	IsMoreGeography bool
}

func (t *Target) Test(ctx *fiber.Ctx) error {
	assigns := AssignRule{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Test form")
	return ctx.Render("target/floor", assigns, view.LAYOUTMain)
}

func (t *Target) LoadMoreData(ctx *fiber.Ctx) error {
	user := GetUserLogin(ctx)
	if !user.IsFound() {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	key := ctx.Query("key")
	// userId, _ := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	search := ctx.Query("search")
	option := ctx.Query("option")
	filter := ctx.Query("filter")
	selected := ctx.Query("selected")
	page := ctx.Query("page")
	var isSystem bool
	if user.IsAdmin() {
		querySystem := ctx.Query("isSystem")
		if querySystem == "true" {
			isSystem = true
		} else {
			isSystem = false
		}
	}
	var filterTarget payload.FilterTarget
	err := json.Unmarshal([]byte(filter), &filterTarget)
	if err != nil {
		return err
	}
	var selectedTarget []int64
	err = json.Unmarshal([]byte(selected), &selectedTarget)
	if err != nil {
		return err
	}
	switch option {
	case "domain":
		var rows []model.InventoryRecord
		var isMoreData, lastPage bool
		if isSystem {
			rows, isMoreData, lastPage = new(model.Inventory).LoadMoreDataSystem(key, search, selectedTarget)
		} else {
			rows, isMoreData, lastPage = new(model.Inventory).LoadMoreData(key, search, user.Id, selectedTarget)
		}
		if page == "identity" {
			rows = new(model.Target).HandleForIdentity(rows, user)
		}
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/inventory_block.gohtml", rows).String()
		//fmt.Printf("\n rows: %+v \n", rows)
		//fmt.Printf("\n data: %+v \n", data)
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)
	case "adformat":
		rows, isMoreData, lastPage := new(model.AdType).LoadMoreData(key, search, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/adFormat_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)
	case "adsize":
		rows, isMoreData, lastPage := new(model.AdSize).LoadMoreData(key, search, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/adSize_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)
	case "adtag":
		var rows []model.InventoryAdTagRecord
		var isMoreData bool
		var listFilter []int64
		var lastPage bool
		if isSystem {
			rows, isMoreData, listFilter, lastPage = new(model.InventoryAdTag).LoadMoreDataSystem(key, search, filterTarget, selectedTarget)
		} else {
			rows, isMoreData, listFilter, lastPage = new(model.InventoryAdTag).LoadMoreData(key, search, filterTarget, user.Id, selectedTarget)
		}
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/adTag_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		if listFilter == nil {
			loadMoreData.ListFilter = []int64{}
		} else {
			loadMoreData.ListFilter = listFilter
		}
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)
	case "device":
		rows, isMoreData, lastPage := new(model.Device).LoadMoreData(key, search, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/device_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)
	case "geography":
		rows, isMoreData, lastPage := new(model.Country).LoadMoreData(key, search, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/geo_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)

	case "language":
		rows, isMoreData, lastPage := new(model.Language).LoadMoreData(key, search, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/language_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)

	case "channels":
		rows, isMoreData, lastPage := new(model.Channels).LoadMoreData(key, search, user.Id, filterTarget, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/channels_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)

	case "category":
		rows, isMoreData, lastPage := new(model.Category).LoadMoreData(key, search, user.Id, filterTarget, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/category_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)

	case "keywords":
		rows, isMoreData, lastPage := new(model.ContentKeyword).LoadMoreData(key, search, user.Id, filterTarget, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/keywords_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)

	case "videos":
		rows, total, isMoreData, lastPage := new(model.Content).LoadMoreData(key, search, user.Id, filterTarget, selectedTarget)
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/videos_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		loadMoreData.TotalAll = total
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)

	default:
		var rows []model.InventoryRecord
		var isMoreData, lastPage bool
		if isSystem {
			rows, isMoreData, lastPage = new(model.Inventory).LoadMoreDataSystem(key, search, selectedTarget)
		} else {

			rows, isMoreData, lastPage = new(model.Inventory).LoadMoreData(key, search, user.Id, selectedTarget)
		}
		var loadMoreData payload.LoadMoreData
		data := htmlblock.Render("target/block_html/inventory_block.gohtml", rows).String()
		loadMoreData.Data = data
		loadMoreData.IsMoreData = isMoreData
		loadMoreData.CurrentPage, _ = strconv.Atoi(key)
		loadMoreData.LastPage = lastPage
		loadMoreData.Total = len(rows)
		if search != "" {
			loadMoreData.IsSearch = true
		}
		return ctx.JSON(loadMoreData)
	}
}

func (t *Target) FilterAdTag(ctx *fiber.Ctx) error {
	user := GetUserLogin(ctx)
	if !user.IsFound() {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	type AdTagFilter struct {
		ListAdTagID []int64 `json:"tagId"`
		IsSystem    bool    `json:"is_system"`
		ServerType  int     `json:"serverType"`
	}
	var adTagFilter AdTagFilter
	if err := json.Unmarshal(ctx.Body(), &adTagFilter); err != nil {
		return err
	}

	type Data struct {
		InventoryId []payload.ListTarget `json:"inventory_id"`
		AdFormatId  []payload.ListTarget `json:"ad_format_id"`
		AdSizeId    []payload.ListTarget `json:"ad_size_id"`
	}
	var data Data
	//Dùng map để lọc các id trùng nhau
	mapInventory := make(map[int64]payload.ListTarget)
	mapAdFormat := make(map[int64]payload.ListTarget)
	mapAdSize := make(map[int64]payload.ListTarget)
	for _, tagId := range adTagFilter.ListAdTagID {
		recordAdTag := new(model.InventoryAdTag).GetByIdForFilter(tagId)
		if recordAdTag.Id == 0 {
			continue
		}

		var recordInventory model.InventoryRecord
		if adTagFilter.IsSystem {
			recordInventory, _ = new(model.Inventory).GetByIdForFilterSystem(recordAdTag.InventoryId)
		} else {
			recordInventory, _ = new(model.Inventory).GetByIdForFilter(recordAdTag.InventoryId, user.Id)
		}
		if recordInventory.Id != 0 {
			mapInventory[recordAdTag.InventoryId] = payload.ListTarget{
				Id:   recordInventory.Id,
				Name: recordInventory.Name,
			}
		}

		recordAdFormat := new(model.AdType).GetById(recordAdTag.Type.Int())
		if recordAdFormat.Id != 0 {
			mapAdFormat[recordAdTag.Type.Int()] = payload.ListTarget{
				Id:   recordAdFormat.Id,
				Name: recordAdFormat.Name,
			}
		}
		var recordAdSize model.AdSizeRecord
		if recordAdTag.Type == mysql.TYPEDisplay {
			recordAdSize = new(model.AdSize).GetById(recordAdTag.PrimaryAdSize)
			if recordAdSize.Id != 0 {
				mapAdSize[recordAdTag.PrimaryAdSize] = payload.ListTarget{
					Id:   recordAdSize.Id,
					Name: recordAdSize.Name,
				}
			}
		} else if recordAdTag.Type == mysql.TYPEStickyBanner {
			recordAdSize = new(model.AdSize).GetById(recordAdTag.PrimaryAdSize)
			if recordAdSize.Id != 0 {
				mapAdSize[recordAdTag.PrimaryAdSize] = payload.ListTarget{
					Id:   recordAdSize.Id,
					Name: recordAdSize.Name,
				}
			}
		}

	}
	//fmt.Println(mapInventory)
	for _, value := range mapInventory {
		data.InventoryId = append(data.InventoryId, value)
	}
	for _, value := range mapAdFormat {
		data.AdFormatId = append(data.AdFormatId, value)
	}
	if adTagFilter.ServerType != 2 {
		for _, value := range mapAdSize {
			data.AdSizeId = append(data.AdSizeId, value)
		}
	} else {
		data.AdSizeId = []payload.ListTarget{}
	}

	return ctx.JSON(data)
}

func (t *Target) LoadSelected(ctx *fiber.Ctx) error {
	user := GetUserLogin(ctx)
	if !user.IsFound() {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	option := ctx.Query("option")
	filter := ctx.Query("filter")
	var filterTarget payload.FilterTarget
	err := json.Unmarshal([]byte(filter), &filterTarget)
	if err != nil {
		return err
	}
	switch option {
	case "domain":
		selected, _ := new(model.Target).GetAllByUser(user.Id, "domain", filterTarget)
		data := htmlblock.Render("floor/block_html/select_all.gohtml", selected).String()
		var selectAll payload.SelectAll
		selectAll.Data = data
		selectAll.ListSelected = selected
		return ctx.JSON(selectAll)
	case "adformat":
		selected, _ := new(model.Target).GetAllByUser(user.Id, "adformat", filterTarget)
		data := htmlblock.Render("floor/block_html/select_all.gohtml", selected).String()
		var selectAll payload.SelectAll
		selectAll.Data = data
		selectAll.ListSelected = selected
		return ctx.JSON(selectAll)
	case "adsize":
		selected, _ := new(model.Target).GetAllByUser(user.Id, "adsize", filterTarget)
		data := htmlblock.Render("floor/block_html/select_all.gohtml", selected).String()
		var selectAll payload.SelectAll
		selectAll.Data = data
		selectAll.ListSelected = selected
		return ctx.JSON(selectAll)
	case "adtag":
		selected, listFilter := new(model.Target).GetAllByUser(user.Id, "adtag", filterTarget)
		data := htmlblock.Render("floor/block_html/select_all.gohtml", selected).String()
		var selectAll payload.SelectAll
		selectAll.Data = data
		selectAll.ListSelected = selected
		selectAll.ListFilter = listFilter
		return ctx.JSON(selectAll)
	case "device":
		selected, _ := new(model.Target).GetAllByUser(user.Id, "device", filterTarget)
		data := htmlblock.Render("floor/block_html/select_all.gohtml", selected).String()
		var selectAll payload.SelectAll
		selectAll.Data = data
		selectAll.ListSelected = selected
		return ctx.JSON(selectAll)
	case "geography":
		selected, _ := new(model.Target).GetAllByUser(user.Id, "geography", filterTarget)
		data := htmlblock.Render("floor/block_html/select_all.gohtml", selected).String()
		var selectAll payload.SelectAll
		selectAll.Data = data
		selectAll.ListSelected = selected
		return ctx.JSON(selectAll)
	default:
		selected, _ := new(model.Target).GetAllByUser(user.Id, "domain", filterTarget)
		data := htmlblock.Render("floor/block_html/select_all.gohtml", selected).String()
		var selectAll payload.SelectAll
		selectAll.Data = data
		selectAll.ListSelected = selected
		return ctx.JSON(selectAll)
	}
}

func (t *Target) LoadSize(ctx *fiber.Ctx) error {
	data := new(model.AdSize).GetAllApi()
	return ctx.JSON(data)
}
