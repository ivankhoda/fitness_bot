package exercises

import (
	"encoding/json"
	"fitness_bot/internal/domain"
	"log"
	"net/http"
)

type ExercisesHandler struct {
	repo domain.ExerciseRepository
}

func (p *ExercisesHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := p.repo.GetAll(domain.ExercsiesFilter{})
	if err != nil {
		http.Error(w, "Error fetching exercises", http.StatusInternalServerError)
		return
	}

	log.Printf("Fetched exercises: %+v", exercises)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(exercises)
}

func (p *ExercisesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.GetExercises(w, r)
}

func NewExercisesHandler(repo domain.ExerciseRepository) *ExercisesHandler {
	return &ExercisesHandler{repo: repo}
}
