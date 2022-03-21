package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (app *application) HandleDistDir(w http.ResponseWriter, r *http.Request) {
	subpath := chi.URLParam(r, "subpath")
	path := r.URL.Path[1:]
	app.infoLog.Println("subpath, path:", subpath, path)
	app.StaticWithMIME(w, r, path)
}

func (app *application) StaticWithMIME(w http.ResponseWriter, r *http.Request, path string) {
	if app.vueDist == nil {
		app.errorLog.Println("dist is not set up")
		http.Error(w, "dist is not initted", http.StatusBadRequest)
		return
	}
	// what do we have?
	dir, _ := app.vueDist.ReadDir("assets")
	for _, item := range dir {
		app.infoLog.Println(item.Name())
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
	default:
		app.infoLog.Println("extension was", ext)
	}

	if newMIME != "" {
		app.infoLog.Println("mime correction should now be", newMIME)

		w.Header().Set("Content-Type", newMIME)
	} else {
		app.infoLog.Println("no mime correction")
	}

	w.Write(file)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// from: chi examples
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Println("dist router called")
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fsSys := http.StripPrefix(pathPrefix, http.FileServer(root))
		// correct MIME types
		ext := filepath.Ext(r.URL.Path)
		newMIME := ""
		switch ext {
		case ".js":
			newMIME = "text/javascript"
		case ".css":
			newMIME = "text/css"
		default:
			app.infoLog.Println("extension was", ext)
		}

		if newMIME != "" {
			app.infoLog.Println("mime correction should now be", newMIME)

			w.Header().Set("Content-Type", newMIME)
		} else {
			app.infoLog.Println("no mime correction")
		}

		fsSys.ServeHTTP(w, r)
	})
}

func (app *application) routes() http.Handler {
	mux := chi.NewMux()

	mux.Get("/", app.home)

	mux.Route("/assets", func(mux chi.Router) {
		mux.Get("/{subpath}", app.HandleDistDir)
	})

	return mux
}
