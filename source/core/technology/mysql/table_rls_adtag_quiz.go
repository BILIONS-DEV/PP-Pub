package mysql

type TableRlsAdTagQuiz struct {
	Id      int64 `gorm:"column:id" json:"id"`
	AdTagId int64 `gorm:"column:adtag_id" json:"adtag_id"`
	QuizId  int64 `gorm:"column:quiz_id" json:"quiz_id"`
}

func (TableRlsAdTagQuiz) TableName() string {
	return "rls_adtag_quiz"
}
