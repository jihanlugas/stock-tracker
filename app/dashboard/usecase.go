package dashboard

import (
	"fmt"
	"stock-tracker/app/base"
	"stock-tracker/app/item"
	"stock-tracker/jwt"
	"stock-tracker/request"
	"stock-tracker/response"
)

type Usecase interface {
	GetDashboard(loginUser jwt.UserLogin) (res response.Dashboard, err error)
}

type usecase struct {
	baseUsecase    base.Usecase
	itemRepository item.Repository
}

func (u *usecase) GetDashboard(loginUser jwt.UserLogin) (res response.Dashboard, err error) {
	vItems, _, err := u.itemRepository.Page(u.baseUsecase.GetConnection(), request.PageItem{
		Paging: request.Paging{
			Page:  1,
			Limit: 100,
		},
		Preloads: "",
	})
	if err != nil {
		return res, fmt.Errorf("failed to get %s: %v", u.itemRepository.Name(), err)
	}

	res = response.Dashboard{
		Items: vItems,
	}

	return res, nil
}

func NewUsecase(baseUsecase base.Usecase, itemRepository item.Repository) Usecase {
	return &usecase{
		baseUsecase:    baseUsecase,
		itemRepository: itemRepository,
	}
}
