package cmd

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/infra"
	"github.com/miladbarzideh/shortify/internal/controller"
	"github.com/miladbarzideh/shortify/internal/repository"
	"github.com/miladbarzideh/shortify/internal/service"
	"github.com/miladbarzideh/shortify/pkg/generator"
)

type Server struct {
	logger *logrus.Logger
	cfg    *infra.Config
	db     *gorm.DB
	redis  *redis.Client
}

func NewServer(logger *logrus.Logger, cfg *infra.Config, db *gorm.DB, redis *redis.Client) *Server {
	return &Server{
		logger: logger,
		cfg:    cfg,
		db:     db,
		redis:  redis,
	}
}

func (s *Server) Run() error {
	app := echo.New()
	s.mapHandlers(app)
	return app.Start(fmt.Sprintf(":%s", s.cfg.Server.Port))
}

func (s *Server) mapHandlers(app *echo.Echo) {
	urlRepository := repository.NewRepository(s.logger, s.db)
	urlCacheRepository := repository.NewCacheRepository(s.logger, s.redis)
	gen := generator.NewGenerator()
	urlService := service.NewService(s.logger, s.cfg, urlRepository, urlCacheRepository, gen)
	urlHandler := controller.NewHandler(s.logger, s.cfg, urlService)
	groupV1 := app.Group("/api/v1")
	groupV1.POST("/urls/shorten", urlHandler.CreateShortURL())
	groupV1.GET("/urls/:url", urlHandler.RedirectToLongURL())
}

var cmdServer = func(cfg *infra.Config, log *logrus.Logger, postgresDb *gorm.DB, redis *redis.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the URL shortener app",
		Long: `Start the URL shortener app with a customizable port number.
    Usage example: shortify serve -p 8080`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Changed("port") {
				cfg.Server.Port = cmd.Flag("port").Value.String()
			}

			server := NewServer(log, cfg, postgresDb, redis)
			if server.Run() != nil {
				log.Fatal("failed to start the cmd")
			}
		},
	}
}
