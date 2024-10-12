package main

import (
	"log"
	"net/http"
)

var (
	ErrInternalServerError = "internal server error"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error [%s] %s: %v\n", r.Method, r.URL.String(), err.Error())
	msg := "server encountered an error"

	writeJSONError(w, http.StatusInternalServerError, msg)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request [%s] %s: %v\n", r.Method, r.URL.String(), err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found %s %s", r.Method, r.URL.String())

	writeJSONError(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Conflict %s %s", r.Method, r.URL.String())

	writeJSONError(w, http.StatusConflict, err.Error())
}
