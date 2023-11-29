package quiz

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/dto"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"
)

type quiz struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewQuizRepo(db *gorm.DB, cache caching.Cache) *quiz {
	return &quiz{Db: db, Cache: cache}
}

func (t *quiz) insert() {
	var records model.Quiz
	for i := 0; i < 10; i++ {
		iString := strconv.Itoa(i)
		records = model.Quiz{
			ID:          3,
			Title:       "Quiz " + iString,
			UserID:      4,
			CategoryID:  5,
			ContentType: 2,
			Questions: []model.QuizQuestion{
				{
					ID:       6,
					Question: "Quiz abab" + iString + " - Quest 1",
					Answers: []model.QuizAnswer{
						{
							ID:        11,
							Answer:    "Quiz ggg" + iString + " - Quest 1 - Answer 1",
							AnswerImg: "IMG ->>> Quiz " + iString + " - Quest 1 - Answer 1",
							Correct:   false,
						},
						{
							ID:        12,
							Answer:    "Quiz 1 - Quest 1 - Answer 2",
							AnswerImg: "IMG ->>> Quiz " + iString + " - Quest 1 - Answer 2",
							Correct:   false,
						},
					},
				},
				{
					ID:       7,
					Question: "Quiz tttt" + iString + " - Quest 2",
					Answers: []model.QuizAnswer{
						{
							ID:        13,
							Answer:    "Quiz " + iString + " - Quest 2 - Answer 1",
							AnswerImg: "IMG ->>> Quiz " + iString + " - Quest 2 - Answer 1",
							Correct:   false,
						},
						{
							ID:        14,
							Answer:    "Quiz 1 - Quest 2 - Answer 2",
							AnswerImg: "IMG ->>> Quiz " + iString + " - Quest 2 - Answer 2",
							Correct:   false,
						},
					},
				},
			},
		}
		//records = append(records, rec)
	}
	fmt.Println(records)
	t.Db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&records)
}

func (t *quiz) Migrate() {
	err := t.Db.Debug().AutoMigrate(&model.Quiz{}, &model.QuizQuestion{}, &model.QuizAnswer{}, &model.QuizRelationshipQuestion{})
	if err != nil {
		panic(err)
	}
}

func (t *quiz) FindByUser(userID int64) (records []model.Quiz) {
	var quizRec []model.Quiz
	t.Db.Debug().
		Preload("QuizRelationshipQuestion", func(db *gorm.DB) *gorm.DB {
			return db.Or("quiz_id DESC")
		}).
		Preload("Questions", mysql.Paginate(mysql.Deps{Page: 1})).
		//Preload("Questions", func(db *gorm.DB) *gorm.DB {
		//	//return mysql.Paginate(mysql.Deps{Page: 2})
		//	return db.Offset(2).Limit(2).Order("id DESC")
		//}).
		//Preload("Questions.Answers").
		Where("user_id = ?", userID).
		Offset(0).Limit(1).
		Find(&quizRec)
	return
}

func (t *quiz) Create(record *model.Quiz) (err error) {
	err = t.Db.
		//Debug().
		Create(record).
		Error
	return
}

func (t *quiz) Update(record *model.Quiz) (err error) {
	t.Db.
		Unscoped().
		Select(clause.Associations, "Questions").Delete(&record.Questions)
	err = t.Db.
		//Debug().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(record).
		Error
	return
}

func (t *quiz) FindById(id int64) (record model.Quiz, err error) {
	err = t.Db.Preload(clause.Associations).
		Preload("Questions.Answers").
		Find(&record, id).
		Error
	return
}

func (t *quiz) FindByFilters(payload *dto.QuizFilterPayload, userLogin model.User) (response datatable.Response, err error) {
	var records []model.Quiz
	var total int64
	err = t.Db.Where("user_id = ?", userLogin.ID).
		Scopes(
			t.setFilterSearch(payload),
		).
		Model(&records).Count(&total).
		Scopes(
			t.setOrder(payload),
			pagination.Paginate(pagination.Params{
				Limit:  payload.Length,
				Offset: payload.Start,
			}),
		).Find(&records).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf("error")
		}
		return datatable.Response{}, err
	}
	response.Draw = payload.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.makeResponseDatatable(records)
	return
}

func (t *quiz) setFilterSearch(payload *dto.QuizFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool
		// Search from form of datatable <- not use
		if payload.Search != nil && payload.Search.Value != "" {
			flag = true
		}
		// Search from form filter
		if payload.PostData.QuerySearch != "" {
			flag = true
		}
		if !flag {
			return db
		}
		return db.Where("title LIKE ?", "%"+payload.PostData.QuerySearch+"%")
	}
}

func (t *quiz) setOrder(payload *dto.QuizFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(payload.Order) > 0 {
			var orders []string
			for _, order := range payload.Order {
				column := payload.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db.Order("id DESC")
	}
}

type ResponseQuizDatatable struct {
	Title       string `json:"title"`
	ContentType string `json:"content_type"`
	Action      string `json:"action"`
}

func (t *quiz) makeResponseDatatable(quizs []model.Quiz) (records []ResponseQuizDatatable) {
	for _, quiz := range quizs {
		rec := ResponseQuizDatatable{
			Title:       quiz.Title,
			ContentType: quiz.ContentType.String(),
			Action:      htmlblock.Render("quiz/index/block.action.gohtml", quiz).String(),
		}
		records = append(records, rec)
	}
	return
}

func (t *quiz) Delete(record *model.Quiz) (response dto.ResponseDelete, err error) {
	err = t.Db.Delete(record).Error
	if err != nil {
		response = dto.ResponseDelete{
			ID:      record.ID,
			Status:  false,
			Message: "error",
		}
	}

	response = dto.ResponseDelete{
		ID:      record.ID,
		Status:  true,
		Message: "done",
	}
	return
}
