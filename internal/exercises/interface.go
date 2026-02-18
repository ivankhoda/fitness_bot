package exercises

import "net/http"

type ExercisesFetcher interface {
	FetchExercises(r *http.Request) ([]ExerciseRecord, error)
}
