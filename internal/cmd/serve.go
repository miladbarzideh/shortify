package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/domain/controller"
	"github.com/miladbarzideh/shortify/internal/domain/repository"
	"github.com/miladbarzideh/shortify/internal/domain/service"
	"github.com/miladbarzideh/shortify/internal/infra"
	"github.com/miladbarzideh/shortify/pkg/generator"
)

type Server struct {
	logger    *logrus.Logger
	cfg       *infra.Config
	db        *gorm.DB
	redis     *redis.Client
	telemetry *infra.TelemetryProvider
}

func NewServer(
	logger *logrus.Logger,
	cfg *infra.Config,
	db *gorm.DB,
	redis *redis.Client,
	telemetry *infra.TelemetryProvider,
) *Server {
	return &Server{
		logger:    logger,
		cfg:       cfg,
		db:        db,
		redis:     redis,
		telemetry: telemetry,
	}
}

func (s *Server) Run() {
	app := echo.New()
	s.mapHandlers(app)
	// https://echo.labstack.com/docs/cookbook/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		address := fmt.Sprintf(":%s", s.cfg.Server.Port)
		if err := app.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("shutting down the server", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) mapHandlers(app *echo.Echo) {
	urlRepository := repository.NewRepository(s.logger, s.db, s.telemetry)
	urlCacheRepository := repository.NewCacheRepository(s.logger, s.redis, s.telemetry)
	gen := generator.NewGenerator(s.cfg.Shortener.CodeLength)
	urlService := service.NewService(s.logger, s.cfg, urlRepository, urlCacheRepository, gen, s.telemetry)
	urlHandler := controller.NewHandler(s.logger, s.cfg, urlService, s.telemetry)
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

			telemetry, err := infra.NewTelemetry(log, cfg)
			if err != nil {
				log.Fatal(err)
			}

			server := NewServer(log, cfg, postgresDb, redis, telemetry)
			server.Run()
		},
	}
}
