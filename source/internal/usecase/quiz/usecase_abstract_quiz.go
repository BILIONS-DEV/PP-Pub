package quiz

import (
	"source/internal/entity/dto"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
)

type UsecaseQuiz interface {
	Create(payload dto.PayloadQuizSubmit, userLogin model.User, userAdmin model.User) (response dto.Response)
	Update(payload dto.PayloadQuizSubmit, userLogin model.User, userAdmin model.User) (response dto.Response)
	GetById(id int64) (record model.Quiz, err error)
	GetByFilters(payload dto.QuizFilterPayload, userLogin model.User) (response datatable.Response, err error)
	Delete(payload dto.Delete, userLogin model.User, userAdmin model.User) (response dto.ResponseDelete, err error)
}
