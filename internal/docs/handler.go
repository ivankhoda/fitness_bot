package docs

import (
	"embed"
	"fitness_bot/internal/config"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

//go:embed templates templates/partials assets
var templateFS embed.FS

type DocsHandler struct {
	app config.Application
}

func (h *DocsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func NewAssetsHandler() http.Handler {
	assetsFS, err := fs.Sub(templateFS, "assets")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(assetsFS))
}
