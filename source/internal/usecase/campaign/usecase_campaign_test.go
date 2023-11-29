package campaign

import (
	"reflect"
	"source/infrastructure/fakedb"
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/repo"
	// "source/internal/repo/campaign/mocks"
	"testing"
)

// func TestCampaignUsecase_AddCampaign2(t *testing.T) {
// 	recEx := model.CampaignModel{
// 		Name:          "BBBB",
// 		TrafficSource: "Mgid",
// 		DemandSource:  "system1",
// 	}
// 	campRepo := mocks.NewRepoCampaign(t)
// 	campRepo.On("IsExist", recEx.Name).Return(false)
// 	campRepo.On("Save", &recEx).Return(nil)
// 	repos := &repo.Repositories{
// 		Campaign: campRepo,
// 	}
// 	campUC := NewCampaignUsecase(repos, nil)
//
// 	rec := model.CampaignModel{
// 		Name:          "BBBB",
// 		TrafficSource: "Mgid",
// 		DemandSource:  "system1",
// 	}
// 	err := campUC.AddCampaign2(&rec)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	// repos := repo.Repositories{Campaign: }
// 	// campUC := NewCampaignUsecase()
// }
//
// func TestCampaignUsecase_EditCampaign2(t *testing.T) {
// 	rec := model.CampaignModel{
// 		ID:            51,
// 		Name:          "BBBB",
// 		TrafficSource: "Mgid",
// 		DemandSource:  "system1",
// 	}
// 	campRepo := mocks.NewRepoCampaign(t)
// 	// campRepo.On("IsExist", "BBBB").Return(false)
// 	campRepo.On("Save", rec).Return(nil)
// 	repos := &repo.Repositories{
// 		Campaign: campRepo,
// 	}
// 	campUC := NewCampaignUsecase(repos, nil)
//
// 	err := campUC.EditCampaign2(&rec)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func Test_campaignUsecase_GetCampaignById(t1 *testing.T) {

	db, _ := fakedb.NewMysqlAff()
	rps := repo.NewRepositories(&repo.Deps{
		Db:    db,
		Cache: nil,
		Kafka: nil,
	})

	type fields struct {
		repos *repo.Repositories
		Trans *lang.Translation
	}
	type args struct {
		idSearch int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRow model.CampaignModel
	}{
		// TODO: Add test cases.
		{"test 1", fields{repos: rps, Trans: nil}, args{idSearch: 2}, model.CampaignModel{ID: 1}},
		{"test 2", fields{repos: rps, Trans: nil}, args{idSearch: 216}, model.CampaignModel{ID: 216}},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &campaignUsecase{
				repos: tt.fields.repos,
				Trans: tt.fields.Trans,
			}
			if gotRow := t.GetCampaignById(tt.args.idSearch); !reflect.DeepEqual(gotRow, tt.wantRow) {
				t1.Errorf("GetCampaignById() = %v \n, want %v\n", gotRow, tt.wantRow)
			}
		})
	}
}
