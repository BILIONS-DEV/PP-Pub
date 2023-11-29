package _example

import (
	"gorm.io/gorm"
	"source/internal/entity/model"
)

var db *gorm.DB

func Find() {
	var records []model.Quiz
	db.
		Debug().
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
			return db.Select("quiz_answer.question_id").Order("id DESC")
		}).
		Find(&records)

	fmt.Printf("\n records: %+v \n", records)
	return
}

func Update() {
	quizRec := model.Quiz{
		ID:         1,
		Title:      "this is quiz 1.1.1.2",
		UserID:     1,
		CategoryID: 1,
		Questions: []model.QuizQuestion{
			{
				ID:       2,
				Question: "question 1 for QUIZ 2.2.2.2",
				Answers: []model.QuizAnswer{
					{
						ID:        1,
						Answer:    "ans 1 of question 1.1.2.3",
						AnswerImg: "img1.1.2",
						Correct:   true,
					},
					{
						ID:        2,
						Answer:    "ans 2 of question 1.1.2.3",
						AnswerImg: "",
						Correct:   false,
					},
				},
			},
		},
	}
	db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Save(&quizRec)
	//db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Create(&quizRec)
	//db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Updates(&quizRec)
}
