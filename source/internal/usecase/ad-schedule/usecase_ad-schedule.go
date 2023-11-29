package ad_schedule

import (
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	adscheduleRepo "source/internal/repo/ad-schedule"
)

type UsecaseAdSchedule interface {
	Create(record *model.AdScheduleModel) (err error)
	Update(ID int64, record *model.AdScheduleModel) (err error)
	Filter(input *InputFilterUC) (totalRecords int64, records []*model.AdScheduleModel, err error)
	Detail(ID int64) (record *model.AdScheduleModel, err error)
	Delete(ID int64) (err error)
}

func NewAdScheduleUC(repos *repo.Repositories) *adScheduleUC {
	return &adScheduleUC{repos: repos}
}

type adScheduleUC struct {
	repos *repo.Repositories
}

// InputFilterUC : vì không có gì thay đổi nên lớp data của usecase kế thừa lại struct của DTO cho tiện,
// tuy nhiên vẫn phải tạo struct data của lớp usecase để rèn luyện thói quen & đề phòng trường hợp có field dữ liệu không giống vs DTO.
type InputFilterUC struct {
	*dto.PayloadAdScheduleFilter
}

func (t *adScheduleUC) Filter(input *InputFilterUC) (totalRecords int64, records []*model.AdScheduleModel, err error) {
	//=> khởi tạo input của Repo
	inputRepo := adscheduleRepo.InputFilter{
		UserID:                input.UserID,
		ClientTypes:           input.ClientTypes,
		AdBreakTypes:          input.AdBreakTypes,
		AdBreakConfigTypeJoin: input.AdBreakConfigType,
		Search:                input.Search,
		Page:                  input.Page,
		Order:                 input.Order,
	}
	return t.repos.AdSchedule.Filter(&inputRepo)
}

func (t *adScheduleUC) Delete(ID int64) (err error) {
	return t.repos.AdSchedule.DeleteByID(ID)
}

func (t *adScheduleUC) Detail(ID int64) (record *model.AdScheduleModel, err error) {
	return t.repos.AdSchedule.FindByID(ID)
}

func (t *adScheduleUC) Create(record *model.AdScheduleModel) (err error) {
	//=> bản ghi mới không được có ID, nếu không hàm Save của gorm sẽ hiểu là 1 lệnh update
	if record.IsFound() {
		return errors.New(`ID field is not approved`)
	}

	//=> validate lại các trường dữ liệu trong struct model
	if err = record.Validate(); err != nil {
		return
	}

	//=> kiểm tra xem tên bản ghi đã tồn tại trong DB hay chưa
	var (
		isExists      bool
		inputIsExists = adscheduleRepo.InputIsExists{Name: record.Name}
	)
	if isExists = t.repos.AdSchedule.IsExists(&inputIsExists); isExists {
		return errors.NewWithID(`schedule already existss`, "name")
	}

	//=> thực hiện insert vào DB bằng hàm Save của gorm
	if err = t.repos.AdSchedule.Save(record); err != nil {
		return
	}

	return
}

func (t *adScheduleUC) Update(ID int64, newRecord *model.AdScheduleModel) (err error) {
	//=> nếu không nhập ID vào payload & payload không truyền ID vào object model thì không đủ dữ liệu để xử lý.
	if !newRecord.IsFound() {
		return errors.New(`new record must have an ID field`)
	}

	//=> validate lại các trường dữ liệu trong struct model
	if err = newRecord.Validate(); err != nil {
		return
	}

	//=> gọi ra bản ghi hiện tại trong DB
	oldRecord, _ := t.repos.AdSchedule.FindByID(ID)
	if !oldRecord.IsFound() {
		return errors.New(`record not found`)
	}

	//=> kiểm tra xem bản ghi này đã tồn tại trên DB chưa
	var (
		inputIsExists = adscheduleRepo.InputIsExists{Name: newRecord.Name}
		isExists      bool
	)
	if isExists = t.repos.AdSchedule.IsExists(&inputIsExists, newRecord.ID); isExists {
		return errors.New(`schedule already exists`)
	}

	//=> nếu bản ghi trong db vẫn còn relationship thì xóa đi
	if len(oldRecord.AdBreakConfigs) > 0 {
		if err = t.repos.AdSchedule.EmptyConfigs(oldRecord.AdBreakConfigs); err != nil {
			return
		}
	}

	//=> cập nhật các trường dữ liệu không thay đổi từ bản ghi trong DB vào bản ghi mới trc khi được update.
	newRecord.CreatedAt = oldRecord.CreatedAt

	//=> thực hiện lưu vào DB
	if err = t.repos.AdSchedule.Save(newRecord); err != nil {
		return
	}

	return
}
