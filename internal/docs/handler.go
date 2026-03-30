package docs

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type DocsHandler struct {
	builder *DocsHandler
}

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

	ts, err := template.ParseFiles(
		"../../ui/html/docs/base.tmpl",
		"../../ui/html/docs/docs.tmpl",
		"../../ui/html/docs/partials/get.tmpl",
		"../../ui/html/docs/partials/post.tmpl",
		"../../ui/html/docs/partials/put.tmpl",
		"../../ui/html/docs/partials/patch.tmpl",
		"../../ui/html/docs/partials/delete.tmpl",
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
