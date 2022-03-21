package main

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func (app *application) HandleDistDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	app.StaticWithMIME(w, r, path)
}

func (app *application) StaticWithMIME(w http.ResponseWriter, r *http.Request, path string) {
	if app.vueDist == nil {
		app.errorLog.Println("dist is not set up")
		http.Error(w, "dist is not initted", http.StatusBadRequest)
		return
	}

	file, err := app.vueDist.ReadFile("dist/" + path)
	if err != nil {
		app.errorLog.Printf("file not found at %s", path)
		http.Error(w, "no such file", http.StatusNotFound)
		return
	}
	ext := filepath.Ext(path)
	newMIME := ""
	switch ext {
	case ".js":
		newMIME = "text/javascript"
	case ".css":
		newMIME = "text/css"
	case ".map":
		newMIME = "application/json"
	default:
		app.infoLog.Println("extension was", ext)
	}

	if newMIME != "" {
		w.Header().Set("Content-Type", newMIME)
	}

	w.Write(file)
}

func (app *application) routes() http.Handler {
	mux := chi.NewMux()

	mux.Get("/", app.home)

	mux.Route("/assets", func(mux chi.Router) {
		mux.Get("/{subpath}", app.HandleDistDir)
	})

	return mux
}
