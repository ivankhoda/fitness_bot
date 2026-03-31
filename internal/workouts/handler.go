package workouts

import (
	"encoding/json"
	"fitness_bot/internal/config"

	"net/http"
)

type WorkoutHandler struct {
	builder *WorkoutBuilder
	app     *config.Application
}

func (h *WorkoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		h.app.ErrorLog.Print("Method is not allowed:", r.Method)
		h.app.ServerError(w, http.StatusMethodNotAllowed, "Method is not allowed")
		return
	}
	if r.URL.Path != "/generate_workout" {
		h.app.NotFound(w)
		return
	}
	workout, err := h.builder.BuildWorkout(r)

	if err != nil {
		h.app.ErrorLog.Print("Error building workout:", err)
		http.Error(w, "Error building workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

func NewWorkoutHandler(builder *WorkoutBuilder, app config.Application) *WorkoutHandler {
	return &WorkoutHandler{builder: builder, app: &app}
}
