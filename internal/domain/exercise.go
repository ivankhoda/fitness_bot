package domain

import (
	"net/http"
	"time"
)

type Exercise struct {
	ID         string
	Muscles    []string
	Category   string
	Difficulty string
}

type ExerciseRecord struct {
	UUID         string     `json:"uuid"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	MuscleGroups []string   `json:"muscle_groups"`
	Difficulty   string     `json:"difficulty"`
	Category     string     `json:"category"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type ExerciseRepository interface {
	Insert(exercise ExerciseRecord) (int, error)
	Upsert(exercise ExerciseRecord) error
	GetAll(f ExercsiesFilter) ([]ExerciseRecord, error)
	GetByID(id int) (*ExerciseRecord, error)
	Delete(id int) error
}

type ExercsiesFilter struct {
	MuscleGroups []string
	Category     string
	Difficulty   string
	Limit        string
}

type ExercisesFetcher interface {
	FetchExercises(r *http.Request) ([]ExerciseRecord, error)
}
