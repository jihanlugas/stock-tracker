package model

import "time"

const (
	VIEW_USER    = "users_view"
	VIEW_ITEM    = "items_view"
	VIEW_ITEMLOG = "itemlogs_view"
)

type UserLogin struct {
	ExpiredDt     time.Time `json:"expiredDt"`
	UserID        string    `json:"userId"`
	PassVersion   int       `json:"passVersion"`
	CompanyID     string    `json:"companyId"`
	Role          string    `json:"role"`
	UsercompanyID string    `json:"usercompanyId"`
}
