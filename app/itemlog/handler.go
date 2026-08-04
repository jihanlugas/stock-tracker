package itemlog

import (
	"net/http"
	"strings"

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

// Page
// @Tags Itemlog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req query request.PageItemlog false "url query string"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /itemlog [get]
func (h Handler) Page(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.PageItemlog)
	err = c.Bind(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	data, count, err := h.usecase.Page(loginUser, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully retrieved itemlog list", response.PayloadPagination(req, data, count)).SendJSON(c)
}

// GetById
// @Tags Itemlog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /itemlog/{id} [get]
func (h Handler) GetById(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	preloadSlice := []string{}
	preloads := c.QueryParam("preloads")
	if preloads != "" {
		preloadSlice = strings.Split(preloads, ",")
	}

	vItemlog, err := h.usecase.GetById(loginUser, id, preloadSlice...)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully retrieved itemlog detail", vItemlog).SendJSON(c)
}

// Create
// @Tags Itemlog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateItemlog true "json req body"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /itemlog [post]
func (h Handler) Create(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.CreateItemlog)
	err = c.Bind(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	err = h.usecase.Create(loginUser, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully created itemlog", nil).SendJSON(c)
}

// Update
// @Tags Itemlog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param req body request.UpdateItemlog true "json req body"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /itemlog/{id} [put]
func (h Handler) Update(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	req := new(request.UpdateItemlog)
	err = c.Bind(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	err = h.usecase.Update(loginUser, id, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully updated itemlog", nil).SendJSON(c)
}

// Delete
// @Tags Itemlog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /itemlog/{id} [delete]
func (h Handler) Delete(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	err = h.usecase.Delete(loginUser, id)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully deleted itemlog", nil).SendJSON(c)
}
