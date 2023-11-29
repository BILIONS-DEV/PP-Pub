package section

import (
	"gorm.io/gorm"
	"source/internal/entity/model"
)

const KeyName = "RepoSection"

type RepoSection interface {
	Save(record *model.Section) (err error)
	FindBySectionID(sectionID, traffic_source string) (record model.Section, err error)
	FindBySectionName(sectionName, traffic_source string) (record model.Section, err error)
}

type sectionRP struct {
	DB *gorm.DB
}

func NewSectionRP(DB *gorm.DB) *sectionRP {
	return &sectionRP{DB: DB}
}

func (t *sectionRP) Save(record *model.Section) (err error) {
	err = t.DB.Save(record).Error
	return
}

func (t *sectionRP) FindBySectionID(sectionID, traffic_source string) (record model.Section, err error) {
	err = t.DB.
		Where("section_id = ? AND traffic_source = ?", sectionID, traffic_source).
		Find(&record).Error
	return
}

func (t *sectionRP) FindBySectionName(sectionName, traffic_source string) (record model.Section, err error) {
	err = t.DB.
		Where("section_name = ? AND traffic_source = ?", sectionName, traffic_source).
		Find(&record).Error
	return
}
