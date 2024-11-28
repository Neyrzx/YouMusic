package config_test

import (
	"context"
	"testing"

	"github.com/neyrzx/youmusic/internal/config"
	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/assert"
)

// TODO: Покрыть тесткейсами
func TestDatabase_ConnectionURI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		config   *config.Database
		lookuper envconfig.Lookuper
		want     string
	}{
		{
			name: "valid",
			lookuper: envconfig.MapLookuper(map[string]string{
				"DB_NAME":     "youmusic",
				"DB_USER":     "app",
				"DB_PASSWORD": "supersecret",
			}),
			want: "postgresql://app:supersecret@localhost:5432/youmusic?connect_timeout=10&sslmode=disable",
		},
		{
			name: "no-password",
			lookuper: envconfig.MapLookuper(map[string]string{
				"DB_NAME": "youmusic",
				"DB_USER": "app",
			}),
			want: "postgresql://app@localhost:5432/youmusic?connect_timeout=10&sslmode=disable",
		},
		{
			name: "has-password_no-user",
			lookuper: envconfig.MapLookuper(map[string]string{
				"DB_NAME":     "youmusic",
				"DB_PASSWORD": "supersecret",
			}),
			want: "postgresql://localhost:5432/youmusic?connect_timeout=10&sslmode=disable",
		},
		{
			name: "no-password_no-user",
			lookuper: envconfig.MapLookuper(map[string]string{
				"DB_NAME":     "youmusic",
				"DB_PASSWORD": "supersecret",
			}),
			want: "postgresql://localhost:5432/youmusic?connect_timeout=10&sslmode=disable",
		},
		{
			name: "no-port",
			lookuper: envconfig.MapLookuper(map[string]string{
				"DB_NAME":     "youmusic",
				"DB_PASSWORD": "supersecret",
				"DB_PORT":     "",
			}),
			want: "postgresql://localhost/youmusic?connect_timeout=10&sslmode=disable",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var config config.Database

			envconfig.ProcessWith(context.Background(),
				&envconfig.Config{
					Target:   &config,
					Lookuper: test.lookuper,
				})

			got := config.ConnectionURI()

			assert.Equal(t, test.want, got)
		})
	}
}
