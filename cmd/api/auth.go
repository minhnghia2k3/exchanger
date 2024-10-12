package main

import (
	"errors"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := store.User{
		RoleID:    1,
		Username:  payload.Username,
		Email:     payload.Email,
		CreatedAt: time.Now(),
	}

	err := user.Password.Set(payload.Password)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Users.Insert(r.Context(), &user); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

}
