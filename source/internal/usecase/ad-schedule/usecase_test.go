package ad_schedule

import (
	"github.com/stretchr/testify/assert"
	"source/internal/entity/model"
	"source/internal/repo"
	"testing"
)

func TestAdScheduleUC_Detail(t *testing.T) {

	var id int64 = 7

	assertions := assert.New(t)

	adScheduleRepo := new(mockAdScheduleRepo)
	adScheduleRepo.On("FindByID", id).
		Return(&model.AdScheduleModel{ID: 9, Name: "Test"}, nil)

	uc := NewAdScheduleUC(&repo.Repositories{AdSchedule: adScheduleRepo})
	record, err := uc.Detail(id)

	adScheduleRepo.AssertExpectations(t)
	assertions.NoError(err)
	assertions.Equal(record.ID, id)
	assertions.Equal(record.Name, "Test")
}
