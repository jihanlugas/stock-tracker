package itemlog

import (
	"fmt"
	"stock-tracker/app/base"
	"stock-tracker/model"
	"stock-tracker/request"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	base.Repository[model.Itemlog, model.ItemlogView]
	Page(conn *gorm.DB, req request.PageItemlog) (vItemlogs []model.ItemlogView, count int64, err error)
}

type repository struct {
	base.Repository[model.Itemlog, model.ItemlogView]
}

func (r *repository) Page(conn *gorm.DB, req request.PageItemlog) (vItemlogs []model.ItemlogView, count int64, err error) {
	query := conn.Model(&vItemlogs)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.ItemID != "" {
		query = query.Where("item_id = ?", req.ItemID)
	}
	if req.Notes != "" {
		query = query.Where("notes ILIKE ?", "%"+req.Notes+"%")
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Quantity != nil {
		query = query.Where("quantity = ?", *req.Quantity)
	}
	if req.StartQuantity != nil {
		query = query.Where("quantity >= ?", *req.StartQuantity)
	}
	if req.EndQuantity != nil {
		query = query.Where("quantity <= ?", *req.EndQuantity)
	}
	if req.StartCreateDt != nil {
		query = query.Where("create_dt >= ?", *req.StartCreateDt)
	}
	if req.EndCreateDt != nil {
		query = query.Where("create_dt <= ?", *req.EndCreateDt)
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}
	if req.ItemName != "" {
		query = query.Where("item_name ILIKE ?", "%"+req.ItemName+"%")
	}

	query = base.ApplyGlobalSearch(query, req.Search, req)

	err = query.Count(&count).Error
	if err != nil {
		return vItemlogs, count, err
	}

	sortField := "create_dt"
	sortOrder := "desc"
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

	err = query.Find(&vItemlogs).Error
	if err != nil {
		return vItemlogs, count, err
	}

	return vItemlogs, count, nil
}

func NewRepository() Repository {
	return &repository{
		Repository: base.NewRepository[model.Itemlog, model.ItemlogView]("itemlog"),
	}
}
