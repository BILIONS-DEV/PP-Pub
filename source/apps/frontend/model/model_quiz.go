package model

import (
	"source/core/technology/mysql"
)

type Quiz struct{}

type QuizRecord struct {
	mysql.TableQuiz
}

func (t *Quiz) GetById(id int64) (record QuizRecord, err error) {
	record.TableQuiz, err = record.FindById(id)
	return
}

func (t *Quiz) GetByUser(userId int64) (records []QuizRecord, err error) {
	quizs, err := new(QuizRecord).FindByUser(userId)
	for _, quiz := range quizs {
		records = append(records, QuizRecord{quiz})
	}
	return
}

type QizPostsRecord struct {
	mysql.TableQizPosts
}

type QizUsersRecord struct {
	mysql.TableQizUsers
}

// Func GetAllQuizForUserSelect lấy ra toàn bộ quiz từ wp có user là 2 poster(admin) hoặc của user tự tạo
func (t *Quiz) GetAllQuizForUserSelect(email string) (records []QizPostsRecord, err error) {
	var listUserIDGetQuiz []int64
	// Append user poster admin
	listUserIDGetQuiz = append(listUserIDGetQuiz, 2)

	// Get user của DB quiz qua email
	var qizUser QizUsersRecord
	mysql.
		DBQuiz.
		//Debug().
		Where("user_email = ?", email).
		Find(&qizUser)
	if qizUser.ID != 0 {
		listUserIDGetQuiz = append(listUserIDGetQuiz, qizUser.ID)
	}
	// Lấy ra toàn bộ quiz có id là của user và id = 2 là của admin
	err = mysql.
		DBQuiz.
		//Debug().
		Where("post_author in ?", listUserIDGetQuiz).
		Where("post_status = 'publish'").
		Where("post_type = 'advquiz'").
		Find(&records).Error
	return
}
