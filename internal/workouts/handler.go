package workouts

import (
	"encoding/json"
	"fitness_bot/internal/config"

	"net/http"
	"strings"
)

type sessionDebugResponse struct {
	SelectedDifficulty string   `json:"selectedDifficulty"`
	CurrentSource      string   `json:"currentSource"`
	SessionSource      string   `json:"sessionSource"`
	Keys               []string `json:"keys"`
}

type WorkoutHandler struct {
	builder *WorkoutBuilder
	app     *config.Application
}

func (h *WorkoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	workout, err := h.builder.BuildWorkout(r)

	if err != nil {
		h.app.ErrorLog.Print("Error building workout:", err)
		http.Error(w, "Error building workout", http.StatusInternalServerError)
		return
	}

	selectedDifficulty := r.URL.Query().Get("difficulty")
	if selectedDifficulty == "" {
		selectedDifficulty = workout.Difficulty
	}

	requestSource := detectRequestSource(r)
	requestIP := detectIpAddress(r)

	h.app.SessionManager.Put(r.Context(), "selectedDifficulty", selectedDifficulty)
	h.app.SessionManager.Put(r.Context(), "requestSource", requestSource)
	h.app.SessionManager.Put(r.Context(), "requestIP", requestIP)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

func NewWorkoutHandler(builder *WorkoutBuilder, app config.Application) *WorkoutHandler {
	return &WorkoutHandler{builder: builder, app: &app}
}

type SessionDebugHandler struct {
	app *config.Application
}

func (h *SessionDebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := sessionDebugResponse{
		SelectedDifficulty: h.app.SessionManager.GetString(r.Context(), "selectedDifficulty"),
		CurrentSource:      detectRequestSource(r),
		SessionSource:      h.app.SessionManager.GetString(r.Context(), "requestSource"),
		Keys:               h.app.SessionManager.Keys(r.Context()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func NewSessionDebugHandler(app config.Application) *SessionDebugHandler {
	return &SessionDebugHandler{app: &app}
}

func detectIpAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = r.RemoteAddr
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		return ip[:colonIndex]
	}
	return ip
}

func detectRequestSource(r *http.Request) string {
	if source := normalizeRequestSource(r.URL.Query().Get("source")); source != "" {
		return source
	}

	if source := normalizeRequestSource(r.Header.Get("X-Request-Source")); source != "" {
		return source
	}

	userAgent := strings.ToLower(r.UserAgent())
	switch {
	case strings.Contains(userAgent, "postmanruntime"):
		return "postman"
	case strings.Contains(userAgent, "telegrambot"), strings.Contains(userAgent, "telegram"):
		return "telegram"
	case strings.Contains(userAgent, "mozilla/"):
		return "web"
	case strings.Contains(userAgent, "curl/"), strings.Contains(userAgent, "httpie"), strings.Contains(userAgent, "go-http-client"):
		return "api"
	default:
		return "unknown"
	}
}

func normalizeRequestSource(raw string) string {
	source := strings.ToLower(strings.TrimSpace(raw))
	switch source {
	case "web", "browser":
		return "web"
	case "tg", "telegram", "telegram-bot", "telegram_bot":
		return "telegram"
	case "postman":
		return "postman"
	case "api", "mobile", "ios", "android":
		return "api"
	default:
		return ""
	}
}
