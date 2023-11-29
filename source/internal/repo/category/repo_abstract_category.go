package category

import "source/internal/entity/model"

type RepoCategory interface {
	FindAll() (records []model.Category, err error)
}
