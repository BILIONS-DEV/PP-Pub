package history2

import (
	"source/internal/entity/model"
)

type historyBidder struct {
	*historyUsecase
}

func newHistoryBidder(historyUsecase *historyUsecase) *historyBidder {
	return &historyBidder{historyUsecase: historyUsecase}
}

func (t *historyBidder) push(input *inputPush) (err error) {
	oldRecord := input.oldData.(model.BidderModel)
	newRecord := input.newData.(model.BidderModel)
	var rootRecord model.BidderModel
	if newRecord.ID != 0 {
		rootRecord = newRecord
	} else {
		rootRecord = oldRecord
	}
	record := model.HistoryModel{
		CreatorID: input.creatorID,
		UserID:    0,
		//Title:      input.title,
		Object:     rootRecord.TableName(),
		ObjectID:   rootRecord.ID,
		ObjectName: "",
		ObjectType: model.HistoryObjectTYPE(input.objectType),
		DetailType: input.detailType,
		OldData:    oldRecord.ToJSON(),
		NewData:    newRecord.ToJSON(),
	}
	return t.Repos.History.Save(&record)
}

func (t *historyBidder) validate(push *inputPush) (err error) {
	return
}

func (t *historyBidder) compare(record *model.HistoryModel) (resp []model.CompareData, err error) {
	return
}

func (t *historyBidder) filter() {
	return
}
