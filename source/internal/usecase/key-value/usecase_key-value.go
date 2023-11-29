package key_value

import (
	"github.com/lib/pq"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	keyValueRepo "source/internal/repo/key-value"
)

type UsecaseKeyValue interface {
	Create(record *model.KeyModel) (err error)
	Update(ID int64, record *model.KeyModel) (err error)
	Filter(input *InputFilterUC) (totalRecords int64, records []*model.KeyModel, err error)
	GetByID(id int64) (record *model.KeyModel, err error)
}

func NewKeyValueUC(repos *repo.Repositories) *keyValueUC {
	return &keyValueUC{repos: repos}
}

type keyValueUC struct {
	repos *repo.Repositories
}

type InputFilterUC struct {
	datatable.Request
	UserID      int64  `json:"user_id"`
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
}

func (t *keyValueUC) Filter(input *InputFilterUC) (totalRecords int64, records []*model.KeyModel, err error) {
	//=> khởi tạo input của Repo
	inputRepo := keyValueRepo.InputFilter{
		UserID: input.UserID,
		Search: input.QuerySearch,
		Offset: input.Start,
		Limit:  input.Length,
		Order:  input.OrderString(),
	}
	return t.repos.KeyValue.Filter(&inputRepo)
}

func (t *keyValueUC) Create(record *model.KeyModel) (err error) {
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
		inputIsExists = keyValueRepo.InputIsExists{Key: record.KeyName}
	)
	if isExists = t.repos.KeyValue.IsExists(&inputIsExists); isExists {
		return errors.NewWithID(`Key already exists`, "key")
	}

	//=> thực hiện insert vào DB bằng hàm Save của gorm
	if err = t.repos.KeyValue.Save(record); err != nil {
		return
	}

	//=> Tạo cronjob cho Custom Key Value tạo lên GAM
	cronjob := model.CronjobModel{
		Type:         model.TYPECronJobCreateKeyValueGAM,
		ListObjectID: pq.Int64Array{record.ID},
		Status:       model.StatusCronJobQueue,
	}
	if err = t.repos.CronJob.Save(&cronjob); err != nil {
		return
	}
	return
}

func (t *keyValueUC) Update(ID int64, newRecord *model.KeyModel) (err error) {
	//=> nếu không nhập ID vào payload & payload không truyền ID vào object model thì không đủ dữ liệu để xử lý.
	if !newRecord.IsFound() {
		return errors.New(`new record must have an ID field`)
	}

	//=> validate lại các trường dữ liệu trong struct model
	if err = newRecord.Validate(); err != nil {
		return
	}

	//=> gọi ra bản ghi hiện tại trong DB
	oldRecord, _ := t.repos.KeyValue.FindByID(ID)
	if !oldRecord.IsFound() {
		return errors.New(`record not found`)
	}

	//=> kiểm tra xem bản ghi này đã tồn tại trên DB chưa
	var (
		inputIsExists = keyValueRepo.InputIsExists{Key: newRecord.KeyName}
		isExists      bool
	)
	if isExists = t.repos.KeyValue.IsExists(&inputIsExists, newRecord.ID); isExists {
		return errors.New(`key already exists`)
	}

	//=> cập nhật các trường dữ liệu không thay đổi từ bản ghi trong DB vào bản ghi mới trc khi được update.
	newRecord.CreatedAt = oldRecord.CreatedAt

	//=> thực hiện lưu vào DB
	if err = t.repos.KeyValue.Save(newRecord); err != nil {
		return
	}

	//=> Tạo cronjob cho Custom Key Value tạo lên GAM
	cronjob := model.CronjobModel{
		Type:         model.TYPECronJobCreateKeyValueGAM,
		ListObjectID: pq.Int64Array{newRecord.ID},
		Status:       model.StatusCronJobQueue,
	}
	if err = t.repos.CronJob.Save(&cronjob); err != nil {
		return
	}
	return
}

func (t *keyValueUC) GetByID(id int64) (record *model.KeyModel, err error) {
	record, err = t.repos.KeyValue.FindByID(id)
	return
}
