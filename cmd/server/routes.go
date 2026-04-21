package main

import (
	"fitness_bot/internal/config"
	"fitness_bot/internal/docs"
	"fitness_bot/internal/exercises"
	"fitness_bot/internal/workouts"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func routes(app *config.Application, client *exercises.ExercisesClient) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})
	workoutBuilder := workouts.NewWorkoutBuilder(client, *app)

	router.Handler(http.MethodGet, "/docs/assets/*filepath", http.StripPrefix("/docs/assets/", docs.NewAssetsHandler()))
	router.Handler(http.MethodGet, "/exercises", exercises.NewExercisesHandler(app.Exercise))
	router.Handler(http.MethodGet, "/exercises/:id", exercises.NewExerciseHandler(app.Exercise))

	router.Handler(http.MethodGet, "/docs", docs.NewDocsHandler(*app))
	router.Handler(http.MethodGet, "/generate_workout", workouts.NewWorkoutHandler(workoutBuilder, *app))
	router.Handler(http.MethodGet, "/debug/session", workouts.NewSessionDebugHandler(*app))

	standart := alice.New(app.RecoverPanic, app.LogRequest, secureHeaders, app.SessionManager.LoadAndSave)
	return standart.Then(router)
}
