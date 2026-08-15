package dashboard

import (
	"net/http"
	"stock-tracker/jwt"
	"stock-tracker/response"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

// Dashboard
// @Tags Item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /dashboard [get]
func (h Handler) Dashboard(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	res, err := h.usecase.GetDashboard(loginUser)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully get dashboard", res).SendJSON(c)
}
