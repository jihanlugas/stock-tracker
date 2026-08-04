package item

import (
	"errors"
	"fmt"
	"strings"

	"stock-tracker/app/base"
	"stock-tracker/model"
	"stock-tracker/request"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	base.Repository[model.Item, model.ItemView]
	Page(conn *gorm.DB, req request.PageItem) (vItems []model.ItemView, count int64, err error)
	AddStock(conn *gorm.DB, id string, stock int) error
	AddSent(conn *gorm.DB, id string, sent int) error
}

type repository struct {
	base.Repository[model.Item, model.ItemView]
}

func (r *repository) Page(conn *gorm.DB, req request.PageItem) (vItems []model.ItemView, count int64, err error) {
	query := conn.Model(&vItems)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Notes != "" {
		query = query.Where("notes ILIKE ?", "%"+req.Notes+"%")
	}
	if req.Stock != nil {
		query = query.Where("stock = ?", *req.Stock)
	}
	if req.StartStock != nil {
		query = query.Where("stock >= ?", *req.StartStock)
	}
	if req.EndStock != nil {
		query = query.Where("stock <= ?", *req.EndStock)
	}
	if req.Sent != nil {
		query = query.Where("sent = ?", *req.Sent)
	}
	if req.StartSent != nil {
		query = query.Where("sent >= ?", *req.StartSent)
	}
	if req.EndSent != nil {
		query = query.Where("sent <= ?", *req.EndSent)
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}
	if req.StartCreateDt != nil {
		query = query.Where("create_dt >= ?", req.StartCreateDt)
	}
	if req.EndCreateDt != nil {
		query = query.Where("create_dt <= ?", req.EndCreateDt)
	}

	query = base.ApplyGlobalSearch(query, req.Search, req)

	err = query.Count(&count).Error
	if err != nil {
		return vItems, count, err
	}

	sortField := "create_dt"
	sortOrder := "asc"
	if req.SortField != "" {
		sortField = req.SortField
	}
	if req.SortOrder != "" {
		sortOrder = req.SortOrder
	}
	query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vItems).Error
	if err != nil {
		return vItems, count, err
	}

	return vItems, count, nil
}

func (r *repository) AddStock(conn *gorm.DB, id string, stock int) error {
	if stock < 0 {
		return errors.New("stock value must be non-negative")
	}

	var item model.Item
	err := conn.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&item).Error
	if err != nil {
		return err
	}

	item.Stock += stock
	return conn.Save(&item).Error
}

func (r *repository) AddSent(conn *gorm.DB, id string, sent int) error {
	if sent < 0 {
		return errors.New("sent value must be non-negative")
	}

	var item model.Item
	err := conn.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&item).Error
	if err != nil {
		return err
	}

	if item.Stock < sent {
		return errors.New("insufficient stock")
	}

	item.Sent += sent
	item.Stock -= sent
	return conn.Save(&item).Error
}

func NewRepository() Repository {
	return &repository{
		Repository: base.NewRepository[model.Item, model.ItemView]("item"),
	}
}
