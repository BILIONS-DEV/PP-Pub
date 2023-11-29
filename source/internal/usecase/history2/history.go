package history2

import (
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/lang"
	"source/internal/repo"
)

type UsecaseHistory interface {
	Push(input *dto.PayloadHistoryPush) (err error)
	Compare(record *model.HistoryModel) (resp []model.CompareData, err error)
	Filter()
}

type historyImplInterface interface {
	push(push *inputPush) (err error)
	validate(push *inputPush) (err error)
	compare(record *model.HistoryModel) (resp []model.CompareData, err error)
	filter()
}

type historyUsecase struct {
	Repos *repo.Repositories
	Lang  *lang.Translation
}

func NewHistoryUsecase(repos *repo.Repositories, lang *lang.Translation) *historyUsecase {
	return &historyUsecase{Repos: repos, Lang: lang}
}

type inputPush struct {
	creatorID        int64
	title            string
	app              string
	detailType       string
	objectType       model.HistoryObjectTYPE
	oldData, newData interface{}
}

func (t *historyUsecase) Push(input *dto.PayloadHistoryPush) (err error) {
	var impl historyImplInterface
	if impl, err = t.makeImpl(input.DetailType); err != nil {
		return
	}

	objectType := t.makeObjectType(input.OldData, input.NewData)
	inp := &inputPush{
		creatorID:  input.CreatorID,
		title:      t.makeTitle(input.DetailType, objectType),
		detailType: input.DetailType,
		objectType: objectType,
		oldData:    input.OldData,
		newData:    input.NewData,
	}
	if err = impl.validate(inp); err != nil {
		return
	}
	return impl.push(inp)
}

func (t *historyUsecase) makeImpl(detailType string) (impl historyImplInterface, err error) {
	switch detailType {
	case
		string(model.DetailInventorySubmitFE),
		string(model.DetailInventoryConfigFE):
		impl = newHistoryInventory(t)

	case
		string(model.DetailBidderFE),
		string(model.DetailBidderBE):
		impl = newHistoryBidder(t)
	default:
		return nil, errors.New("Malformed")
	}
	return
}

func (t *historyUsecase) makeTitle(detailType string, objectType model.HistoryObjectTYPE) (title string) {
	switch detailType {
	case string(model.DetailInventorySubmitFE):
		title = "Submit Inventory"
		return
	case string(model.DetailInventoryConfigFE):
		title = objectType.String() + "Inventory Config"
		return
	case string(model.DetailInventoryConsentFE):
		title = objectType.String() + "Inventory Consent"
		return
	case string(model.DetailBidderFE):
		return
	case string(model.DetailBidderBE):
		return
	default:
		return ""
	}
}

func (t *historyUsecase) makeObjectType(oldData, newData interface{}) (objectType model.HistoryObjectTYPE) {
	if oldData == nil && newData != nil {
		return model.HistoryObjectTYPEAdd // Create
	}
	if oldData != nil && newData != nil {
		return model.HistoryObjectTYPEUpdate // Update
	}
	if oldData != nil && newData == nil {
		return model.HistoryObjectTYPEDelete // Delete
	}
	return
}

func (t *historyUsecase) Compare(record *model.HistoryModel) (resp []model.CompareData, err error) {
	var impl historyImplInterface
	if impl, err = t.makeImpl(record.DetailType); err != nil {
		return
	}
	resp, err = impl.compare(record)
	return
}

func (t *historyUsecase) Filter() {
	return
}

func makeResponseCompare(text string, recOldData *string, recNewData *string, objectType model.HistoryObjectTYPE) (res model.CompareData, err error) {
	if recOldData == nil && recNewData == nil {
		return res, errors.New("no response")
	}
	var oldData, newData string
	if recOldData != nil {
		oldData = *recOldData
	} else {
		oldData = ""
	}
	if recNewData != nil {
		newData = *recNewData
	} else {
		newData = ""
	}
	var action string
	switch objectType {
	case model.HistoryObjectTYPEAdd:
		action = "add"
		break
	case model.HistoryObjectTYPEUpdate:
		action = compareStringData(recOldData, recNewData)
	case model.HistoryObjectTYPEDelete:
		action = "delete"
	}
	res = model.CompareData{
		Action:  action,
		Text:    text,
		OldData: oldData,
		NewData: newData,
	}
	return
}

func compareStringData(oldData *string, newData *string) (action string) {
	if oldData == nil && newData != nil {
		action = "add"
		return
	}
	if oldData != nil && newData == nil {
		action = "delete"
		return
	}

	if oldData != nil && newData != nil {
		if *oldData != *newData {
			action = "update"
			return
		}
		if *oldData == *newData {
			action = "none"
			return
		}
	}

	return
}
