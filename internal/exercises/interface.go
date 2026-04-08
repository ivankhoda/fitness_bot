package exercises

import (
	"fitness_bot/internal/domain"
	"net/http"
)

type ExercisesFetcher interface {
	FetchExercises(r *http.Request) ([]domain.ExerciseRecord, error)
}
