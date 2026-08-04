package user

import (
	"fmt"
	"strings"

	"stock-tracker/app/base"
	"stock-tracker/model"
	"stock-tracker/request"
	"stock-tracker/utils"

	"gorm.io/gorm"
)

type Repository interface {
	base.Repository[model.User, model.UserView]
	GetByUsername(conn *gorm.DB, username string, preloads ...string) (tUser model.User, err error)
	GetByEmail(conn *gorm.DB, email string, preloads ...string) (tUser model.User, err error)
	GetByPhoneNumber(conn *gorm.DB, phoneNumber string, preloads ...string) (tUser model.User, err error)
	GetViewByUsername(conn *gorm.DB, username string, preloads ...string) (vUser model.UserView, err error)
	GetViewByEmail(conn *gorm.DB, email string, preloads ...string) (vUser model.UserView, err error)
	GetViewByPhoneNumber(conn *gorm.DB, phoneNumber string, preloads ...string) (vUser model.UserView, err error)
	Page(conn *gorm.DB, req request.PageUser) (vUsers []model.UserView, count int64, err error)
}

type repository struct {
	base.Repository[model.User, model.UserView]
}

func (r repository) GetByUsername(conn *gorm.DB, username string, preloads ...string) (tUser model.User, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("username = ? ", username).First(&tUser).Error
	return tUser, err
}

func (r repository) GetByEmail(conn *gorm.DB, email string, preloads ...string) (tUser model.User, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("email = ? ", email).First(&tUser).Error
	return tUser, err
}

func (r repository) GetByPhoneNumber(conn *gorm.DB, phoneNumber string, preloads ...string) (tUser model.User, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("no_hp = ? ", utils.FormatPhoneTo62(phoneNumber)).First(&tUser).Error
	return tUser, err
}

func (r repository) GetViewByUsername(conn *gorm.DB, username string, preloads ...string) (vUser model.UserView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("username = ? ", username).First(&vUser).Error
	return vUser, err
}

func (r repository) GetViewByEmail(conn *gorm.DB, email string, preloads ...string) (vUser model.UserView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("email = ? ", email).First(&vUser).Error
	return vUser, err
}

func (r repository) GetViewByPhoneNumber(conn *gorm.DB, phoneNumber string, preloads ...string) (vUser model.UserView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("no_hp = ? ", phoneNumber).First(&vUser).Error
	return vUser, err
}

func (r repository) Page(conn *gorm.DB, req request.PageUser) (vUsers []model.UserView, count int64, err error) {
	query := conn.Model(&vUsers)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.Fullname != "" {
		query = query.Where("fullname ILIKE ?", "%"+req.Fullname+"%")
	}
	if req.Email != "" {
		query = query.Where("email ILIKE ?", "%"+req.Email+"%")
	}
	if req.Username != "" {
		query = query.Where("username ILIKE ?", "%"+req.Username+"%")
	}
	if req.PhoneNumber != "" {
		query = query.Where("no_hp ILIKE ?", "%"+utils.FormatPhoneTo62(req.PhoneNumber)+"%")
	}
	if req.Username != "" {
		query = query.Where("username ILIKE ?", "%"+utils.FormatPhoneTo62(req.Username)+"%")
	}
	if req.Address != "" {
		query = query.Where("address ILIKE ?", "%"+utils.FormatPhoneTo62(req.Address)+"%")
	}
	if req.BirthPlace != "" {
		query = query.Where("birth_place ILIKE ?", "%"+utils.FormatPhoneTo62(req.BirthPlace)+"%")
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
		return vUsers, count, err
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

	err = query.Find(&vUsers).Error
	if err != nil {
		return vUsers, count, err
	}

	return vUsers, count, err
}

func NewRepository() Repository {
	return &repository{
		Repository: base.NewRepository[model.User, model.UserView]("user"),
	}
}
