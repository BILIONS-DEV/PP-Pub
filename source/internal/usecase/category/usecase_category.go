package category

import (
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/repo"
)

type category struct {
	Repos *repo.Repositories
	Trans *lang.Translation
}

func (t *category) GetAll() (records []model.Category, err error) {
	records, err = t.Repos.Category.FindAll()
	return
}

func NewCategoryUsecase(repos *repo.Repositories, trans *lang.Translation) *category {
	return &category{Repos: repos, Trans: trans}
}
