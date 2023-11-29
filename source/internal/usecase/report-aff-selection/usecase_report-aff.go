package reportAffSelection

import (
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/repo"
	reportAffSelection "source/internal/repo/report-aff-selection"
)

type UsecaseReportAffSelection interface {
	Filter(input *InputFilter) (totalRecords int64, records []*model.ReportAffSelectionModel, recordTotal *model.ReportAffSelectionModel, err error)
	GetAllCampaign() (records []*model.ReportAffSelectionModel, err error)
	GetAllTrafficSource() (records []*model.ReportAffSelectionModel, err error)
	GetAllPublisher() (records []*model.ReportAffSelectionModel, err error)
	GetAllSection() (records []*model.ReportAffSelectionModel, err error)
}

type reportAffUsecase struct {
	repos *repo.Repositories
	Trans *lang.Translation
}

func NewReportAffSelectionUsecase(repos *repo.Repositories) *reportAffUsecase {
	return &reportAffUsecase{repos: repos}
}

type InputFilter struct {
	datatable.Request
	UserID         int64
	QuerySearch    string
	StartDate      string
	EndDate        string
	Campaigns      interface{}
	TrafficSources interface{}
	PublisherID    interface{}
	GroupBy        interface{}
	SectionID      interface{}
}

func (t *reportAffUsecase) Filter(input *InputFilter) (totalRecords int64, records []*model.ReportAffSelectionModel, recordTotal *model.ReportAffSelectionModel, err error) {
	//=> khởi tạo input của Repo
	inputRepo := reportAffSelection.InputFilter{
		UserID:         input.UserID,
		Search:         input.QuerySearch,
		Offset:         input.Start,
		Limit:          input.Length,
		Order:          input.OrderString(),
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Campaigns:      input.Campaigns,
		TrafficSources: input.TrafficSources,
		PublisherID:    input.PublisherID,
		SectionID:      input.SectionID,
		GroupBy:        input.GroupBy,
	}
	return t.repos.ReportAffSelection.Filter(&inputRepo)
}

func (t *reportAffUsecase) GetAllCampaign() (records []*model.ReportAffSelectionModel, err error) {
	return t.repos.ReportAffSelection.FindAllByGroup("campaign_id")
}

func (t *reportAffUsecase) GetAllTrafficSource() (records []*model.ReportAffSelectionModel, err error) {
	return t.repos.ReportAffSelection.FindAllByGroup("traffic_source")
}

func (t *reportAffUsecase) GetAllPublisher() (records []*model.ReportAffSelectionModel, err error) {
	return t.repos.ReportAffSelection.FindAllByGroup("publisher_id")
}

func (t *reportAffUsecase) GetAllSection() (records []*model.ReportAffSelectionModel, err error) {
	return t.repos.ReportAffSelection.FindAllByGroup("section_id")
}
