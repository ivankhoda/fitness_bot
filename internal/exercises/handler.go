package exercises

import (
	"encoding/json"
	"net/http"
)

type ExercisesHandler struct {
	client ExercisesFetcher
}

func (p *ExercisesHandler) GetExercises(w http.ResponseWriter, r *http.Request) {

	exercises, err := p.client.FetchExercises(r)
	if err != nil {
		http.Error(w, "Error fetching exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}

func NewExercisesHandler(client ExercisesFetcher) *ExercisesHandler {
	return &ExercisesHandler{client: client}
}
