package main

import (
	"context"
	"fitness_bot/internal/assert"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDatabaseURLFromEnv_BuildsURLFromParts(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("POSTGRES_USER", "test_postgress_user")
	t.Setenv("POSTGRES_PASSWORD", "test_postgress_password")
	t.Setenv("POSTGRES_DB", "TEST_DB")
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_SSLMODE", "disable")

	got := databaseURLFromEnv()
	want := "postgres://test_postgress_user:test_postgress_password@localhost:5432/TEST_DB?sslmode=disable"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDatabaseURLFromEnv_FallbackValues(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("POSTGRES_DB", "TEST_DB")
	t.Setenv("POSTGRES_HOST", "")
	t.Setenv("POSTGRES_SSLMODE", "")

	got := databaseURLFromEnv()
	want := "postgres://user:password@localhost:5432/TEST_DB?sslmode=disable"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDatabaseURLFromEnv_Table(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{
			name: "DATABASE_URL takes precedence",
			env: map[string]string{
				"DATABASE_URL":      "postgres://user:password@localhost:5432/TEST_DB?sslmode=disable",
				"POSTGRES_USER":     "other_user",
				"POSTGRES_PASSWORD": "other_password",
				"POSTGRES_DB":       "OTHER_DB",
			},
			want: "postgres://user:password@localhost:5432/TEST_DB?sslmode=disable",
		},
		{name: "POSTGRES_HOST and POSTGRES_SSLMODE fallback",
			env: map[string]string{
				"DATABASE_URL":      "",
				"POSTGRES_USER":     "user",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "TEST_DB",
				"POSTGRES_HOST":     "db_host",
				"POSTGRES_SSLMODE":  "require",
			},
			want: "postgres://user:password@db_host:5432/TEST_DB?sslmode=require",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			got := databaseURLFromEnv()
			assert.Equal(t, tt.want, got)

		})
	}
}

func TestDBConnectionStringWithEnv(t *testing.T) {
	const DATABASE_URL_TEST = "postgres://postgres:@localhost:5432/fitness_bot_test?sslmode=disable"
	db, dBerr := pgxpool.New(context.Background(), DATABASE_URL_TEST)
	if dBerr != nil {
		log.Fatal(dBerr)
	}
	defer db.Close()
}
