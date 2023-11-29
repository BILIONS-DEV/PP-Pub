package controllerv2

import (
	"github.com/gofiber/fiber/v2"
	history2 "source/apps/history"
	"source/core/technology/mysql"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"time"
)

type history struct {
	*handler
}

func (h *handler) InitRoutesHistory(app fiber.Router) {
	this := history{h}
	app.Post("/history", this.Filter)
	app.Post("/history/:hid", this.Detail)
}

func (t *history) Detail(ctx *fiber.Ctx) (err error) {

	var ID int
	if ID, err = ctx.ParamsInt("hid"); err != nil {
		return ctx.JSON(dto.MakeResponseError(err))
	}

	// get user login
	user := getUserLogin(ctx)

	// check login
	if !user.IsFound() {
		return ctx.JSON(dto.MakeResponseErrorString("Permission denied"))
	}
	record := t.useCases.History.DetailByUser(int64(ID), user.ID)

	// Convert dữ liệu sang phiên bản cũ để dùng lại package history của Khánh dùng cho chuyển đổi dữ liệu của JSON
	// Vì phần này là code không đồng bộ (source mới nhưng dùng package history ở source cũ) nên xử lý ngoài controller để không bị ảnh hưởng đến logic trong usecase
	recordOldVersion := mysql.TableHistory{
		Id:          record.ID,
		CreatorId:   record.CreatorID,
		Title:       record.Title,
		Object:      record.Object,
		ObjectId:    record.ObjectID,
		ObjectType:  mysql.TYPEObjectType(record.ObjectType),
		DetailType:  record.DetailType,
		App:         record.App,
		UserId:      record.UserID,
		InventoryId: record.InventoryID,
		BidderId:    record.BidderID,
		OldData:     record.OldData,
		NewData:     record.NewData,
		CreatedAt:   record.CreatedAt,
	}
	compareData, err := history2.CompareDataHistory(recordOldVersion)
	if err != nil {
		return ctx.JSON(dto.MakeResponseError(err))
	}
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return ctx.JSON(dto.MakeResponseError(err))
	}
	createTime := record.CreatedAt.In(loc).Format("01/02/2006 15:04:05 AM")
	output := struct {
		Row        model.History
		Compare    []history2.ResponseCompare
		CreateTime string
	}{
		Row:        record,
		Compare:    compareData,
		CreateTime: createTime,
	}
	return ctx.JSON(dto.MakeResponseSuccess(output))
}

func (t *history) Filter(ctx *fiber.Ctx) (err error) {
	// get user login
	user := getUserLogin(ctx)

	// check login
	if !user.IsFound() {
		return ctx.JSON(dto.MakeResponseErrorString("Permission denied"))
	}

	// get payload from http request
	var payload dto.PayloadHistoryFilter
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.JSON(dto.MakeResponseError(err))
	}

	// run usecase
	var (
		total     int64
		histories []model.History
	)
	if total, histories, err = t.useCases.History.FilterByUser(user.ID, &payload); err != nil {
		return ctx.JSON(dto.MakeResponseError(err))
	}
	// make response
	object := fiber.Map{
		"records":     histories,
		"totalRecord": total,
	}
	return ctx.JSON(dto.MakeResponseSuccess(object, "success"))
}
