package main

import (
	"errors"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
)

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

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
