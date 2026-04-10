package main

import (
	"fitness_bot/internal/config"
	"fitness_bot/internal/docs"
	"fitness_bot/internal/exercises"
	"fitness_bot/internal/workouts"
	"net/http"

	"github.com/justinas/alice"
)

func routes(app *config.Application, client *exercises.ExercisesClient) http.Handler {

	workoutBuilder := workouts.NewWorkoutBuilder(client, *app)
	mux := http.NewServeMux()
	mux.Handle("/exercises", exercises.NewExercisesHandler(client, app.Exercise))
	mux.Handle("/docs", docs.NewDocsHandler(*app))

	mux.Handle("/generate_workout", workouts.NewWorkoutHandler(workoutBuilder, *app))

	standart := alice.New(app.RecoverPanic, app.LogRequest, secureHeaders)
	return standart.Then(mux)
}
