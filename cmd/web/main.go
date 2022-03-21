package main

import (
	"embed"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"vite-tester/internal/vueglue"
	// "github.com/vearutop/statigz"
	// "github.com/vearutop/statigz/brotli"
)

//go:embed dist/assets
var dist embed.FS

type config struct {
	port int
	env  string
}

type application struct {
	config   config
	ViteGlue *vueglue.VueGlue
	vueDist  *embed.FS
	infoLog  *log.Logger
	errorLog *log.Logger
}

var app *application

func initConfig() config {
	var config config
	config.port = 4000
	config.env = "development"
	return config
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting server in %s mode on port %d", app.config.env, app.config.port)

	return srv.ListenAndServe()

}

func main() {

	// support vue fields in templates.
	gob.Register(vueglue.VueGlue{})

	// s, err := fs.Sub(dist, "dist")
	// if err != nil {
	// 	fmt.Println("sub failed")
	// 	return
	// }

	glue, err := vueglue.NewVueGlue(&dist)
	if err != nil {
		fmt.Println(err)
		return
	}

	var appl application
	appl.ViteGlue = glue
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app = &appl
	app.vueDist = &dist
	app.infoLog = infoLogger
	app.errorLog = errorLogger
	app.config = initConfig()

	// err = http.ListenAndServe(":5001", statigz.FileServer(s.(fs.ReadDirFS), brotli.AddEncoding))
	// if err != nil {
	// 	fmt.Println("done", err)
	// }
	err = app.serve()
	if err != nil {
		errorLogger.Println(err)
	}
}
