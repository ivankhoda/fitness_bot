package testutils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"

	localtestcontainers "fitness_bot/internal/models/testcontainers"
)

type testServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func NewTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	db, _, err := newTestDBContainer(t)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func newTestDBContainer(t *testing.T) (*pgxpool.Pool, testcontainers.Container, error) {
	t.Helper()

	ctx := context.Background()
	pgC, err := localtestcontainers.StartTestPGContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	t.Cleanup(func() {
		_ = pgC.Terminate(context.Background())
	})

	host, err := pgC.Host(ctx)
	if err != nil {
		return nil, nil, err
	}
	port, err := pgC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, nil, err
	}

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, err
	}

	var db *pgxpool.Pool
	for i := 0; i < 10; i++ {
		db, err = pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			err = db.Ping(ctx)
			if err == nil {
				break
			}
			db.Close()
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, nil, err
	}

	t.Cleanup(db.Close)

	if err := applyTestMigrations(ctx, db); err != nil {
		return nil, nil, err
	}

	return db, pgC, nil
}

func applyTestMigrations(ctx context.Context, db *pgxpool.Pool) error {
	migrationsDir, err := testMigrationsDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	var upMigrations []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		upMigrations = append(upMigrations, entry.Name())
	}

	sort.Strings(upMigrations)
	for _, name := range upMigrations {
		script, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			return err
		}

		if _, err := db.Exec(ctx, string(script)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
	}

	return nil
}

func testMigrationsDir() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve current file path")
	}

	return filepath.Join(filepath.Dir(filepath.Dir(currentFile)), "migrations"), nil
}
