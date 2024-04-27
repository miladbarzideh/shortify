package app

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/miladbarzideh/shortify/internal/domain/handler"
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
	s.mapHandlers(app)
	return app.Start(":8080")
}

func (s *Server) mapHandlers(app *echo.Echo) {

	urlHandler := handler.NewHandler(s.logger)

	// Map routes
	group := app.Group("/api/v1")
	group.POST("/shorten", urlHandler.CreateShortURL())
	group.GET(":url", urlHandler.GetShortURL())
}
