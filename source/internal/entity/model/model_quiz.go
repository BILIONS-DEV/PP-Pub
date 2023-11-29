package model

import "gorm.io/gorm"

func (Quiz) TableName() string {
	return "quiz"
}

type Quiz struct {
	ID           int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title        string          `gorm:"column:title" json:"title"`
	Illustration string          `gorm:"column:illustration"`
	UserID       int64           `gorm:"column:user_id" json:"user_id"`
	CategoryID   int64           `gorm:"column:category_id" json:"category_id"`
	ContentType  TYPEContentType `gorm:"column:content_type" json:"content_type"`
	Questions    []QuizQuestion  `gorm:"many2many:quiz_rls_question;constraint:OnDelete:CASCADE;foreignKey:ID;joinForeignKey:QuizID;References:ID;joinReferences:QuestionID" json:"questions"`
	gorm.Model
}

func (QuizQuestion) TableName() string {
	return "quiz_question"
}

type QuizQuestion struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Question   string         `gorm:"column:question" json:"question"`
	Background string         `gorm:"column:background" json:"background"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	Answers    []QuizAnswer   `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE" json:"answers"`
}

func (QuizAnswer) TableName() string {
	return "quiz_answer"
}

type QuizAnswer struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	QuestionID int64          `gorm:"column:question_id" json:"question_id"`
	Answer     string         `gorm:"column:answer" json:"answer"`
	AnswerImg  string         `gorm:"column:answer_img" json:"answer_img"`
	Correct    bool           `gorm:"column:correct" json:"correct"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (QuizRelationshipQuestion) TableName() string {
	return "quiz_rls_question"
}

type QuizRelationshipQuestion struct {
	QuizID     int64 `gorm:"column:quiz_id;primaryKey" json:"quiz_id"`
	QuestionID int64 `gorm:"column:question_id;primaryKey" json:"question_id"`
}

type TYPEContentType int

const (
	TYPEContentTypeRelated TYPEContentType = iota + 1
	TYPEContentTypeQuiz
)

func (t TYPEContentType) Int() int64 {
	switch t {
	case TYPEContentTypeRelated:
		return 1
	case TYPEContentTypeQuiz:
		return 2
	default:
		return 0
	}
}

func (t TYPEContentType) String() string {
	switch t {
	case TYPEContentTypeRelated:
		return "Related"
	case TYPEContentTypeQuiz:
		return "Quiz 1"
	default:
		return ""
	}
}
