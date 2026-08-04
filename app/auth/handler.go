package auth

import (
	"net/http"

	"stock-tracker/jwt"
	"stock-tracker/request"
	"stock-tracker/response"
	"stock-tracker/utils"

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

// SignIn
// @Tags Authentication
// @Accept json
// @Produce json
// @Param req body request.Signin true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/sign-in [post]
func (h Handler) SignIn(c echo.Context) error {
	var err error

	req := new(request.Signin)
	err = c.Bind(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	token, userLogin, err := h.usecase.SignIn(*req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully signed in", response.Payload{
		"token":     token,
		"userLogin": userLogin,
	}).SendJSON(c)
}

// SignOut Sign out user
// @Tags Authentication
// @Accept json
// @Produce json
// // @Param req body request.Signin true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/sign-out [get]
func (h Handler) SignOut(c echo.Context) error {
	return response.Success(http.StatusOK, "Successfully signed out", nil).SendJSON(c)
}

// RefreshToken
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/refresh-token [get]
func (h Handler) RefreshToken(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	token, err := h.usecase.RefreshToken(loginUser)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully refreshed token", response.Payload{
		"token": token,
	}).SendJSON(c)
}

// Init
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/init [get]
func (h Handler) Init(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	user, err := h.usecase.Init(loginUser)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	res := response.Init{
		User: user,
	}

	return response.Success(http.StatusOK, "Successfully initialized", res).SendJSON(c)
}
