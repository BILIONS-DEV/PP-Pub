package model

import (
	"fmt"
	"golang.org/x/oauth2"
	"source/apps/frontend/ggapi"
	"source/apps/frontend/lang"
	"source/core/technology/mysql"
	"source/pkg/utility"
)

type Gam struct{}

type GamRecord struct {
	mysql.TableGam
}

func (GamRecord) TableName() string {
	return mysql.Tables.Gam
}

func (t *Gam) GetById(gamId int64) (rec GamRecord) {
	mysql.Client.Last(&rec, gamId)
	return
}

func (t *Gam) GetByUser(userId int64, email string) (rec GamRecord) {
	mysql.Client.Where("user_id = ? AND gg_email = ?", userId, email).Last(&rec)
	return
}

func (t *Gam) GetByIdOfUser(gamId int64, userId int64) (rec GamRecord) {
	mysql.Client.Where("user_id = ?", userId).Last(&rec, gamId)
	return
}

func (t *Gam) GetByEmail(email string) (rec GamRecord) {
	mysql.Client.Where("gg_email = ?", email).Last(&rec)
	return
}

func (t *Gam) CheckEmail(email string, userId int64) (rec GamRecord) {
	mysql.Client.Where("gg_email = ? and user_id = ?", email, userId).Last(&rec)
	return
}

func (t *Gam) UpdateRefreshTokenByEmail(refreshToken string, email string) {
	mysql.Client.Model(GamRecord{}).Where("gg_email = ?", email).Update("refresh_token", refreshToken)
	return
}

func (t *Gam) Push(userId int64, inputFromGoogle *oauth2.Token, userFromGoogle ggapi.User) (rec GamRecord, err error) {
	// Trường hợp refresh token rỗng là email đã đăng nhập trên tài khoản này hoặc trên tài khoản khác trong hệ thống
	if inputFromGoogle.RefreshToken == "" {
		// Check email đăng nhập trong user này chưa
		rec = t.CheckEmail(userFromGoogle.Email, userId)
		if rec.Id != 0 {
			// Nếu tồn tại rồi thì return luôn
			return
		}
		// Nếu đăng nhập trong user khác
		// Get GAM từ email và tiến hành lấy refresh token
		rec = t.GetByEmail(userFromGoogle.Email)
		// Tiến hành clone gam
		rec.Id = 0
		rec.UserId = userId
		err = mysql.Client.Create(&rec).Error
		if err != nil {
			if !utility.IsWindow() {
				err = fmt.Errorf(lang.Translate.Errors.GamError.Add.ToString())
				return
			}
		}
	} else { // Trường hợp có refresh token tài khoản chưa đăng nhập bao giờ hoặc đã đăng nhập nhưng remove access để đăng nhập lại
		// Kiểm tra xem đã có gam cho email nào chưa
		rec = t.GetByEmail(userFromGoogle.Email)
		if rec.Id != 0 { // Nếu đã có tồn tại ít nhất 1 gam
			// Update lại refresh token cho toàn bộ email
			t.UpdateRefreshTokenByEmail(inputFromGoogle.RefreshToken, userFromGoogle.Email)
		}
		// Check email đăng nhập trong user này chưa
		rec = t.CheckEmail(userFromGoogle.Email, userId)
		if rec.Id != 0 {
			// Nếu tồn tại rồi thì return luôn
			return
		}
		// Nếu chưa tồn tại tiến hành create gam như bình thường
		rec.UserId = userId
		rec.RefreshToken = inputFromGoogle.RefreshToken
		rec.TokenType = inputFromGoogle.Type()
		rec.IsDisabled = mysql.GamEnable
		rec.GgEmail = userFromGoogle.Email
		rec.GgUserId = userFromGoogle.Id
		rec.GgUserRole = userFromGoogle.Role
		err = mysql.Client.Create(&rec).Error
		if err != nil {
			if !utility.IsWindow() {
				err = fmt.Errorf(lang.Translate.Errors.GamError.Add.ToString())
				return
			}
		}
	}
	return
}

func (t *Gam) Callback() {

}
