package docs

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"
)

//go:embed templates/base.tmpl templates/docs.tmpl templates/partials/*.tmpl
var templateFS embed.FS

type DocsHandler struct{}

func (h *DocsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/docs" {
		http.NotFound(w, r)
		return
	}
	log.Println("Serving docs.html")

	ts, err := template.ParseFS(
		templateFS,
		"templates/base.tmpl",
		"templates/docs.tmpl",
		"templates/partials/*.tmpl",
	)
	if err != nil {
		log.Println("template parse error:", err)
		http.Error(w, "Internal Server Error", 500)
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
		log.Println("template execute error:", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}
