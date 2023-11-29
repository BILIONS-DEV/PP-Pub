package quiz

import (
	"github.com/asaskevich/govalidator"
	"source/internal/entity/dto"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/repo"
)

type quiz struct {
	Repos *repo.Repositories
	Trans *lang.Translation
}

func NewQuizUsecase(repos *repo.Repositories, trans *lang.Translation) *quiz {
	return &quiz{Repos: repos, Trans: trans}
}

func (t *quiz) Create(payload dto.PayloadQuizSubmit, userLogin model.User, userAdmin model.User) (response dto.Response) {
	// validate
	errs := t.validate(payload)
	if len(errs) > 0 {
		response.Status = false
		response.Errors = errs
		return
	}

	// make record
	record := t.makeRecord(payload, userLogin)
	if err := t.Repos.Quiz.Create(&record); err != nil {
		response.Status = false
		response.Errors = append(response.Errors, dto.Error{
			ID:      "",
			Message: "error",
		})
		return
	}

	response.Status = true

	// Reset cache
	if err := t.Repos.Inventory.ResetCacheByUser(userLogin.ID); err != nil {
		response.Status = false
		response.Errors = append(response.Errors, dto.Error{
			ID:      "",
			Message: "cache error",
		})
	}
	return
}

func (t *quiz) Update(payload dto.PayloadQuizSubmit, userLogin model.User, userAdmin model.User) (response dto.Response) {
	// validate
	errs := t.validate(payload)
	if len(errs) > 0 {
		response.Status = false
		response.Errors = errs
		return
	}

	// get recordOld
	recordOld, _ := t.Repos.Quiz.FindById(payload.ID)
	payload.ContentType = recordOld.ContentType.Int()

	record := t.makeRecord(payload, userLogin)
	if err := t.Repos.Quiz.Update(&record); err != nil {
		response.Status = false
		response.Errors = append(response.Errors, dto.Error{
			ID:      "",
			Message: "error",
		})
		return
	}

	record, _ = t.Repos.Quiz.FindById(record.ID)
	response.Status = true
	response.Data = record

	// Reset cache
	if err := t.Repos.Inventory.ResetCacheByUser(userLogin.ID); err != nil {
		response.Status = false
		response.Errors = append(response.Errors, dto.Error{
			ID:      "",
			Message: "cache error",
		})
	}
	return
}

func (t *quiz) GetById(id int64) (record model.Quiz, err error) {
	record, err = t.Repos.Quiz.FindById(id)
	return
}

func (t *quiz) validate(payload dto.PayloadQuizSubmit) (errs []dto.Error) {
	if govalidator.IsNull(payload.Title) {
		errs = append(errs, dto.Error{
			ID:      "title",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	if payload.Category == 0 {
		errs = append(errs, dto.Error{
			ID:      "category",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	if len(payload.Questions) == 0 {
		errs = append(errs, dto.Error{
			ID:      "question",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	return
}

func (t *quiz) makeRecord(payload dto.PayloadQuizSubmit, user model.User) (record model.Quiz) {
	record.ID = payload.ID
	record.Title = payload.Title
	record.Illustration = payload.Illustration
	record.UserID = user.ID
	record.CategoryID = payload.Category
	record.ContentType = model.TYPEContentType(payload.ContentType)
	var questions []model.QuizQuestion
	for _, question := range payload.Questions {
		var answers []model.QuizAnswer
		for _, answer := range question.Answer {
			answers = append(answers, model.QuizAnswer{
				ID:        answer.ID,
				Answer:    answer.Text,
				AnswerImg: answer.Img,
				Correct:   answer.Correct,
			})
		}
		questions = append(questions, model.QuizQuestion{
			ID:         question.ID,
			Question:   question.Question,
			Background: question.Background,
			Answers:    answers,
		})
	}
	record.Questions = questions
	return
}

func (t *quiz) GetByFilters(payload dto.QuizFilterPayload, userLogin model.User) (response datatable.Response, err error) {
	response, err = t.Repos.Quiz.FindByFilters(&payload, userLogin)
	return
}

func (t *quiz) Delete(payload dto.Delete, userLogin model.User, userAdmin model.User) (response dto.ResponseDelete, err error) {
	record, _ := t.Repos.Quiz.FindById(payload.ID)
	response, err = t.Repos.Quiz.Delete(&record)
	return
}
