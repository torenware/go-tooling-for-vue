package main

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ServeVueAssets(mux *chi.Mux, fileSys fs.ReadFileFS, pathToAssets string) error {

	sub, err := fs.Sub(fileSys, pathToAssets)
	if err != nil {
		return err
	}
	assetServer := http.FileServer(http.FS(sub))
	mux.Handle("/assets/*", http.StripPrefix("/assets", assetServer))
	return nil
}

func (app *application) routes() http.Handler {
	mux := chi.NewMux()

	mux.Get("/", app.home)

	err := ServeVueAssets(mux, app.vueDist, "dist/assets")
	if err != nil {
		app.errorLog.Println(err)
	}

	return mux
}
