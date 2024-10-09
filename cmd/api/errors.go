package main

import (
	"log"
	"net/http"
)

var (
	ErrInternalServerError = "internal server error"
)

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error [%s] %s: %v\n", r.Method, r.URL.String(), err.Error())
	msg := "server encountered an error"
	http.Error(w, msg, http.StatusInternalServerError)
}
