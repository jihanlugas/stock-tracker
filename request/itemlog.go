package request

import "time"

type CreateItemlog struct {
	ItemID   string `json:"itemId" form:"itemId" query:"itemId" validate:"required"`
	Notes    string `json:"notes" form:"notes" query:"notes" validate:""`
	Quantity int    `json:"quantity" form:"quantity" query:"quantity" validate:"required,number"`
	Type     string `json:"type" form:"type" query:"type" validate:"required"`
}

type UpdateItemlog struct {
	Notes string `json:"notes" form:"notes" validate:""`
}

type PageItemlog struct {
	Paging
	ItemID        string     `json:"itemId" form:"itemId" query:"itemId"`
	Notes         string     `json:"notes" form:"notes" query:"notes" search:"notes"`
	Type          string     `json:"type" form:"type" query:"type"`
	Quantity      *int       `json:"quantity" form:"quantity" query:"quantity" search:"quantity"`
	StartQuantity *int       `json:"startQuantity" form:"startQuantity" query:"startQuantity"`
	EndQuantity   *int       `json:"endQuantity" form:"endQuantity" query:"endQuantity"`
	StartCreateDt *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt   *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	CreateName    string     `json:"createName" form:"createName" query:"createName" search:"create_name"`
	ItemName      string     `json:"itemName" form:"itemName" query:"itemName" search:"item_name"`
	Search        string     `json:"search" form:"search" query:"search"`
	Preloads      string     `json:"preloads" form:"preloads" query:"preloads"`
}
