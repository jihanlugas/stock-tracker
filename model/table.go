package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                string         `gorm:"primaryKey"`
	Role              string         `gorm:"not null"`
	Email             string         `gorm:"not null"`
	Username          string         `gorm:"not null"`
	PhoneNumber       string         `gorm:"not null"`
	Address           string         `gorm:"not null"`
	Fullname          string         `gorm:"not null"`
	Passwd            string         `gorm:"not null"`
	PassVersion       int            `gorm:"not null"`
	IsActive          bool           `gorm:"not null"`
	PhotoID           string         `gorm:"not null"`
	LastLoginDt       *time.Time     `gorm:"null"`
	BirthDt           *time.Time     `gorm:"null"`
	BirthPlace        string         `gorm:"not null"`
	AccountVerifiedDt *time.Time     `gorm:"null"`
	CreateBy          string         `gorm:"not null"`
	CreateDt          time.Time      `gorm:"not null"`
	UpdateBy          string         `gorm:"not null"`
	UpdateDt          time.Time      `gorm:"not null"`
	DeleteDt          gorm.DeletedAt `gorm:"null"`
}

type Item struct {
	ID       string         `gorm:"primaryKey"`
	Name     string         `gorm:"not null"`
	Notes    string         `gorm:"not null"`
	Stock    int            `gorm:"not null"`
	Sent     int            `gorm:"not null"`
	CreateBy string         `gorm:"not null"`
	CreateDt time.Time      `gorm:"not null"`
	UpdateBy string         `gorm:"not null"`
	UpdateDt time.Time      `gorm:"not null"`
	DeleteDt gorm.DeletedAt `gorm:"null"`
}

type Itemlog struct {
	ID       string         `gorm:"primaryKey"`
	ItemID   string         `gorm:"not null"`
	Type     string         `gorm:"not null"` // "STOCK" | "SENT"
	Notes    string         `gorm:"not null"`
	Quantity int            `gorm:"not null"`
	CreateBy string         `gorm:"not null"`
	CreateDt time.Time      `gorm:"not null"`
	UpdateBy string         `gorm:"not null"`
	UpdateDt time.Time      `gorm:"not null"`
	DeleteDt gorm.DeletedAt `gorm:"null"`
}
