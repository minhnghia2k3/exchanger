package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type envelop map[string]any

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data interface{}) error {
	return writeJSON(w, status, envelop{"data": data})
}

func writeJSONError(w http.ResponseWriter, status int, msg string) error {
	return writeJSON(w, status, envelop{"error": msg})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, value any) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(value); err != nil {
		return err
	}
	return nil
}
