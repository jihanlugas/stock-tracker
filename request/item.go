package request

import "time"

type CreateItem struct {
	Name  string `json:"name" form:"name" validate:"required,lte=200"`
	Notes string `json:"notes" form:"notes" validate:""`
	Stock int    `json:"stock" form:"stock" validate:"required,number"`
	Sent  int    `json:"sent" form:"sent" validate:"required,number"`
}

type UpdateItem struct {
	Name  string `json:"name" form:"name" validate:"required,lte=200"`
	Notes string `json:"notes" form:"notes" validate:""`
	Stock int    `json:"stock" form:"stock" validate:"required,number"`
	Sent  int    `json:"sent" form:"sent" validate:"required,number"`
}

type PageItem struct {
	Paging
	Name          string     `json:"name" form:"name" query:"name" search:"name"`
	Notes         string     `json:"notes" form:"notes" query:"notes" search:"notes"`
	CreateName    string     `json:"createName" form:"createName" query:"createName" search:"create_name"`
	Stock         *int       `json:"stock" form:"stock" query:"stock" search:"stock"`
	StartStock    *int       `json:"startStock" form:"startStock" query:"startStock"`
	EndStock      *int       `json:"endStock" form:"endStock" query:"endStock"`
	Sent          *int       `json:"sent" form:"sent" query:"sent" search:"sent"`
	StartSent     *int       `json:"startSent" form:"startSent" query:"startSent"`
	EndSent       *int       `json:"endSent" form:"endSent" query:"endSent"`
	StartCreateDt *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt   *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Search        string     `json:"search" form:"search" query:"search"`
	Preloads      string     `json:"preloads" form:"preloads" query:"preloads"`
}
