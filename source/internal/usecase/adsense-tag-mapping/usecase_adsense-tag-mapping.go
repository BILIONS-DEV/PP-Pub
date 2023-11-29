package adsense_tag_mapping

import (
	"github.com/asaskevich/govalidator"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	"source/internal/repo/adsense-tag-mapping"
)

type UsecaseAdsenseTagMapping interface {
	Create(record *model.AdsenseTagMappingModel) (err error)
	Update(ID int64, record *model.AdsenseTagMappingModel) (err error)
	Filter(input *InputFilterUC) (totalRecords int64, records []*model.AdsenseTagMappingModel, err error)
	GetByID(id int64) (record *model.AdsenseTagMappingModel, err error)
}

func NewAdsenseTagMappingUC(repos *repo.Repositories) *keyValueUC {
	return &keyValueUC{repos: repos}
}

type keyValueUC struct {
	repos *repo.Repositories
}

type InputFilterUC struct {
	datatable.Request
	UserID      int64  `json:"user_id"`
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
}

func (t *keyValueUC) Filter(input *InputFilterUC) (totalRecords int64, records []*model.AdsenseTagMappingModel, err error) {
	//=> khởi tạo input của Repo
	inputRepo := adsense_tag_mapping.InputFilter{
		UserID: input.UserID,
		Search: input.QuerySearch,
		Offset: input.Start,
		Limit:  input.Length,
		Order:  input.OrderString(),
	}
	return t.repos.AdsenseTagMapping.Filter(&inputRepo)
}

func (t *keyValueUC) validateSubmit(record *model.AdsenseTagMappingModel) (err error) {
	if record.AdsenseAdUnit == 0 {
		err = errors.NewWithID("adsense adunit is required", "adsense_adunit")
		return
	}
	if govalidator.IsNull(string(record.AdUnitType)) {
		err = errors.NewWithID("adunit type is required", "adunit_type")
		return
	}
	if record.TagID == 0 {
		err = errors.NewWithID("tag_id is required", "tag_id")
		return
	}
	if record.BidderID == 0 {
		err = errors.NewWithID("bidder_id is required", "select-bidder")
		return
	}
	return
}

func (t *keyValueUC) Create(record *model.AdsenseTagMappingModel) (err error) {
	//=> bản ghi mới không được có ID, nếu không hàm Save của gorm sẽ hiểu là 1 lệnh update
	if record.IsFound() {
		return errors.New(`ID field is not approved`)
	}

	//=> validate lại các trường dữ liệu trong struct model
	if err = t.validateSubmit(record); err != nil {
		return
	}

	//=> thực hiện insert vào DB bằng hàm Save của gorm
	if err = t.repos.AdsenseTagMapping.Save(record); err != nil {
		return
	}

	//=> Reset cache cho domain của tag
	inventoryAdTag, err := t.repos.InventoryAdTag.FindByID(record.TagID)
	if err != nil {
		return
	}
	if inventoryAdTag.ID == 0 {
		err = errors.NewWithID("not found tag_id", "tag_id")
		return
	}
	if err = t.repos.Inventory.ResetCacheByID(inventoryAdTag.InventoryID); err != nil {
		return
	}
	return
}

func (t *keyValueUC) Update(ID int64, newRecord *model.AdsenseTagMappingModel) (err error) {
	oldRecord, err := t.repos.AdsenseTagMapping.FindByID(ID)
	if err != nil {
		return
	}
	if oldRecord.ID == 0 {
		return errors.NewWithID("not found", "")
	}
	newRecord.ID = oldRecord.ID
	//=> validate lại các trường dữ liệu trong struct model
	if err = t.validateSubmit(newRecord); err != nil {
		return
	}

	//=> thực hiện insert vào DB bằng hàm Save của gorm
	if err = t.repos.AdsenseTagMapping.Save(newRecord); err != nil {
		return
	}

	//=> Reset cache cho domain của tag
	inventoryAdTag, err := t.repos.InventoryAdTag.FindByID(newRecord.TagID)
	if err != nil {
		return
	}
	if inventoryAdTag.ID == 0 {
		err = errors.NewWithID("not found tag_id", "tag_id")
		return
	}
	if err = t.repos.Inventory.ResetCacheByID(inventoryAdTag.InventoryID); err != nil {
		return
	}
	return
}

func (t *keyValueUC) GetByID(id int64) (record *model.AdsenseTagMappingModel, err error) {
	record, err = t.repos.AdsenseTagMapping.FindByID(id)
	return
}
