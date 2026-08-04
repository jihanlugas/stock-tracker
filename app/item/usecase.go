package item

import (
	"fmt"
	"time"

	"stock-tracker/app/base"
	"stock-tracker/jwt"
	"stock-tracker/model"
	"stock-tracker/request"
	"stock-tracker/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageItem) (vItems []model.ItemView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vItem model.ItemView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateItem) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateItem) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	baseUsecase base.Usecase
	repository  Repository
}

func (u *usecase) Page(loginUser jwt.UserLogin, req request.PageItem) (vItems []model.ItemView, count int64, err error) {
	conn := u.baseUsecase.GetConnection()

	vItems, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vItems, count, fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	return vItems, count, nil
}

func (u *usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vItem model.ItemView, err error) {
	conn := u.baseUsecase.GetConnection()

	vItem, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vItem, fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	return vItem, nil
}

func (u *usecase) Create(loginUser jwt.UserLogin, req request.CreateItem) error {
	conn := u.baseUsecase.GetConnection()

	tItem := model.Item{
		ID:       utils.GetUniqueID(),
		Name:     req.Name,
		Notes:    req.Notes,
		Stock:    req.Stock,
		Sent:     req.Sent,
		CreateBy: loginUser.UserID,
		UpdateBy: loginUser.UserID,
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := u.repository.Create(tx, tItem)
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

func (u *usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateItem) error {
	conn := u.baseUsecase.GetConnection()

	tItem, err := u.repository.GetTableById(conn, id)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	tItem.Name = req.Name
	tItem.Notes = req.Notes
	tItem.Stock = req.Stock
	tItem.Sent = req.Sent
	tItem.UpdateBy = loginUser.UserID
	tItem.UpdateDt = time.Now()

	err = u.repository.Save(tx, tItem)
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

func (u *usecase) Delete(loginUser jwt.UserLogin, id string) error {
	conn := u.baseUsecase.GetConnection()

	tItem, err := u.repository.GetTableById(conn, id)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err = u.repository.Delete(tx, tItem)
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
