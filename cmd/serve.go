package cmd

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/domain"
	"github.com/miladbarzideh/shortify/infra"
)

type Server struct {
	logger *logrus.Logger
	cfg    *infra.Config
	db     *gorm.DB
}

func NewServer(logger *logrus.Logger, cfg *infra.Config, db *gorm.DB) *Server {
	return &Server{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}

func (s *Server) Run() error {
	app := echo.New()
	app.Validator = &CustomValidator{validator: validator.New()}
	s.mapHandlers(app)
	return app.Start(fmt.Sprintf(":%s", s.cfg.Server.Port))
}

func (s *Server) mapHandlers(app *echo.Echo) {
	urlRepository := domain.NewRepository(s.logger, s.db)
	urlService := domain.NewService(s.logger, s.cfg, urlRepository)
	urlHandler := domain.NewHandler(s.logger, s.cfg, urlService)
	groupV1 := app.Group("/api/v1")
	groupV1.POST("/urls/shorten", urlHandler.CreateShortURL())
	groupV1.GET("/urls/:url", urlHandler.RedirectToLongURL())
}

var cmdServer = func(cfg *infra.Config, log *logrus.Logger, postgresDb *gorm.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the URL shortener app",
		Long: `Start the URL shortener app with a customizable port number.
    Usage example: shortify serve -p 8080`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Changed("port") {
				cfg.Server.Port = cmd.Flag("port").Value.String()
			}

			server := NewServer(log, cfg, postgresDb)
			if server.Run() != nil {
				log.Fatal("failed to start the cmd")
			}
		},
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
