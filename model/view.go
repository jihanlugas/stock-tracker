package model

import (
	"time"

	"gorm.io/gorm"
)

type UserView struct {
	ID                string         `json:"id"`
	Role              string         `json:"role"`
	Email             string         `json:"email"`
	Username          string         `json:"username"`
	PhoneNumber       string         `json:"phoneNumber"`
	Address           string         `json:"address"`
	Fullname          string         `json:"fullname"`
	Passwd            string         `json:"-"`
	PassVersion       int            `json:"passVersion"`
	IsActive          bool           `json:"isActive"`
	PhotoID           string         `json:"photoId"`
	PhotoUrl          string         `json:"photoUrl"`
	LastLoginDt       *time.Time     `json:"lastLoginDt"`
	BirthDt           *time.Time     `json:"birthDt"`
	BirthPlace        string         `json:"birthPlace"`
	AccountVerifiedDt *time.Time     `json:"accountVerifiedDt"`
	CreateBy          string         `json:"createBy"`
	CreateDt          time.Time      `json:"createDt"`
	UpdateBy          string         `json:"updateBy"`
	UpdateDt          time.Time      `json:"updateDt"`
	DeleteDt          gorm.DeletedAt `json:"deleteDt"`
	CreateName        string         `json:"createName"`
	UpdateName        string         `json:"updateName"`
}

func (UserView) TableName() string {
	return VIEW_USER
}

type ItemView struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Notes      string         `json:"notes"`
	Stock      int            `json:"stock"`
	Sent       int            `json:"sent"`
	CreateBy   string         `json:"createBy"`
	CreateDt   time.Time      `json:"createDt"`
	UpdateBy   string         `json:"updateBy"`
	UpdateDt   time.Time      `json:"updateDt"`
	DeleteDt   gorm.DeletedAt `json:"deleteDt"`
	CreateName string         `json:"createName"`
	UpdateName string         `json:"updateName"`
}

func (ItemView) TableName() string {
	return VIEW_ITEM
}

type ItemlogView struct {
	ID         string         `json:"id"`
	ItemID     string         `json:"itemId"`
	Type       string         `json:"type"` // "STOCK" | "SENT"
	Notes      string         `json:"notes"`
	Quantity   int            `json:"quantity"`
	CreateBy   string         `json:"createBy"`
	CreateDt   time.Time      `json:"createDt"`
	UpdateBy   string         `json:"updateBy"`
	UpdateDt   time.Time      `json:"updateDt"`
	DeleteDt   gorm.DeletedAt `json:"deleteDt"`
	CreateName string         `json:"createName"`
	UpdateName string         `json:"updateName"`
	ItemName   string         `json:"itemName"`

	Item *ItemView `json:"item,omitempty"`
}

func (ItemlogView) TableName() string {
	return VIEW_ITEMLOG
}
