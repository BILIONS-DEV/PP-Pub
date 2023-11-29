package ad_schedule

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"source/internal/entity/model"
	"source/internal/repo/ad-schedule"
)

type mockAdScheduleRepo struct {
	mock.Mock
}

func (t *mockAdScheduleRepo) Migrate() {
	//TODO implement me
	panic("implement me")
}

func (t *mockAdScheduleRepo) DB() *gorm.DB {
	//TODO implement me
	panic("implement me")
}

func (t *mockAdScheduleRepo) Filter(input *ad_schedule.InputFilter) (totalRecord int64, records []*model.AdScheduleModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockAdScheduleRepo) FindByID(ID int64) (record *model.AdScheduleModel, err error) {
	args := t.Called(ID)
	return args.Get(0).(*model.AdScheduleModel), args.Error(1)
}

func (t *mockAdScheduleRepo) IsExists(input *ad_schedule.InputIsExists, IDs ...int64) (exists bool) {
	//TODO implement me
	panic("implement me")
}

func (t *mockAdScheduleRepo) Save(record *model.AdScheduleModel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockAdScheduleRepo) EmptyConfigs(configs []model.AdScheduleConfigModel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockAdScheduleRepo) DeleteByID(ID int64) error {
	//TODO implement me
	panic("implement me")
}
