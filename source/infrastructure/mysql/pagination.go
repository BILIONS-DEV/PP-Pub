package mysql

import (
	"gorm.io/gorm"
)

type Deps struct {
	Page   int
	Limit  int
	Offset int
}

func Paginate(deps Deps) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Page
		page := deps.Page
		if page < 1 {
			page = 1
		}
		// Limit
		limit := deps.Limit
		if limit == 0 {
			limit = 30
		}
		// Offset
		var offset int
		if deps.Offset > 0 {
			offset = deps.Offset
		} else if deps.Offset == 0 && deps.Page > 1 {
			if page == 1 {
				offset = 0
			} else {
				offset = (page - 1) * limit
			}
		}
		return db.Offset(offset).Limit(limit)
	}
}

//func Pagination(key string, limit int) (page, begin int) {
//	if key == "" {
//		return 1, 0
//	}
//	page, _ = strconv.Atoi(key)
//	if page < 1 {
//		return 1, 0
//	}
//	begin = (limit * page) - limit
//	return
//}
