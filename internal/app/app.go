package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/miladbarzideh/shortify/internal/domain/handler"
	"github.com/miladbarzideh/shortify/internal/domain/service"
)

type Server struct {
	logger *logrus.Logger
}

func NewServer(logger *logrus.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Run() error {
	app := echo.New()
	app.Validator = &CustomValidator{validator: validator.New()}
	s.mapHandlers(app)
	return app.Start(":8080")
}

func (s *Server) mapHandlers(app *echo.Echo) {
	urlService := service.NewService(s.logger)
	urlHandler := handler.NewHandler(s.logger, urlService)
	group := app.Group("/api/v1")
	group.POST("/urls/shorten", urlHandler.CreateShortURL())
	group.GET("/urls/:url", urlHandler.RedirectToLongURL())
}
