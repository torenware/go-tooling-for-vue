package main

import (
	// New import

	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	app.RenderTemplate(w, r, "home", nil)
}
