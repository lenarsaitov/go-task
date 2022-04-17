package cards

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go-task/pkg/bind"
	"net/http"
	"strconv"
)

type CardHandler struct {
	service *CardService
}

func NewCardHandler(service *CardService) *CardHandler {
	return &CardHandler{service}
}

type StatusResponse string

var OK StatusResponse = "OK"
var Error StatusResponse = "Error"

type Response struct {
	Status  StatusResponse `json:"status"`
	Message string         `json:"message,omitempty"`
}

var INDENT = "  "

func (h *CardHandler) Setup(root *echo.Group) {
	g := root.Group("/cards")

	g.GET("", h.ListCards)
	g.GET("/:id", h.CardItem)

	g.POST("", h.AddCardItem)

	g.PUT("/:id", h.UpdateCardItem)
	g.DELETE("/:id", h.DeleteCardItem)

	g.POST("/transfer", h.TransferAmount)
	g.POST("/:id", h.RefillBalance)
}

func (h *CardHandler) CardItem(c echo.Context) error {
	cardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.HandleError(c, http.StatusBadRequest, err)
	}

	p, err := h.service.GetCard(c.Request().Context(), cardID)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == nil {
		return c.JSONPretty(http.StatusNotFound, Response{Status: OK, Message: "Not Found"}, INDENT)
	}

	return c.JSONPretty(http.StatusOK, p, INDENT)
}

func (h *CardHandler) ListCards(c echo.Context) error {
	var sizeInt, pageInt, userIDInt int
	var err error

	pageStr := c.FormValue("page")
	sizeStr := c.FormValue("size")
	userID := c.FormValue("user_id")

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

	if len(userID) != 0 {
		userIDInt, err = strconv.Atoi(userID)
		if err != nil {
			return h.HandleError(c, http.StatusBadRequest, err)
		}
	}

	params := &FilterParams{UserID: userIDInt, Page: pageInt, Size: sizeInt}
	p, err := h.service.GetListCards(c.Request().Context(), params)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}

	return c.JSONPretty(http.StatusOK, p, INDENT)
}

func (h *CardHandler) AddCardItem(c echo.Context) error {
	params := &AddCardRequestParams{}
	err := bind.DecodeJSONBody(c, params)
	if err != nil {
		var mr *bind.MalformedRequest
		if errors.As(err, &mr) {
			return h.HandleError(c, mr.Status, err)
		}
		return h.HandleError(c, http.StatusInternalServerError, err)
	}

	p, err := h.service.AddCard(c.Request().Context(), params)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == nil {
		return c.JSON(http.StatusNotFound, Response{Status: Error, Message: "User Not Found"})
	}

	return h.HandleSuccess(c, http.StatusOK, fmt.Sprintf("Successfully added. Card ID: %d", *p))
}

func (h *CardHandler) UpdateCardItem(c echo.Context) error {
	cardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.HandleError(c, http.StatusBadRequest, err)
	}

	params := &UpdateCardRequestParams{CardID: cardID}
	err = bind.DecodeJSONBody(c, params)
	if err != nil {
		var mr *bind.MalformedRequest
		if errors.As(err, &mr) {
			return h.HandleError(c, mr.Status, err)
		}
		return h.HandleError(c, http.StatusInternalServerError, err)
	}

	p, err := h.service.UpdateCard(c.Request().Context(), params)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == false {
		return c.JSON(http.StatusNotFound, Response{Status: Error, Message: "Not Found"})
	}

	return h.HandleSuccess(c, http.StatusOK, "")
}

func (h *CardHandler) DeleteCardItem(c echo.Context) error {
	cardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.HandleError(c, http.StatusBadRequest, err)
	}

	p, err := h.service.DeleteCard(c.Request().Context(), cardID)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == false {
		return c.JSON(http.StatusNotFound, Response{Status: Error, Message: "Not Found"})
	}

	return h.HandleSuccess(c, http.StatusOK, "")
}

func (h *CardHandler) RefillBalance(c echo.Context) error {
	cardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.HandleError(c, http.StatusBadRequest, err)
	}

	params := &RefillCardRequestParams{CardID: cardID}
	err = bind.DecodeJSONBody(c, params)
	if err != nil {
		var mr *bind.MalformedRequest
		if errors.As(err, &mr) {
			return h.HandleError(c, mr.Status, err)
		}
		return h.HandleError(c, http.StatusInternalServerError, err)
	}

	if params.AddBalance < 0 {
		return h.HandleError(c, http.StatusBadRequest, errors.New("Can add balance only. Not subtract"))
	}

	p, err := h.service.RefillCard(c.Request().Context(), params)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if p == false {
		return c.JSON(http.StatusNotFound, Response{Status: Error, Message: "Not Found"})
	}

	return h.HandleSuccess(c, http.StatusOK, "")
}

func (h *CardHandler) TransferAmount(c echo.Context) error {
	params := &TransferBalanceCardRequestParams{}
	err := bind.DecodeJSONBody(c, params)
	if err != nil {
		var mr *bind.MalformedRequest
		if errors.As(err, &mr) {
			return h.HandleError(c, mr.Status, err)
		}
		return h.HandleError(c, http.StatusInternalServerError, err)
	}

	switch {
	case params.AddBalance < 0:
		return h.HandleError(c, http.StatusBadRequest, errors.New("Can add balance only. Not subtract"))
	case params.CardTo < 0:
		return h.HandleError(c, http.StatusBadRequest, errors.New("Invalid value of cardID to"))
	case params.CardFrom < 0:
		return h.HandleError(c, http.StatusBadRequest, errors.New("Invalid value of cardID from"))
	}

	exist, err := h.service.TransferBalanceCard(c.Request().Context(), params)
	if err != nil {
		return h.HandleError(c, http.StatusInternalServerError, err)
	}
	if exist == false {
		return c.JSON(http.StatusNotFound, Response{Status: Error, Message: "Not Found"})
	}

	return h.HandleSuccess(c, http.StatusOK, "")
}

func (h *CardHandler) HandleSuccess(c echo.Context, statusCode int, Message string) error {
	return c.JSONPretty(statusCode, Response{Status: OK, Message: Message}, INDENT)
}

func (h *CardHandler) HandleError(c echo.Context, statusCode int, err error) error {
	return c.JSONPretty(statusCode, Response{Status: Error, Message: fmt.Sprintf("%s", err)}, INDENT)
}
