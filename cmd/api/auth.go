package main

import (
	"errors"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
)

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// Login account
//
//	@Summary		Login user account
//	@Description	login user account and generate access token
//	@Tags			tokens
//	@Accept			json
//	@Produce		json
//	@Param			input	body		LoginPayload	true	"Login payload"
//	@Success		201		{object}	envelop
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/tokens/authentication [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.store.Users.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrUnauthorized):
			app.unauthorizedResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Generate jwt token
	tokenString, err := app.generateToken(user.ID, user.Email)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusCreated, tokenString); err != nil {
		app.internalServerError(w, r, err)
	}
}

type ActivationUserInvitations struct {
	Token string `json:"token"`
}

// Activate user
//
//	@Summary		Activate user account
//	@Description	Activate user by invitation token
//	@Tags			tokens
//	@Accept			json
//	@Produce		json
//	@Param			input	body	ActivationUserInvitations	true	"Invitation token"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/tokens/activate [put]
func (app *application) activateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload ActivationUserInvitations

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Users.Activate(r.Context(), payload.Token); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
