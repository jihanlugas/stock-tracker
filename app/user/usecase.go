package user

import (
	"errors"
	"fmt"
	"time"

	"stock-tracker/app/base"
	"stock-tracker/constant"
	"stock-tracker/cryption"
	"stock-tracker/jwt"
	"stock-tracker/model"
	"stock-tracker/request"
	"stock-tracker/utils"

	"gorm.io/gorm"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageUser) (vUsers []model.UserView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vUser model.UserView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateUser) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateUser) error
	ChangePassword(loginUser jwt.UserLogin, req request.ChangePassword) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	baseUsecase base.Usecase
	repository  Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageUser) (vUsers []model.UserView, count int64, err error) {
	conn := u.baseUsecase.GetConnection()

	vUsers, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vUsers, count, err
	}

	return vUsers, count, nil
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vUser model.UserView, err error) {
	conn := u.baseUsecase.GetConnection()

	vUser, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vUser, fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	return vUser, nil
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateUser) error {
	conn := u.baseUsecase.GetConnection()

	now := time.Now()

	encodePasswd, err := cryption.EncryptAES64(req.Passwd)
	if err != nil {
		return fmt.Errorf("failed to encode password: %v", err)
	}

	tUser := model.User{
		ID:                utils.GetUniqueID(),
		Role:              constant.RoleUser,
		Email:             req.Email,
		Username:          req.Username,
		PhoneNumber:       utils.FormatPhoneTo62(req.PhoneNumber),
		Address:           req.Address,
		Fullname:          req.Fullname,
		Passwd:            encodePasswd,
		PassVersion:       1,
		IsActive:          true,
		PhotoID:           "",
		LastLoginDt:       nil,
		BirthDt:           req.BirthDt,
		BirthPlace:        req.BirthPlace,
		AccountVerifiedDt: &now,
		CreateBy:          loginUser.UserID,
		UpdateBy:          loginUser.UserID,
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err = u.repository.Create(tx, tUser)
	if err != nil {
		_ = tx.Rollback().Error
		return fmt.Errorf("failed to create %s: %v", u.repository.Name(), err)
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateUser) error {
	conn := u.baseUsecase.GetConnection()

	tUser, err := u.repository.GetTableById(conn, id)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if tUser.Email != req.Email {
		_, err = u.repository.GetByEmail(tx, req.Email)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to update %s: %v", u.repository.Name(), err)
			}
		} else {
			return errors.New("email already exist")
		}
	}

	if tUser.PhoneNumber != utils.FormatPhoneTo62(req.PhoneNumber) {
		_, err = u.repository.GetByPhoneNumber(tx, utils.FormatPhoneTo62(req.PhoneNumber))
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to update %s: %v", u.repository.Name(), err)
			}
		} else {
			return errors.New("phone number already exist")
		}
	}

	tUser.Fullname = req.Fullname
	tUser.Email = req.Email
	tUser.PhoneNumber = utils.FormatPhoneTo62(req.PhoneNumber)
	tUser.Username = req.Username
	tUser.Address = req.Address
	tUser.BirthDt = req.BirthDt
	tUser.BirthPlace = req.BirthPlace
	tUser.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tUser)
	if err != nil {
		_ = tx.Rollback().Error
		return fmt.Errorf("failed to update %s: %v", u.repository.Name(), err)
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (u usecase) ChangePassword(loginUser jwt.UserLogin, req request.ChangePassword) error {
	conn := u.baseUsecase.GetConnection()

	tUser, err := u.repository.GetTableById(conn, loginUser.UserID)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err = cryption.CheckAES64(req.CurrentPasswd, tUser.Passwd)
	if err != nil {
		_ = tx.Rollback().Error
		return fmt.Errorf("invalid current password")
	}

	encodePasswd, err := cryption.EncryptAES64(req.Passwd)
	if err != nil {
		_ = tx.Rollback().Error
		return fmt.Errorf("failed to encode password: %v", err)
	}

	tUser.Passwd = encodePasswd
	tUser.PassVersion += 1
	tUser.UpdateBy = loginUser.UserID
	err = u.repository.Save(tx, tUser)
	if err != nil {
		_ = tx.Rollback().Error
		return fmt.Errorf("failed to update password: %v", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	conn := u.baseUsecase.GetConnection()

	tUser, err := u.repository.GetTableById(conn, id)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err = u.repository.Delete(tx, tUser)
	if err != nil {
		_ = tx.Rollback().Error
		return fmt.Errorf("failed to delete %s: %v", u.repository.Name(), err)
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func NewUsecase(baseUsecase base.Usecase, repository Repository) Usecase {
	return &usecase{
		baseUsecase: baseUsecase,
		repository:  repository,
	}
}
