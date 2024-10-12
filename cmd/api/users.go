package main

import (
	"errors"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
)

type UpdateUserPayload struct {
	Username string  `json:"username" validate:"omitempty,min=3,lte=50"`
	Email    string  `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=8,lte=72"`
}

func (app *application) activateTokenHandler(w http.ResponseWriter, r *http.Request) {
}

// Get user handler
//
//	@Summary		Get user
//	@Description	get user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"currency ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userCtx).(*store.User)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Update user
//
//	@Summary		Update user
//	@Description	update user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int					true	"user ID"
//	@Param			input	body		UpdateUserPayload	true	"Update user payload"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/{userID} [patch]
func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdateUserPayload

	user := r.Context().Value(userCtx).(*store.User)

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Username != "" && payload.Username != user.Username {
		user.Username = payload.Username
	}

	if payload.Email != "" && payload.Email != user.Email {
		user.Email = payload.Email
	}

	if payload.Password != nil {
		err := user.Password.Set(*payload.Password)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	}

	if err := app.store.Users.Update(r.Context(), user); err != nil {
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

// Delete user
//
//	@Summary		Delete user
//	@Description	delete user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"user ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/{userID} [delete]
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userCtx).(*store.User)

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
