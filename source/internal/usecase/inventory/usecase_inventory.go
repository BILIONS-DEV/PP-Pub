package inventory

import (
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/repo"
	inventoryRepo "source/internal/repo/inventory"
)

type UsecaseInventory interface {
	Filter(input *InputFilterUC) (totalRecords int64, records []*model.InventoryModel, err error)
	Submit(domain string) (err error)
	GetById(id int64) (record *model.InventoryModel, err error)
	CountConnectWaiting(record *model.InventoryModel) (count int)
	GetByPublisher(userID int64) (records []model.InventoryModel, err error)
	GetAdTagByID(adTagID int64) (record *model.InventoryAdTagModel, err error)
}

func NewInventoryUC(repos *repo.Repositories) *inventoryUC {
	return &inventoryUC{repos: repos}
}

type inventoryUC struct {
	repos *repo.Repositories
}

type InputFilterUC struct {
	datatable.Request
	UserID      int64       `json:"user_id"`
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	AdsSync     interface{} `query:"f_ads_sync[]" json:"f_ads_sync" form:"f_ads_sync[]"`
}

func (t *inventoryUC) Filter(input *InputFilterUC) (totalRecords int64, records []*model.InventoryModel, err error) {
	// => khởi tạo input của Repo
	inputRepo := inventoryRepo.InputFilter{
		UserID:     input.UserID,
		Status:     input.Status,
		SupplyType: input.Type,
		SyncAdsTxt: input.AdsSync,
		Search:     input.QuerySearch,
		Offset:     input.Start,
		Limit:      input.Length,
		Order:      input.OrderString(),
	}
	return t.repos.Inventory.Filter(&inputRepo)
}

func (t *inventoryUC) Submit(domain string) (err error) {
	return
}

func (t *inventoryUC) GetById(id int64) (record *model.InventoryModel, err error) {
	record, err = t.repos.Inventory.FindByID(id)
	return
}

func (t *inventoryUC) CountConnectWaiting(record *model.InventoryModel) (count int) {
	// => Lấy ra toàn bộ RlsBidderSystem hiện tại
	for _, rlsBidderSystem := range record.RlsBidderSystem {
		// => Từ rls tìm ra bidder, kiểm tra loại bỏ bidder là type là google
		bidder, _ := t.repos.Bidder.FindByID(rlsBidderSystem.BidderID)
		if bidder.ID == 0 || bidder.BidderTemplateId == 1 { // Bỏ qua nếu k tìm thấy bidder hoặc là bidder google
			continue
		}
		// status := new(model.InventoryConnectionDemand).GetStatus(row.Id, bidder.Id)
		// if status == 2 {
		//	assigns.CountConnectWaiting++
		// }
	}
	return
}

func (t *inventoryUC) GetByPublisher(userID int64) (records []model.InventoryModel, err error) {
	return t.repos.Inventory.GetByPublisher(userID)
}

func (t *inventoryUC) GetAdTagByID(adTagID int64) (record *model.InventoryAdTagModel, err error) {
	return t.repos.InventoryAdTag.FindByID(adTagID)
}
