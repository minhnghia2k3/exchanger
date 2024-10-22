package main

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrMissingJWT    = errors.New("missing token")
	ErrInvalidJWT    = errors.New("invalid token")
	ErrExpiredJWT    = errors.New("expired token")
	ErrClaimsMissing = errors.New("missing claims")
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[%s] %s -- Internal server error: %v\n", r.Method, r.URL.String(), err.Error())
	msg := "server encountered an error"

	writeJSONError(w, http.StatusInternalServerError, msg)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[%s] %s -- Bad request: %v\n", r.Method, r.URL.String(), err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[%s] %s -- Not found\n", r.Method, r.URL.String())

	writeJSONError(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[%s] %s -- Conflict\n", r.Method, r.URL.String())

	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[%s] %s -- Unauthorized\n", r.Method, r.URL.String())

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s -- Forbidden\n", r.Method, r.URL.String())
	msg := "You are not allowed to access this route"
	writeJSONError(w, http.StatusForbidden, msg)
}
