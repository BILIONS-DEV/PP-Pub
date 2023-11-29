package key_value_gam

import (
	"gorm.io/gorm"
	"source/internal/entity/model"
)

type RepoKeyValueGam interface {
	FindKeyInGAM(networkID int64, keyName string) (record *model.KeyValueGamModel, err error)
	FindKeyValueInGAM(networkID int64, keyName string, value string) (record *model.KeyValueGamModel, err error)
	Save(record *model.KeyValueGamModel) (err error)
}

type keyValueGamRepo struct {
	Db *gorm.DB
}

func NewKeyValueRepo(DB *gorm.DB) *keyValueGamRepo {
	return &keyValueGamRepo{Db: DB}
}

func (t *keyValueGamRepo) FindKeyValueInGAM(networkID int64, keyName string, value string) (record *model.KeyValueGamModel, err error) {
	err = t.Db.
		//Debug().
		Where("network_id = ?", networkID).
		Where("name = ?", keyName).
		Where("value = ?", value).
		Find(&record).
		Error
	return
}

func (t *keyValueGamRepo) FindKeyInGAM(networkID int64, keyName string) (record *model.KeyValueGamModel, err error) {
	err = t.Db.
		Debug().
		Where("network_id = ?", networkID).
		Where("name = ?", keyName).
		Group("name").
		Find(&record).
		Error
	return
}

func (t *keyValueGamRepo) Save(record *model.KeyValueGamModel) (err error) {
	err = t.Db.
		Debug().
		Save(record).
		Error
	return
}
