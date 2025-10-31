package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

type dbConnection struct {
	user   string
	passwd string
	host   string
	port   string
}

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Path to static assets")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbConnectionOpts := &dbConnection{
		user:   os.Getenv("MYSQL_USER"),
		passwd: os.Getenv("MYSQL_PASSWORD"),
		host:   os.Getenv("DB_HOST"),
		port:   os.Getenv("DB_PORT"),
	}
	db, err := openDB(dbConnectionOpts)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		logger: logger,
	}

	app.logger.Info("starting server...", slog.String("addr", fmt.Sprintf("http://localhost%s", cfg.addr)))

	if err := http.ListenAndServe(cfg.addr, app.routes()); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(connOpts *dbConnection) (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/snippetbox?parseTime=true",
			connOpts.user,
			connOpts.passwd,
			connOpts.host,
			connOpts.port,
		),
	)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
