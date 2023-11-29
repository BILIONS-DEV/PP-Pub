package dto

type PayloadHistoryPush struct {
	DetailType string      `json:"detail_type"`
	CreatorID  int64       `json:"creator_id"`
	OldData    interface{} `json:"old_data"`
	NewData    interface{} `json:"new_data"`
}

type PayloadHistoryFilter struct {
	ObjectPage string `json:"object_page"`
	ObjectID   string `json:"object_id"`
	Search     string `json:"search"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Page       int    `json:"page"`
}

func (t *PayloadHistoryFilter) Validate() (err error) {
	return
}

func (t *PayloadHistoryFilter) ToCondition() map[string]interface{} {
	var condition = make(map[string]interface{})
	if t.ObjectPage != "" {
		condition["page"] = t.ObjectPage
	}
	if t.ObjectID != "" {
		condition["object_id"] = t.ObjectID
	}
	return condition
}

type PayloadGetHistory struct {
	Id     string  `json:"id" form:"id"`
	Object string `json:"object" form:"object"`
}

func (t *PayloadGetHistory) Validate() (err error) {
	return
}
