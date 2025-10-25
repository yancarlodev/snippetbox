package main

import (
	"log/slog"
	"net/http"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))

	httpError := http.StatusInternalServerError
	http.Error(w, http.StatusText(httpError), httpError)
}

func (app *application) httpError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}
