package item

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
// @Tags Item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req query request.PageItem false "url query string"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /item [get]
func (h Handler) Page(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.PageItem)
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

	return response.Success(http.StatusOK, "Successfully retrieved item list", response.PayloadPagination(req, data, count)).SendJSON(c)
}

// GetById
// @Tags Item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /item/{id} [get]
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

	vItem, err := h.usecase.GetById(loginUser, id, preloadSlice...)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, "Successfully retrieved item detail", vItem).SendJSON(c)
}

// Create
// @Tags Item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateItem true "json req body"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /item [post]
func (h Handler) Create(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.CreateItem)
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

	return response.Success(http.StatusOK, "Successfully created item", nil).SendJSON(c)
}

// Update
// @Tags Item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param req body request.UpdateItem true "json req body"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /item/{id} [put]
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

	req := new(request.UpdateItem)
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

	return response.Success(http.StatusOK, "Successfully updated item", nil).SendJSON(c)
}

// Delete
// @Tags Item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object} response.Response
// @Failure      500  {object} response.Response
// @Router /item/{id} [delete]
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

	return response.Success(http.StatusOK, "Successfully deleted item", nil).SendJSON(c)
}
