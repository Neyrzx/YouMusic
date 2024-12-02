package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/neyrzx/youmusic/docs"
	"github.com/neyrzx/youmusic/internal/config"
	"github.com/neyrzx/youmusic/internal/delivery/rest"
	"github.com/neyrzx/youmusic/internal/domain/repositories"
	"github.com/neyrzx/youmusic/internal/domain/services"
	"github.com/neyrzx/youmusic/internal/gateways"
	"github.com/neyrzx/youmusic/pkg/httpclient"
	"github.com/neyrzx/youmusic/pkg/logger"
	"github.com/neyrzx/youmusic/pkg/validator"
	"github.com/sethvargo/go-envconfig"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	ctx := context.Background()

	mode := config.GetCurrentRunningMode()

	l := logger.DefaultLogger()
	l.With().Str("mode", string(mode))

	if mode == config.ModeLocal {
		if err := godotenv.Load(); err != nil {
			l.Error().Err(err).Msg("failed to godotenv.Load")
		}
	}

	var cfg config.App
	if err := envconfig.ProcessWith(ctx, &envconfig.Config{Target: &cfg}); err != nil {
		l.Error().Err(err).Msg("failed to envconfig.ProcessWith")
	}

	e := echo.New()

	valid, err := validator.NewValidator()
	if err != nil {
		l.Error().Err(err).Msg("failed to validator.NewValidator")
	}
	e.Validator = valid

	// Middlewares
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			l.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("method", v.Method).
				Dur("latency", v.Latency).
				Str("ip", v.RemoteIP).
				Msg("request")
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Dependencies
	db, err := pgxpool.New(ctx, cfg.Database.ConnectionURI())
	if err != nil {
		l.Error().Err(err).Msg("failed to pgxpool.New")
	}

	client := httpclient.NewHTTPClient(cfg.GatewayMusicInfo)
	tracksRepository := repositories.NewTracksRepository(db)
	musicInfoGateway := gateways.NewMusicInfoGateway(client, cfg.GatewayMusicInfo)
	tracksService := services.NewTracksService(tracksRepository, musicInfoGateway)

	// Routes
	rest.InitAPI(e, tracksService)
	e.GET(cfg.SwaggerDocPath, echoSwagger.WrapHandler)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		if err = e.Start(cfg.Server.ServerAddr); err != nil && errors.Is(err, http.ErrServerClosed) {
			l.Error().Err(err).Msg("failed to e.Start")
		}
	}()

	<-ctx.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.Server.GracefulShoutdownTimeout)
	defer cancelFunc()

	if err = e.Shutdown(ctx); err != nil {
		l.Error().Err(err).Msg("failed to e.Shutdown")
	}

	l.Info().Msg("server successfuly shutdown")
}
