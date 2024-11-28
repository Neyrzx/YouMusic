package tracks_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neyrzx/youmusic/internal/domain/entities"
	"github.com/neyrzx/youmusic/internal/domain/repositories/tracks"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestRepositroy(t *testing.T) {
	tests := []struct {
		name  string
		track entities.Track
	}{
		{
			name: "case: create track",
			track: entities.Track{
				Title:    "Track Title",
				Artist:   "Artist Name",
				Link:     "https://linkto.com/asdo-wd02-3c22-d2c2",
				Released: time.Now(),
				Lyrics:   []string{"verse1", "verse2"},
			},
		},
	}

	ctx := context.Background()
	_, connectionString := setupPostgresContainer(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// t.Cleanup(func() {
			// 	require.NoError(t, container.Restore(ctx))
			// })

			db, err := pgxpool.New(ctx, connectionString)
			require.NoError(t, err)

			repo := tracks.NewRepository(db)
			require.NoError(t, repo.Create(ctx, test.track))
		})
	}
}

func setupPostgresContainer(t *testing.T) (*postgres.PostgresContainer, string) {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)

	testcontainers.CleanupContainer(t, container)
	require.NoError(t, err)

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	m, err := migrate.New("file://../../../../migrations", connectionString)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	return container, connectionString
}
