package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (TableQuiz) TableName() string {
	return "quiz"
}

type TableQuiz struct {
	ID           int64               `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title        string              `gorm:"column:title" json:"title"`
	Illustration string              `gorm:"column:illustration"`
	UserID       int64               `gorm:"column:user_id" json:"user_id"`
	CategoryID   int64               `gorm:"column:category_id" json:"category_id"`
	ContentType  int64               `gorm:"column:content_type" json:"content_type"`
	Questions    []TableQuizQuestion `gorm:"many2many:quiz_rls_question;constraint:OnDelete:CASCADE;foreignKey:ID;joinForeignKey:QuizID;References:ID;joinReferences:QuestionID" json:"questions"`
	gorm.Model
}

func (TableQuizQuestion) TableName() string {
	return "quiz_question"
}

type TableQuizQuestion struct {
	ID         int64             `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Question   string            `gorm:"column:question" json:"question"`
	Background string            `gorm:"column:background" json:"background"`
	DeletedAt  gorm.DeletedAt    `gorm:"column:deleted_at" json:"deleted_at"`
	Answers    []TableQuizAnswer `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE" json:"answers"`
}

func (TableQuizAnswer) TableName() string {
	return "quiz_answer"
}

type TableQuizAnswer struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	QuestionID int64          `gorm:"column:question_id" json:"question_id"`
	Answer     string         `gorm:"column:answer" json:"answer"`
	AnswerImg  string         `gorm:"column:answer_img" json:"answer_img"`
	Correct    bool           `gorm:"column:correct" json:"correct"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableQuizRelationshipQuestion) TableName() string {
	return "quiz_rls_question"
}

type TableQuizRelationshipQuestion struct {
	QuizID     int64 `gorm:"column:quiz_id;primaryKey" json:"quiz_id"`
	QuestionID int64 `gorm:"column:question_id;primaryKey" json:"question_id"`
}

func (t *TableQuiz) FindById(id int64) (record TableQuiz, err error) {
	err = Client.Preload(clause.Associations).
		Preload("Questions.Answers").
		Find(&record, id).
		Error
	return
}

func (t *TableQuiz) FindByUser(userId int64) (records []TableQuiz, err error) {
	err = Client.Preload(clause.Associations).
		Preload("Questions.Answers").
		Where("user_id = ?", userId).
		Find(&records).
		Error
	return
}

func (t *TableQuiz) FindByInId(listId []int64) (records []TableQuiz, err error) {
	err = Client.Preload(clause.Associations).
		Preload("Questions.Answers").
		Where("id in ?", listId).
		Find(&records).
		Error
	return
}

func (t *TableQuiz) Create(record TableQuiz) (err error) {
	err = Client.
		Create(&record).
		Error
	return
}

func (t *TableQuiz) Save(record TableQuiz) (err error) {
	err = Client.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&record).
		Error
	return
}
