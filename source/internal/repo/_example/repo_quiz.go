package quiz

import (
	"fmt"
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
)

type quiz struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewQuizRepo(db *gorm.DB, cache caching.Cache) *quiz {
	return &quiz{Db: db, Cache: cache}
}

func (t *quiz) demoQuery() (records []model.Quiz) {

	//quizRec := model.Quiz{}
	//t.Db.Debug().Preload("Questions").Preload("Questions.Answers").Find(&quizRec, 1)
	//
	//fmt.Printf("\n quizRec: %+v \n", quizRec)

	//err := t.Db.Debug().AutoMigrate(&model.Quiz{}, &model.QuizQuestion{}, &model.QuizAnswer{}, &model.QuizRelationshipQuestion{})
	//if err != nil {
	//	panic(err)
	//}
	//return

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
	//t.Db.Debug().Create(&quizRec)
	a := t.Db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Save(&quizRec)
	fmt.Printf("\n a: %+v \n", a)
	return
}

func (t *quiz) demoQuery2() (records []model.Quiz) {
	t.Db.
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

func (t *quiz) demoQuery3() (records []model.Quiz) {
	//err := t.Db.SetupJoinTable(&model.Quiz{}, "Questions", &model.QuizRelationshipQuestion{})
	//if err != nil {
	//	panic(err)
	//}
	t.Db.
		Debug().
		Preload("Questions").
		//Preload("Questions", func(db *gorm.DB) *gorm.DB {
		//	return db.Order("id ASC")
		//}).
		//Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
		//	return db.Select("quiz_answer.question_id").Order("id DESC")
		//}).
		Find(&records)
	fmt.Printf("\n records: %+v \n", records)
	return
}

func (t *quiz) FindByUser2() (records []model.Quiz) {
	return t.demoQuery()
}

func (t *quiz) Creates(record *model.Quiz) (err error) {
	t.Db.
		Debug().
		Create(record)
	return
}

func (t *quiz) Updates(record *model.Quiz) (err error) {
	return
}
