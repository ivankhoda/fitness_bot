package main

import (
	"fitness_bot/internal/config"
	"fitness_bot/internal/docs"
	"fitness_bot/internal/exercises"
	"fitness_bot/internal/workouts"
	"net/http"
)

func routes(app *config.Application, client *exercises.ExercisesClient) *http.ServeMux {

	workoutBuilder := workouts.NewWorkoutBuilder(client, *app)
	mux := http.NewServeMux()
	mux.Handle("/exercises", exercises.NewExercisesHandler(client, app.Exercise))
	mux.Handle("/docs", docs.NewDocsHandler(*app))

	mux.Handle("/generate_workout", workouts.NewWorkoutHandler(workoutBuilder, *app))
	return mux
}
