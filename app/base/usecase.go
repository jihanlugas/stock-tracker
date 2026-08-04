package base

import (
	"stock-tracker/db"

	"gorm.io/gorm"
)

type Usecase interface {
	GetConnection() *gorm.DB
}

type usecase struct{}

func NewUsecase() Usecase {
	return &usecase{}
}

func (u *usecase) GetConnection() *gorm.DB {
	return db.GetPostgresConnection()
}
