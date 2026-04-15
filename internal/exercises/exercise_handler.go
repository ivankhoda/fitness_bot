package exercises

import (
	"encoding/json"
	"fitness_bot/internal/domain"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ExerciseHandler struct {
	repo domain.ExerciseRepository
}

func (p *ExerciseHandler) GetExercise(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	idStr := params.ByName("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
			return
		}
		exercise, err := p.repo.GetByID(id)
		if err != nil {
			http.Error(w, "Exercise not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercise)
		return
	}

}

func (p *ExerciseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.GetExercise(w, r)
}

func NewExerciseHandler(repo domain.ExerciseRepository) *ExerciseHandler {
	return &ExerciseHandler{repo: repo}
}
