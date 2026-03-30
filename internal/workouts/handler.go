package workouts

import (
	"encoding/json"
	"net/http"
)

type WorkoutHandler struct {
	builder *WorkoutBuilder
}

func (h *WorkoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/generate_workout" {
		http.NotFound(w, r)
		return
	}
	workout, err := h.builder.BuildWorkout(r)

	if err != nil {
		http.Error(w, "Error building workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

func NewWorkoutHandler(builder *WorkoutBuilder) *WorkoutHandler {
	return &WorkoutHandler{builder: builder}
}
