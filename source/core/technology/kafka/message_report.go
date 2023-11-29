package kafka

import "time"

type ReportMessage struct {
	Key          string
	UserId       int
	UserUid      string
	AdSetId      int
	Date         time.Time
	Domain       string
	Device       string
	Browser      string
	Location     string
	IsValidImp   bool
	IsValidClick bool
	IsImp        bool
	IsClick      bool
	Expenses     float64
	Ip           string
	Ua           string
}

func (this *ReportMessage) GetTopicName() string {
	return Topics.Report
}
