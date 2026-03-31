package main

import (
	"fitness_bot/internal/config"
	"fitness_bot/internal/exercises"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// errEnv := godotenv.Load("../../.env")
	// if errEnv != nil {
	// 	fmt.Println("Error loading .env file")
	// 	return
	// }
	addr := flag.String("addr", ":5000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	client := exercises.NewExercisesClient(os.Getenv("BOT_TOKEN"), os.Getenv("EXTERNAL_PROVIDER_URL")+"/api/exercises")

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  routes(app, client),
	}

	app.InfoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
