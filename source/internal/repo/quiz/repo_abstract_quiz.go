package quiz

import (
	"source/internal/entity/dto"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
)

type RepoQuiz interface {
	Migrate()
	FindByUser(userID int64) (records []model.Quiz)
	Create(record *model.Quiz) (err error)
	Update(record *model.Quiz) (err error)
	FindById(id int64) (record model.Quiz, err error)
	FindByFilters(payload *dto.QuizFilterPayload, userLogin model.User) (response datatable.Response, err error)
	Delete(record *model.Quiz) (response dto.ResponseDelete, err error)
}
