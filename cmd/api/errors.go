package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrMissingJWT    = errors.New("missing token")
	ErrInvalidJWT    = errors.New("invalid token")
	ErrExpiredJWT    = errors.New("expired token")
	ErrClaimsMissing = errors.New("missing claims")
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.LogAttrs(context.Background(),
		slog.LevelError,
		"Internal server error:",
		slog.String("URL", r.URL.String()),
		slog.String("method", r.Method),
		slog.String("error", err.Error()),
	)
	msg := "server encountered an error"

	writeJSONError(w, http.StatusInternalServerError, msg)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.LogAttrs(context.Background(),
		slog.LevelWarn,
		"Bad request:",
		slog.String("URL", r.URL.String()),
		slog.String("method", r.Method),
		slog.String("error", err.Error()),
	)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.LogAttrs(context.Background(),
		slog.LevelWarn,
		"Not found:",
		slog.String("URL", r.URL.String()),
		slog.String("method", r.Method),
		slog.String("error", err.Error()),
	)

	writeJSONError(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.LogAttrs(context.Background(),
		slog.LevelWarn,
		"Conflict record:",
		slog.String("URL", r.URL.String()),
		slog.String("method", r.Method),
		slog.String("error", err.Error()),
	)

	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.LogAttrs(context.Background(),
		slog.LevelWarn,
		"Unauthorized:",
		slog.String("URL", r.URL.String()),
		slog.String("method", r.Method),
		slog.String("error", err.Error()),
	)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.LogAttrs(context.Background(),
		slog.LevelWarn,
		"Forbidden:",
		slog.String("URL", r.URL.String()),
		slog.String("method", r.Method),
	)
	msg := "You are not allowed to access this route"
	writeJSONError(w, http.StatusForbidden, msg)
}
