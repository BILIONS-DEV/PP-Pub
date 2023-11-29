package quiz

import (
	"fmt"
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
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

	//err := t.Db.Debug().AutoMigrate(&model.Quiz{}, &model.QuizQuestion{}, &model.QuizAnswer{}, &model.QuizRelationshipQuestion{})
	//if err != nil {
	//	panic(err)
	//}
	//return
	//
	var quizRec []model.Quiz

	t.Db.Debug().
		Preload("Questions", mysql.Paginate(mysql.Deps{Page: 1})).
		//Preload("Questions", func(db *gorm.DB) *gorm.DB {
		//	//return mysql.Paginate(mysql.Deps{Page: 2})
		//	return db.Offset(2).Limit(2).Order("id DESC")
		//}).
		Preload("Questions.Answers").Offset(1).Limit(1).Find(&quizRec)

	fmt.Printf("\n quizRec ID: %+v \n", quizRec[0].ID)
	fmt.Printf("\n quizRec Title: %+v \n", quizRec[0].Title)
	fmt.Printf("\n total Questions: %+v \n", len(quizRec[0].Questions))
	for _, v := range quizRec[0].Questions {
		fmt.Printf("\n questID: %+v \n", v.ID)
		fmt.Printf("\n QuestionText: %+v \n", v.Question)
		fmt.Println("------------------------------", true)
	}

	//quizRec := model.Quiz{
	//	Title:      "QUIZ: 3",
	//	UserID:     1,
	//	CategoryID: 1,
	//	Questions: []model.QuizQuestion{
	//		{
	//			Question: "QUEST 1 - Quiz 3",
	//			Answers: []model.QuizAnswer{
	//				{
	//					Answer:    "Ans 1 - Quest 1 - Quiz 3",
	//					AnswerImg: "img1.1.3",
	//					Correct:   true,
	//				},
	//				{
	//					Answer:    "Ans 2 - Quest 1 - Quiz 3",
	//					AnswerImg: "img2.1.3",
	//					Correct:   false,
	//				},
	//			},
	//		},
	//		{
	//			Question: "QUEST 2 - Quiz 3",
	//			Answers: []model.QuizAnswer{
	//				{
	//					Answer:    "Ans 1 - Quest 2 - Quiz 3",
	//					AnswerImg: "img1.2.3",
	//					Correct:   true,
	//				},
	//				{
	//					Answer:    "Ans 2 - Quest 2 - Quiz 3",
	//					AnswerImg: "img2.2.3",
	//					Correct:   false,
	//				},
	//			},
	//		},
	//		{
	//			Question: "QUEST 3 - Quiz 3",
	//			Answers: []model.QuizAnswer{
	//				{
	//					Answer:    "Ans 1 - Quest 3 - Quiz 3",
	//					AnswerImg: "img1.3.3",
	//					Correct:   true,
	//				},
	//				{
	//					Answer:    "Ans 2 - Quest 3 - Quiz 3",
	//					AnswerImg: "img2.3.3",
	//					Correct:   false,
	//				},
	//			},
	//		},
	//	},
	//}
	//t.Db.Debug().Save(&quizRec)
	//a := t.Db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Save(&quizRec)
	//a := t.Db.Debug().Select(clause.Associations).Delete(&quizRec)
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

func (t *quiz) FindByUser() (records []model.Quiz) {
	return t.demoQuery()
}

func (t *quiz) Create(record *model.Quiz) (err error) {
	t.Db.
		Debug().
		Create(record)
	return
}

func (t *quiz) Update(record *model.Quiz) (err error) {
	return
}
