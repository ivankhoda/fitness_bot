package main

import (
	"context"
	"fitness_bot/internal/config"
	"fitness_bot/internal/exercises"
	"fitness_bot/internal/scheduler"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

	client := exercises.NewExercisesClient(os.Getenv("BOT_TOKEN"), os.Getenv("EXTERNAL_PROVIDER_URL")+os.Getenv("PATH_TO_EXERCISES"))

	db, dBerr := pgxpool.New(context.Background(), *dsn)
	if dBerr != nil {
		errorLog.Fatal(dBerr)
	}
	defer db.Close()
	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Exercise: &exercises.ExerciseModel{DB: db},
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
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  routes(app, client),
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
