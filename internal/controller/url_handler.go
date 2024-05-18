package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/miladbarzideh/shortify/infra"
	"github.com/miladbarzideh/shortify/internal/model"
	"github.com/miladbarzideh/shortify/internal/service"
	"github.com/miladbarzideh/shortify/pkg/generator"
)

const (
	msgInvalidURLError       = "invalid URL"
	msgInvalidShortCodeError = "invalid short code"
	msgInternalServerError   = "internal server error"
)

type URLHandler interface {
	CreateShortURL() echo.HandlerFunc
	RedirectToLongURL() echo.HandlerFunc
}

type handler struct {
	logger         *logrus.Logger
	cfg            *infra.Config
	service        service.URLService
	tracer         trace.Tracer
	getReqCount    infra.Counter
	createReqCount infra.Counter
}

func NewHandler(logger *logrus.Logger, cfg *infra.Config, service service.URLService, telemetry *infra.TelemetryProvider) URLHandler {
	tracer := telemetry.TraceProvider.Tracer("urlHandler")
	meter := telemetry.MeterProvider.Meter("urlHandler")
	getReqCount := infra.NewCounter(meter, "url.gets")
	createReqCount := infra.NewCounter(meter, "url.creates")

	return &handler{
		logger:         logger,
		cfg:            cfg,
		service:        service,
		tracer:         tracer,
		getReqCount:    getReqCount,
		createReqCount: createReqCount,
	}
}

func (h *handler) CreateShortURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := h.tracer.Start(c.Request().Context(), "urlHandler.create")
		defer span.End()
		longURL := new(model.URLData)
		if err := c.Bind(longURL); err != nil {
			h.logger.Error(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if !longURL.Validate() {
			h.logger.Error(msgInvalidURLError)
			span.RecordError(errors.New(msgInvalidURLError))
			span.SetStatus(codes.Error, msgInvalidURLError)
			return echo.NewHTTPError(http.StatusBadRequest, msgInvalidURLError)
		}

		span.SetAttributes(attribute.String("url", longURL.URL))
		shortURL, err := h.service.CreateShortURL(ctx, longURL.URL)
		if err != nil {
			h.logger.Error(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, msgInternalServerError)
		}

		h.createReqCount.Inc(ctx)

		return c.JSON(http.StatusOK, &model.URLData{
			URL: shortURL,
		})
	}
}

func (h *handler) RedirectToLongURL() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := h.tracer.Start(c.Request().Context(), "urlHandler.redirect")
		defer span.End()
		shortCode := c.Param("url")
		if !generator.IsValidBase62(shortCode) {
			h.logger.Errorf("%s: %s", msgInvalidShortCodeError, shortCode)
			span.RecordError(errors.New(msgInvalidShortCodeError))
			span.SetStatus(codes.Error, msgInvalidShortCodeError)
			return echo.NewHTTPError(http.StatusBadRequest, msgInvalidShortCodeError)
		}

		longURL, err := h.service.GetLongURL(ctx, shortCode)
		if err != nil {
			h.logger.Error(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			if errors.Is(err, service.ErrURLNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			return echo.NewHTTPError(http.StatusInternalServerError, msgInternalServerError)
		}

		h.getReqCount.Inc(ctx)

		return c.Redirect(http.StatusMovedPermanently, longURL)
	}
}
