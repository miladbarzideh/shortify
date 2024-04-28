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
		longURL := new(URL)
		if err := c.Bind(longURL); err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}

		if !longURL.Validate() {
			h.logger.Error("invalid url")
			return echo.NewHTTPError(http.StatusBadRequest, "invalid url")
		}

		shortURL, err := h.service.CreateShortURL(longURL.Url)
		if err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create short url")
		}

		return c.JSON(http.StatusOK, &URL{
			Url: shortURL,
		})
	}
}

func (h *Handler) RedirectToLongURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		shortCode := c.Param("url")
		longURL, err := h.service.GetLongURL(shortCode)
		if err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return c.Redirect(http.StatusMovedPermanently, longURL)
	}
}
