package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"snippetbox.srcrer.duckdns.org/internal/web"
)

type config struct {
	addr      string
	staticDir string
}

var (
	cfg    config
)

type application struct {
	logger *slog.Logger
}

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Path to static assets")

	flag.Parse()

	mux := http.NewServeMux()

	app := &application{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	fileServer := http.FileServer(web.NeuteredFileSystem{Fs: http.Dir(cfg.staticDir)})
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	app.logger.Info("starting server...", slog.String("addr", fmt.Sprintf("http://localhost%s", cfg.addr)))

	err := http.ListenAndServe(cfg.addr, mux)

	app.logger.Error(err.Error())
	os.Exit(1)
}
