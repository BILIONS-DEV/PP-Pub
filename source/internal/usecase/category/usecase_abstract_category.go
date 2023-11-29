package category

import "source/internal/entity/model"

type UsecaseCategory interface {
	GetAll() (records []model.Category, err error)
}
