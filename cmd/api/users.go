package main

import (
	"errors"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,lte=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,lte=72"`
}

type UpdateUserPayload struct {
	Username string `json:"username" validate:"omitempty,lte=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,lte=72"`
}

func (app *application) activateTokenHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userCtx).(store.User)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdateUserPayload

	user := r.Context().Value(userCtx).(store.User)

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	model := new(store.User)

	if payload.Username != "" && payload.Username != user.Username {
		model.Username = payload.Username
	}

	if payload.Email != "" && payload.Email != user.Email {
		model.Email = payload.Email
	}

	if payload.Password != "" {
		err := model.Password.Set(payload.Password)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	}

	if err := app.store.Users.Update(r.Context(), model); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictErrorResponse(w, r, err)
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userCtx).(store.User)

	if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
