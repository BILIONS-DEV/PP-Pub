package quiz

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
	"source/internal/entity/model"
)

// https://techmaster.vn/posts/36796/golang-gorm-many-to-many-relationship <- xem ví dụ

func (m *Member) TableName() string {
	return "tungdt_member"
}

func (c *Club) TableName() string {
	return "tungdt_club"
}

func (mc *MemberClub) TableName() string {
	return "tungdt_member_club"
}

type Member struct {
	Id    string `gorm:"primaryKey"`
	Name  string
	Clubs []Club `gorm:"many2many:tungdt_member_club;"` // thể hiện mối quan hệ nhiều nhiều với bảng member_club
}

type Club struct {
	Id   string `gorm:"primaryKey"`
	Name string
	//Members []Member `gorm:"many2many:member_club;"` // thể hiện mối quan hệ nhiều nhiều với bảng member_club
}

type MemberClub struct {
	MemberId string `gorm:"primaryKey" column:"member_id"`
	ClubId   string `gorm:"primaryKey" column:"club_id"`
	Active   bool
}

func (t *quiz) GetMemberByName(name string) (member Member, err error) {
	err = t.Db.Debug().Preload("Clubs", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Find(&member, "name = ?", name).Error
	if err != nil {
		return Member{}, err
	}
	return member, nil
}

func (t *quiz) AddMemberToClub2() (err error) {
	member := Member{
		Id:   "quacquac",
		Name: "TungDT",
		Clubs: []Club{{
			Id:   "hiphop",
			Name: "CandyCrew HoaBinh",
		}},
	}
	//t.Db.Debug().Create(&member)
	t.Db.Debug().Select("Clubs").Delete(&member)
	return
}

func (t *quiz) demoQueryTest() (records []model.Quiz) {
	_ = t.AddMemberToClub2()
	return

	mem, err := t.GetMemberByName("TungDT")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n mem: %+v \n", mem)

	//err := t.AddMemberToClub()
	//if err != nil {
	//	return nil
	//}
	//return
	//club := Club{}
	//member := Member{}
	//err := t.Db.Debug().
	//	//Preload("Members").Find(&club, "name = ?", "Sport").Error
	//	Preload("Clubs").Find(&member, "name = ?", "Bob").Error
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("\n club: %+v \n", club)
	//fmt.Printf("\n member: %+v \n", member)

	return
}

func (t *quiz) createTable() {
	err := t.Db.Debug().AutoMigrate(&Member{}, &Club{}, &MemberClub{})
	if err != nil {
		panic(err)
	}
}

func NewID(length ...int) (id string) {
	id, _ = gonanoid.New(8)
	return
}

func (t *quiz) AddMemberToClub() (err error) {
	// Khởi tạo transaction
	tx := t.Db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	//---- Tạo members
	john := Member{
		Id:   NewID(),
		Name: "John",
	}

	anna := Member{
		Id:   NewID(),
		Name: "Anna",
	}

	bob := Member{
		Id:   NewID(),
		Name: "Bob",
	}

	alice := Member{
		Id:   NewID(),
		Name: "Alice",
	}

	// Thêm các đối tượng vừa tạo vào mảng members
	var members []Member
	members = append(members, john, anna, bob, alice)

	// Insert các bản ghi member vào trong CSDL
	if err := tx.Create(&members).Error; err != nil {
		tx.Rollback()
		return err
	}

	//--- Club
	math := Club{
		Id:   NewID(),
		Name: "Math",
	}

	sport := Club{
		Id:   NewID(),
		Name: "Sport",
	}

	music := Club{
		Id:   NewID(),
		Name: "Music",
	}

	// Thêm các đối tượng vừa tạo vào mảng clubs
	var clubs []Club
	clubs = append(clubs, math, sport, music)

	// Insert các bản ghi club vào trong CSDL
	if err := tx.Create(&clubs).Error; err != nil {
		tx.Rollback()
		return err
	}

	//---- Thêm các thành viên vào club
	err = assignMembersToClub(tx, math, []Member{john, anna})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = assignMembersToClub(tx, sport, []Member{bob, alice})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = assignMembersToClub(tx, music, []Member{john, bob, alice})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Function thực hiện gắn các bản ghi member vào 1 club cụ thể
func assignMembersToClub(tx *gorm.DB, club Club, members []Member) (err error) {
	for _, member := range members {
		err := tx.Create(&MemberClub{
			MemberId: member.Id,
			ClubId:   club.Id,
			Active:   true, //random true or false
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
