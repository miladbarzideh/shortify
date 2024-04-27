package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger *logrus.Logger
}

func NewHandler(logger *logrus.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) CreateShortURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "CreateShortURL")
	}
}

func (h *Handler) GetShortURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "GetShortURL")
	}
}
