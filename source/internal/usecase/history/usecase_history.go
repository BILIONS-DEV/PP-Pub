package history

import (
	"errors"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/repo"
	historyRepo "source/internal/repo/history"
	"time"
)

type UsecaseHistory interface {
	Filter(payload *dto.PayloadHistoryFilter) (total int64, histories []model.History, err error)
	Detail(ID int64) (record model.History)
	DetailByUser(ID, userID int64) (record model.History)
	FilterByUser(userID int64, payload *dto.PayloadHistoryFilter) (total int64, histories []model.History, err error)
	GetAllPage() (page []string)
	LoadHistories(payload *InputHistoryUC) (histories []model.History, err error)
}

type history struct {
	Repos *repo.Repositories
	Lang  *lang.Translation
}

func NewHistoryUsecase(repos *repo.Repositories, lang *lang.Translation) *history {
	return &history{Repos: repos, Lang: lang}
}

func (t history) DetailByUser(ID, userID int64) (record model.History) {
	record = t.Repos.History.FindByID(ID, userID)
	return
}

func (t *history) FilterByUser(userID int64, payload *dto.PayloadHistoryFilter) (total int64, histories []model.History, err error) {
	// validate dữ liệu cơ bản
	if err = payload.Validate(); err != nil {
		return
	}
	// chuyển đổi và validate dữ liệu time
	var startDate, endDate time.Time
	if payload.StartDate != "" && payload.EndDate != "" {
		if startDate, err = time.Parse("2006-01-02", payload.StartDate); err != nil {
			return
		}
		if endDate, err = time.Parse("2006-01-02", payload.EndDate); err != nil {
			return
		}
		if startDate.After(endDate) {
			err = errors.New("the start date must be the previous day")
			return
		}
	}
	// gọi vào repo để thực hiện filter
	total, histories, err = t.Repos.History.FindByUser(userID, &historyRepo.FilterByUserInput{
		Condition: payload.ToCondition(),
		Search:    payload.Search,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      payload.Page,
		Order:     "id DESC",
	})
	return
}

func (t *history) Filter(payload *dto.PayloadHistoryFilter) (total int64, histories []model.History, err error) {
	// validate dữ liệu cơ bản
	if err = payload.Validate(); err != nil {
		return
	}
	// chuyển đổi và validate dữ liệu time
	var startDate, endDate time.Time
	if payload.StartDate != "" && payload.EndDate != "" {
		if startDate, err = time.Parse("2006-01-02", payload.StartDate); err != nil {
			return
		}
		if endDate, err = time.Parse("2006-01-02", payload.EndDate); err != nil {
			return
		}
		if startDate.After(endDate) {
			err = errors.New("the start date must be the previous day")
			return
		}
	}
	// gọi vào repo để thực hiện filter
	total, histories, err = t.Repos.History.FindByFilter(&historyRepo.FilterInput{
		Condition:  payload.ToCondition(),
		Search:     payload.Search,
		ObjectPage: payload.ObjectPage,
		ObjectID:   payload.ObjectID,
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       payload.Page,
		Order:      "id DESC",
	})
	return
}

type InputHistoryUC struct {
	Id     string `json:"id" form:"id"`
	Object string `json:"object" form:"object"`
}

func (t *InputHistoryUC) Validate() (err error) {
	return
}

func (t *history) LoadHistories(payload *InputHistoryUC) (histories []model.History, err error) {
	// validate dữ liệu cơ bản
	if err = payload.Validate(); err != nil {
		return
	}

	// gọi vào repo để thực hiện filter
	histories, err = t.Repos.History.LoadHistoriesByObject(&historyRepo.InputObject{
		Id:         payload.Id,
		ObjectType: payload.Object,
	})
	return
}

func (t history) Detail(ID int64) (record model.History) {
	record = t.Repos.History.GetByID(ID)
	return
}

func (t *history) GetAllPage() (pages []string) {
	pages = t.Repos.History.GetAllHistoryPage()
	return
}
