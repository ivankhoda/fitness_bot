package main

import (
	"context"
	"fitness_bot/internal/config"
	"fitness_bot/internal/exercises"
	"fitness_bot/internal/models"
	"fitness_bot/internal/scheduler"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load(".env", "../../.env")
	if errEnv != nil {
		fmt.Println("No .env file loaded, using process environment")
	}
	addr := flag.String("addr", ":5000", "HTTP network address")
	dsn := flag.String("dsn", databaseURLFromEnv(), "PostgreSQL DSN")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	exercisesAPIURL := exercisesAPIURLFromEnv()
	infoLog.Printf("Using exercises API endpoint: %s", exercisesAPIURL)
	client := exercises.NewExercisesClient(os.Getenv("BOT_TOKEN"), exercisesAPIURL)

	db, dBerr := pgxpool.New(context.Background(), *dsn)
	if dBerr != nil {
		errorLog.Fatal(dBerr)
	}
	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &config.Application{
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		Exercise:       &models.ExerciseModel{DB: db},
		SessionManager: sessionManager,
	}

	syncStateStore := scheduler.NewSyncStateStore(db)
	syncExercisesService := scheduler.NewSyncExercisesService(client, app.Exercise, syncStateStore, app)
	infoLog.Printf("Starting initial sync of exercises")
	if err := syncExercisesService.Run(); err != nil {
		errorLog.Printf("initial exercise sync failed: %v", err)
	}
	jobScheduler := scheduler.NewScheduler(syncExercisesService)
	if err := jobScheduler.Start(); err != nil {
		errorLog.Fatal(err)
	}
	defer jobScheduler.Stop()

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      routes(app, client),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.InfoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}

func databaseURLFromEnv() string {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" && !strings.Contains(databaseURL, "${") {
		return databaseURL
	}

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	sslMode := os.Getenv("POSTGRES_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		host,
		os.Getenv("POSTGRES_DB"),
		sslMode,
	)
}

func exercisesAPIURLFromEnv() string {
	baseURL := normalizeEnvValue(os.Getenv("EXTERNAL_PROVIDER_URL"))
	path := normalizeEnvValue(os.Getenv("PATH_TO_EXERCISES"))

	if baseURL == "" {
		return path
	}

	if path == "" {
		return baseURL
	}

	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func normalizeEnvValue(value string) string {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.Trim(trimmed, `"`)
	trimmed = strings.Trim(trimmed, `'`)
	return strings.TrimSpace(trimmed)
}
