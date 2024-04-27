package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/miladbarzideh/shortify/internal/domain/service"
)

type Handler struct {
	logger  *logrus.Logger
	service *service.Service
}

func NewHandler(logger *logrus.Logger, service *service.Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) CreateShortURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		longURL := c.FormValue("url")
		shortURL, err := h.service.CreateShortURL(longURL)
		if err != nil {
			h.logger.Error("failed to create short url")
		}

		return c.String(http.StatusOK, shortURL)
	}
}

func (h *Handler) RedirectToLongURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		shortURL := c.Param("url")
		longURL, err := h.service.GetLongURL(shortURL)
		if err != nil {
			h.logger.Errorf("failed to find a relevant url: %s", shortURL)
		}

		return c.String(http.StatusOK, longURL)
	}
}
