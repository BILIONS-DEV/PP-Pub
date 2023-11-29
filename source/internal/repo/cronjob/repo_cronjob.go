package cronjob

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
)

type RepoCronJob interface {
	FindQueue(limit ...int) (record []model.CronjobModel, err error)
	Save(record *model.CronjobModel) (err error)
}

type cronJobRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewCronJobRepo(db *gorm.DB) *cronJobRepo {
	return &cronJobRepo{Db: db}
}

func (t *cronJobRepo) FindQueue(limit ...int) (records []model.CronjobModel, err error) {
	err = t.Db.
		//Debug().
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(limit) > 0 {
				return db.Limit(limit[0])
			}
			return db
		}).
		Where("status = ?", model.StatusCronJobQueue).
		First(&records).Error
	return
}

func (t *cronJobRepo) Save(record *model.CronjobModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}
