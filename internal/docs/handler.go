package docs

import (
	"embed"
	"fitness_bot/internal/config"
	"html/template"
	"net/http"
	"time"
)

//go:embed templates templates/partials
var templateFS embed.FS

type DocsHandler struct {
	app config.Application
}

func (h *DocsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		h.app.ClientError(w, 405)
		return
	}
	if r.URL.Path != "/docs" {
		h.app.NotFound(w)
		return
	}
	h.app.InfoLog.Print("Serving docs.html")

	ts, err := template.ParseFS(
		templateFS,
		"templates/base.tmpl",
		"templates/docs.tmpl",
		"templates/partials/*.tmpl",
	)
	if err != nil {
		h.app.ErrorLog.Println("template parse error:", err)

		h.app.ServerError(w, 500, "Internal Server Error")
		return
	}

	lang := r.URL.Query().Get("lang")
	if lang != "ru" && lang != "en" {
		lang = "ru"
	}

	data := map[string]any{
		"Lang": lang,
		"Year": time.Now().Year(),
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.app.ErrorLog.Print("template execute error:", err)
		h.app.ServerError(w, 500, "Internal Server Error")
		return
	}
}

func NewDocsHandler(app config.Application) *DocsHandler {
	return &DocsHandler{app: app}
}
