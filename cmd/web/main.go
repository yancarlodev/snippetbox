package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

var (
	cfg config
)

type application struct {
	logger *slog.Logger
}

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Path to static assets")

	flag.Parse()

	app := &application{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	app.logger.Info("starting server...", slog.String("addr", fmt.Sprintf("http://localhost%s", cfg.addr)))

	err := http.ListenAndServe(cfg.addr, app.routes())

	app.logger.Error(err.Error())
	os.Exit(1)
}
