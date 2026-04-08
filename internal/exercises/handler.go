package exercises

import (
	"encoding/json"
	"fitness_bot/internal/domain"
	"log"
	"net/http"
)

type ExercisesHandler struct {
	client domain.ExercisesFetcher
	repo   domain.ExerciseRepository
}

func (p *ExercisesHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := p.client.FetchExercises(r)
	if err != nil {
		http.Error(w, "Error fetching exercises", http.StatusInternalServerError)
		return
	}

	// for _, exercise := range exercises {
	// 	_, err = p.repo.Insert(exercise)
	// 	if err != nil {
	// 		http.Error(w, "Error saving exercise", http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	log.Printf("Fetched exercises: %+v", exercises)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(exercises)
}

func (p *ExercisesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	p.GetExercises(w, r)
}

func NewExercisesHandler(client domain.ExercisesFetcher, repo domain.ExerciseRepository) *ExercisesHandler {
	return &ExercisesHandler{client: client, repo: repo}
}
