package token_adsense

import (
	"gorm.io/gorm"
	"source/internal/entity/model"
)

const KeyName = "RepoTokenAdsense"

type RepoTokenAdsense interface {
	Migrate()
	Save(record *model.TokenAdsense) (err error)
	FindAll() (records []*model.TokenAdsense, err error)
}

type tokenAdsenseRP struct {
	DB *gorm.DB
}

func NewTokenAdsenseRP(DB *gorm.DB) *tokenAdsenseRP {
	return &tokenAdsenseRP{DB: DB}
}

func (t *tokenAdsenseRP) Migrate() {
	// os.Exit(1)
	err := t.DB.AutoMigrate(
		&model.TokenAdsense{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *tokenAdsenseRP) Save(record *model.TokenAdsense) (err error) {
	err = t.DB.Save(record).Error
	return
}

func (t *tokenAdsenseRP) FindAll() (records []*model.TokenAdsense, err error) {
	err = t.DB.Where(model.TokenAdsense{
		Status: "approved",
	}).Find(&records).Error
	return
}
