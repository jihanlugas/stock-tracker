package auth

import (
	"errors"
	"fmt"
	"time"

	"stock-tracker/app/base"
	"stock-tracker/app/user"
	"stock-tracker/config"
	"stock-tracker/cryption"
	"stock-tracker/jwt"
	"stock-tracker/model"
	"stock-tracker/request"
	"stock-tracker/utils"
)

type Usecase interface {
	SignIn(req request.Signin) (token string, userLogin jwt.UserLogin, err error)
	RefreshToken(userLogin jwt.UserLogin) (token string, err error)
	Init(userLogin jwt.UserLogin) (vUser model.UserView, err error)
}

type usecase struct {
	baseUsecase    base.Usecase
	userRepository user.Repository
}

func (u usecase) SignIn(req request.Signin) (token string, userLogin jwt.UserLogin, err error) {

	var tUser model.User

	conn := u.baseUsecase.GetConnection()

	if utils.IsValidEmail(req.Username) {
		tUser, err = u.userRepository.GetByEmail(conn, req.Username)
	} else {
		tUser, err = u.userRepository.GetByUsername(conn, req.Username)
	}

	if err != nil {
		return "", userLogin, err
	}

	err = cryption.CheckAES64(req.Passwd, tUser.Passwd)
	if err != nil {
		return "", userLogin, errors.New("invalid username or password")
	}

	if !tUser.IsActive {
		return "", userLogin, errors.New("user not active")
	}

	now := time.Now()
	tx := conn.Begin()

	tUser.LastLoginDt = &now
	tUser.UpdateBy = tUser.ID
	err = u.userRepository.Update(tx, model.User{
		ID:          tUser.ID,
		LastLoginDt: &now,
		UpdateBy:    tUser.ID,
	})
	if err != nil {
		return "", userLogin, err
	}

	err = tx.Commit().Error
	if err != nil {
		return "", userLogin, fmt.Errorf("failed to commit transaction: %w", err)
	}

	expiredAt := time.Now().Add(time.Minute * time.Duration(config.AuthTokenExpiredMinute))
	userLogin.ExpiredDt = expiredAt
	userLogin.UserID = tUser.ID
	userLogin.Role = tUser.Role
	userLogin.PassVersion = tUser.PassVersion
	token, err = jwt.CreateToken(userLogin)
	if err != nil {
		return "", userLogin, err
	}

	return token, userLogin, err
}

func (u usecase) RefreshToken(userLogin jwt.UserLogin) (token string, err error) {
	userLogin.ExpiredDt = time.Now().Add(time.Minute * time.Duration(config.AuthTokenExpiredMinute))

	token, err = jwt.CreateToken(userLogin)
	if err != nil {
		return token, err
	}

	return token, err
}

func (u usecase) Init(userLogin jwt.UserLogin) (vUser model.UserView, err error) {
	conn := u.baseUsecase.GetConnection()

	vUser, err = u.userRepository.GetViewById(conn, userLogin.UserID)
	if err != nil {
		return vUser, err
	}

	return vUser, err
}

func NewUsecase(baseUsecase base.Usecase, userRepository user.Repository) Usecase {
	return &usecase{
		baseUsecase:    baseUsecase,
		userRepository: userRepository,
	}
}
