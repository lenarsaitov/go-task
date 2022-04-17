package users

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lenarsaitov/go-task/pkg/bind"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service}
}

type StatusResponse string

var OK StatusResponse = "OK"
var Error StatusResponse = "Error"

type Response struct {
	Status  StatusResponse `json:"status"`
	Message string         `json:"message,omitempty"`
}

var INDENT = "  "

func (h *UserHandler) Setup(root *echo.Group) {
	g := root.Group("/users")

	g.GET("", h.ListUsers)
	g.GET("/:id", h.UserItem)

	g.POST("", h.AddUserItem)

	g.PUT("/:id", h.UpdateUserItem)
	g.DELETE("/:id", h.DeleteUserItem)
}

func (h *UserHandler) UserItem(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.HandleError(c, http.StatusBadRequest, err)
	}

	p, err := h.service.GetUser(c.Request().Context(), userID)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == nil {
		return c.JSON(http.StatusNotFound, Response{Status: OK, Message: "Not Found"})
	}

	return c.JSONPretty(http.StatusOK, p, INDENT)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	var sizeInt, pageInt int
	var err error
	pageStr := c.FormValue("page")
	sizeStr := c.FormValue("size")
	userName := c.FormValue("user_name")

	if len(pageStr) != 0 {
		pageInt, err = strconv.Atoi(pageStr)
		if err != nil {
			return h.HandleError(c, http.StatusBadRequest, err)
		}
	}
	if len(sizeStr) != 0 {
		sizeInt, err = strconv.Atoi(sizeStr)
		if err != nil {
			return h.HandleError(c, http.StatusBadRequest, err)
		}
	}

	params := &FilterParams{UserName: userName, Page: pageInt, Size: sizeInt}

	p, err := h.service.GetListUsers(c.Request().Context(), params)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}

	return c.JSONPretty(http.StatusOK, p, INDENT)
}

func (h *UserHandler) AddUserItem(c echo.Context) error {
	params := &AddUserRequestParams{}
	err := bind.DecodeJSONBody(c, params)
	if err != nil {
		var mr *bind.MalformedRequest
		if errors.As(err, &mr) {
			return h.HandleError(c, mr.Status, err)
		}
		return h.HandleError(c, http.StatusInternalServerError, err)
	}

	p, err := h.service.AddUser(c.Request().Context(), params)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == nil {
		return c.JSON(http.StatusNotFound, Response{Status: "Error", Message: "Not Found"})
	}

	return h.HandleSuccess(c, http.StatusOK, fmt.Sprintf("Successfully added. User ID: %d", *p))
}

func (h *UserHandler) UpdateUserItem(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.HandleError(c, http.StatusBadRequest, err)
	}

	params := &UpdateUserRequestParams{UserID: userID}
	err = bind.DecodeJSONBody(c, params)
	if err != nil {
		var mr *bind.MalformedRequest
		if errors.As(err, &mr) {
			return h.HandleError(c, mr.Status, err)
		}
		return h.HandleError(c, http.StatusInternalServerError, err)
	}

	p, err := h.service.UpdateUser(c.Request().Context(), params)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == false {
		return c.JSON(http.StatusNotFound, Response{Status: "Error", Message: "Not Found"})
	}

	return h.HandleSuccess(c, http.StatusOK, "")
}

func (h *UserHandler) DeleteUserItem(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.HandleError(c, http.StatusBadRequest, err)
	}

	p, err := h.service.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == false {
		return c.JSON(http.StatusNotFound, Response{Status: "Error", Message: "Not Found"})
	}

	return h.HandleSuccess(c, http.StatusOK, "")
}

func (h *UserHandler) HandleSuccess(c echo.Context, statusCode int, Message string) error {
	return c.JSONPretty(statusCode, Response{Status: OK, Message: Message}, INDENT)
}

func (h *UserHandler) HandleError(c echo.Context, statusCode int, err error) error {
	return c.JSONPretty(statusCode, Response{Status: Error, Message: fmt.Sprintf("%s", err)}, INDENT)
}
