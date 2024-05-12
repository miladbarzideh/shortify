package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/miladbarzideh/shortify/infra"
	"github.com/miladbarzideh/shortify/internal/model"
	"github.com/miladbarzideh/shortify/internal/service"
	"github.com/miladbarzideh/shortify/pkg/generator"
)

type URLHandler interface {
	CreateShortURL() echo.HandlerFunc
	RedirectToLongURL() echo.HandlerFunc
}

type handler struct {
	logger  *logrus.Logger
	cfg     *infra.Config
	service service.URLService
}

func NewHandler(logger *logrus.Logger, cfg *infra.Config, service service.URLService) URLHandler {
	return &handler{
		logger:  logger,
		cfg:     cfg,
		service: service,
	}
}

func (h *handler) CreateShortURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		longURL := new(model.URLData)
		if err := c.Bind(longURL); err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if !longURL.Validate() {
			h.logger.Error("invalid url")
			return echo.NewHTTPError(http.StatusBadRequest, "invalid url")
		}

		shortURL, err := h.service.CreateShortURL(longURL.URL)
		if err != nil {
			h.logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong")
		}

		return c.JSON(http.StatusOK, &model.URLData{
			URL: shortURL,
		})
	}
}

func (h *handler) RedirectToLongURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		shortCode := c.Param("url")
		if !generator.IsValidBase62(shortCode) {
			h.logger.Errorf("Invalid short code: %s", shortCode)
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid short code")
		}

		longURL, err := h.service.GetLongURL(ctx, shortCode)
		if err != nil {
			h.logger.Error(err.Error())
			if errors.Is(err, service.ErrURLNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong")
		}

		return c.Redirect(http.StatusMovedPermanently, longURL)
	}
}
