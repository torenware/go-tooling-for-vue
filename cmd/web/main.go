package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	vueglue "github.com/torenware/vite-go"
)

//go:embed dist
var dist embed.FS

type config struct {
	port int
	env  string
}

type application struct {
	config   config
	VueGlue  *vueglue.VueGlue
	vueDist  *embed.FS
	infoLog  *log.Logger
	errorLog *log.Logger
}

var app *application

func initConfig() config {
	var config config
	config.port = 4000
	config.env = "production"
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

	config := &vueglue.ViteConfig{
		Environment: "production",
		FS:          dist,
		AssetsPath:  "dist",
		URLPrefix:   "/src/",
	}

	glue, err := vueglue.NewVueGlue(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	var appl application
	appl.VueGlue = glue
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
