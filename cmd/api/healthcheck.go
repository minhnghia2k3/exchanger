package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"version":     version,
		"environment": "development",
		"status":      "ok",
	}

	app.writeJSON(w, r, http.StatusOK, data)
}
