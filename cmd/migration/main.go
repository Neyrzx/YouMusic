package main

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/neyrzx/youmusic/internal/config"
	"github.com/neyrzx/youmusic/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	l := logger.DefaultLogger().With().Str("app", "migrations").Logger()

	if err := godotenv.Load(); err != nil {
		l.Fatal().Err(err).Msg("failed to godotenv.Load")
	}

	if err := migrationUp(ctx); err != nil {
		l.Fatal().Err(err).Msg("failed to migrationUp")
	}

	l.Info().Msg("migrate done successful")
}

func migrationUp(ctx context.Context) error {
	var (
		err    error
		config config.App
		m      *migrate.Migrate
	)

	if err = envconfig.ProcessWith(ctx, &envconfig.Config{Target: &config}); err != nil {
		return fmt.Errorf("failed envconfig.ProcessWith(...): %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s/", config.MigrationsSource)
	connectionURL := config.Database.ConnectionURI()

	if m, err = migrate.New(sourceURL, connectionURL); err != nil {
		return fmt.Errorf("failed migrate.New(%s, %s): %w", sourceURL, connectionURL, err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate.Up() with %s, %s: %w", sourceURL, connectionURL, err)
	}

	sourceErr, databaseErr := m.Close()

	if sourceErr != nil {
		return fmt.Errorf("source error: %w", sourceErr)
	}

	if databaseErr != nil {
		return fmt.Errorf("database error: %w", databaseErr)
	}

	return nil
}
