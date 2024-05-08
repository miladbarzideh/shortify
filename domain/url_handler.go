package domain

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/miladbarzideh/shortify/domain/generator"
	"github.com/miladbarzideh/shortify/infra"
)

type URLData struct {
	URL string `json:"url" validate:"required,url"`
}

type Handler struct {
	logger  *logrus.Logger
	cfg     *infra.Config
	service *Service
}

func NewHandler(logger *logrus.Logger, cfg *infra.Config, service *Service) *Handler {
	return &Handler{
		logger:  logger,
		cfg:     cfg,
		service: service,
	}
}

func (h *Handler) CreateShortURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		longURL := new(URLData)
		if err := c.Bind(longURL); err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(longURL); err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		shortURL, err := h.service.CreateShortURL(longURL.URL)
		if err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &URLData{
			URL: shortURL,
		})
	}
}

func (h *Handler) RedirectToLongURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		shortCode := c.Param("url")
		if !generator.IsValidBase62(shortCode) {
			h.logger.Errorf("Invalid short code: %s", shortCode)
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid short code")
		}

		longURL, err := h.service.GetLongURL(shortCode)
		if err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return c.Redirect(http.StatusMovedPermanently, longURL)
	}
}
