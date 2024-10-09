package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(data); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, value any) {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(value); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
