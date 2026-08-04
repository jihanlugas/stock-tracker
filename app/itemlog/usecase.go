package itemlog

import (
	"fmt"

	"stock-tracker/app/base"
	"stock-tracker/app/item"
	"stock-tracker/constant"
	"stock-tracker/jwt"
	"stock-tracker/model"
	"stock-tracker/request"
	"stock-tracker/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageItemlog) (vItemlogs []model.ItemlogView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vItemlog model.ItemlogView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateItemlog) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateItemlog) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	baseUsecase    base.Usecase
	repository     Repository
	itemRepository item.Repository
}

func (u *usecase) Page(loginUser jwt.UserLogin, req request.PageItemlog) (vItemlogs []model.ItemlogView, count int64, err error) {
	conn := u.baseUsecase.GetConnection()

	vItemlogs, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vItemlogs, count, fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	return vItemlogs, count, nil
}

func (u *usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vItemlog model.ItemlogView, err error) {
	conn := u.baseUsecase.GetConnection()

	vItemlog, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vItemlog, fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	return vItemlog, nil
}

func (u *usecase) Create(loginUser jwt.UserLogin, req request.CreateItemlog) error {
	conn := u.baseUsecase.GetConnection()

	tItemlog := model.Itemlog{
		ID:       utils.GetUniqueID(),
		ItemID:   req.ItemID,
		Type:     req.Type,
		Notes:    req.Notes,
		Quantity: req.Quantity,
		CreateBy: loginUser.UserID,
		UpdateBy: loginUser.UserID,
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	switch constant.ItemlogType(tItemlog.Type) {
	case constant.ITEMLOG_TYPE_SENT:
		err := u.itemRepository.AddSent(tx, tItemlog.ItemID, tItemlog.Quantity)
		if err != nil {
			_ = tx.Rollback().Error
			return fmt.Errorf("failed to add sent: %v", err)
		}
	case constant.ITEMLOG_TYPE_STOCK:
		err := u.itemRepository.AddStock(tx, tItemlog.ItemID, tItemlog.Quantity)
		if err != nil {
			_ = tx.Rollback().Error
			return fmt.Errorf("failed to add stock: %v", err)
		}
	default:
		_ = tx.Rollback().Error
		return fmt.Errorf("invalid type")
	}

	err := u.repository.Create(tx, tItemlog)
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

func (u *usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateItemlog) error {
	conn := u.baseUsecase.GetConnection()

	tItemlog, err := u.repository.GetTableById(conn, id)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	tItemlog.Notes = req.Notes
	tItemlog.UpdateBy = loginUser.UserID

	err = u.repository.Save(tx, tItemlog)
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

	tItemlog, err := u.repository.GetTableById(conn, id)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", u.repository.Name(), err)
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err = u.repository.Delete(tx, tItemlog)
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

func NewUsecase(baseUsecase base.Usecase, repository Repository, itemRepository item.Repository) Usecase {
	return &usecase{
		baseUsecase:    baseUsecase,
		repository:     repository,
		itemRepository: itemRepository,
	}
}
