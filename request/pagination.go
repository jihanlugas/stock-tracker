package request

import "stock-tracker/config"

type Paging struct {
	Page      int    `json:"page,omitempty" form:"page" query:"page" example:"1"`
	Limit     int    `json:"limit,omitempty" form:"limit" query:"limit" example:"10"`
	SortField string `json:"sortField,omitempty" form:"sortField" query:"sortField" example:""`
	SortOrder string `json:"sortOrder,omitempty" form:"sortOrder" query:"sortOrder" example:""`
}

func (p *Paging) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}

	return p.Page
}

func (p *Paging) GetLimit() int {
	if p.Limit == 0 {
		return config.DefaultDataPerPage
	} else {
		return p.Limit
	}
}

func (p *Paging) SetLimit(lim int) {
	p.Limit = lim
}

func (p *Paging) SetPage(page int) {
	p.Page = page
}

type IPaging interface {
	GetPage() int
	GetLimit() int
	SetLimit(int)
	SetPage(int)
}
