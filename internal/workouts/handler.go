package workouts

import (
	"encoding/json"
	"net/http"
)

type WorkoutHandler struct {
	builder *WorkoutBuilder
}

func (h *WorkoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
