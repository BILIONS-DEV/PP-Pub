package controllerv2

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/view"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	errors2 "source/internal/errors"
	adschedule "source/internal/usecase/ad-schedule"
	"strconv"
)

type advertising struct {
	*handler
}

func (h *handler) InitRoutesAdvertising(app fiber.Router) {
	t := advertising{handler: h}
	app.Get(config.URIAdvertisingSchedulesAdd, t.Add)

	app.Post(config.URIAdvertisingSchedules, t.IndexPost)
	app.Post(config.URIAdvertisingSchedulesAdd, t.AddPost)
	app.Post(config.URIAdvertisingSchedulesEdit, t.EditPost)
	app.Post(config.URIAdvertisingSchedulesDetail, t.DetailPost)
	app.Post(config.URIAdvertisingSchedulesDelete, t.DeletePost)

	app.Get("/t", func(ctx *fiber.Ctx) error {

		var listErrs []error
		for i := 0; i < 3; i++ {
			idString := strconv.Itoa(i)
			err := errors2.New("err origin: " + idString)
			err = errors2.Wrap(err, "err trans: "+idString, idString)
			listErrs = append(listErrs, err)
		}
		listErrs = append(listErrs, errors2.New("fuck"))
		listErrs = append(listErrs, errors2.Wrap(nil, "err nil"))
		listErrs = append(listErrs, errors2.NewWithID("my error", "my_id"))

		fail := dto.Fail(listErrs...)
		fail = dto.Fail(errors.New("error default"))

		myErr := errors.New("deo hieu kieu gi")
		myErrWrap := errors2.Wrap(myErr, "deo hieu translate...")
		myErrWrap2 := errors2.Wrap(nil, "deo hieu translate...")
		//fail = dto.Fail(errors2.Origin(myErrWrap))

		return ctx.JSON(fiber.Map{
			"myErrWrap2": myErrWrap2.Error(),
			"Origin2":    errors2.Origin(myErrWrap2).Error(),
			"myErrWrap":  myErrWrap.Error(),
			"Origin":     errors2.Origin(myErrWrap).Error(),
			"tran":       dto.Fail(myErrWrap),
			"orgi":       dto.Fail(errors2.Origin(myErrWrap)),
		})

		//origin := errors2.Origin(listErrs[2])
		//fmt.Printf("\n org: %+v \n", origin)
		return ctx.JSON(fiber.Map{
			"fail": fail,
			//"org":  origin.Error(),
		})
	})
}

type AssignAdSchedulesAdd struct {
	Assign
}

func (t *advertising) Add(ctx *fiber.Ctx) (err error) {

	// assigns
	assigns := AssignQuizAdd{Assign: newAssign(ctx, "Add Ad Schedules")}
	assigns.Categories, _ = t.useCases.Category.GetAll()
	return ctx.Render("ad_schedules/add", assigns, view.LAYOUTMain)
}

func (t *advertising) IndexPost(ctx *fiber.Ctx) (err error) {
	var payload dto.PayloadAdScheduleFilter
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseError(err),
		)
	}

	if errs := payload.Validate(); len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseErrorWithID(errs...),
		)
	}

	var inputUC = adschedule.InputFilterUC{PayloadAdScheduleFilter: &payload}
	var records []*model.AdScheduleModel
	var totalRecords int64
	if totalRecords, records, err = t.useCases.AdSchedule.Filter(&inputUC); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(err),
		)
	}

	return ctx.JSON(
		dto.Done(fiber.Map{
			"totalRecord": totalRecords,
			"records":     records,
		}, "Schedule has been filtered"),
	)
}

// AddPost : ở hàm này không truyền dto vào usecase nữa bởi vì sẽ không đảm bảo được logic khi test độc lập 1 function usecase.
// ví dụ: lúc trước khi để payload.Validate() ở controller -> xong vẫn truyền payload vào usecase -> trong usecase gọi hàm payload.ToModel()
// => lúc chạy test thì sẽ gọi trực tiếp hàm usecase.Create(payload) -> thì thiếu mất logic validate input của hàm payload.Input()
// => vẫn insert vào db thông qua hàm usecase nhưng không đúng logic
// Giải pháp: bắt buộc DTO phải chuyển đổi sang dữ liệu của usecase chứ không được truyền trực tiếp DTO vào usecase nữa.
//
// Cũng không thể cho hàm payload.Validate() vào usecase vì nó liên quan đến xử lý đầu ra dto.Response.[]Error
// Nếu cho vào usecase thì phần biến đổi dữ liệu DTO lại phải xử lý trong đó => sai rule của Clean Architecture
func (t *advertising) AddPost(ctx *fiber.Ctx) (err error) {

	var payload dto.PayloadAdSchedule
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseError(err),
		)
	}

	var errs []error
	if errs = payload.Validate(); len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(errs...),
		)
	}

	var record = payload.ToModel()
	if err = t.useCases.AdSchedule.Create(&record); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(err),
		)
	}

	return ctx.JSON(
		dto.Done(record, "Schedule has been created"),
	)
}

func (t *advertising) EditPost(ctx *fiber.Ctx) (err error) {

	var payload dto.PayloadAdSchedule

	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(err),
		)
	}

	var errs []error
	if errs = payload.Validate(); len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(errs...),
		)
	}

	//if payload.ID == 0 {
	//	return ctx.Status(fiber.StatusBadRequest).JSON(
	//		dto.Fail(errors2.New(`ID is required`)),
	//	)
	//}

	var record = payload.ToModel()
	if err = t.useCases.AdSchedule.Update(payload.ID, &record); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(err),
		)
	}

	return ctx.JSON(
		dto.Done(record, "Schedule has been updated"),
	)
}

func (t *advertising) DetailPost(ctx *fiber.Ctx) (err error) {
	IdString := ctx.Params("id")
	if IdString == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseErrorString("record not found"),
		)
	}

	var ID int64
	if ID, err = strconv.ParseInt(IdString, 10, 64); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseError(err),
		)
	}

	var record *model.AdScheduleModel
	if record, _ = t.useCases.AdSchedule.Detail(ID); !record.IsFound() {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseErrorString("record not found"),
		)
	}

	return ctx.JSON(
		dto.MakeResponseSuccess(record, ""),
	)
}

func (t *advertising) DeletePost(ctx *fiber.Ctx) (err error) {
	IdString := ctx.Params("id")
	if IdString == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseErrorString("record not found"),
		)
	}

	var ID int64
	if ID, err = strconv.ParseInt(IdString, 10, 64); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseError(err),
		)
	}

	if err = t.useCases.AdSchedule.Delete(ID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseError(err),
		)
	}

	return ctx.JSON(dto.MakeResponseSuccess(nil, "Schedule has been deleted"))
}
