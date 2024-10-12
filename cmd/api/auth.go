package main

import (
	"errors"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
	"time"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=3,lte=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,lte=72"`
}

type UserWithToken struct {
	store.User
	Token string `json:"token"`
}

// Register user
//
//	@Summary		Register user
//	@Description	register user
//	@Tags			authentications
//	@Accept			json
//	@Produce		json
//	@Param			input	body		RegisterUserPayload	true	"Register user payload"
//	@Success		201		{object}	UserWithToken
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		409		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentications/users [post]
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

	// TODO: Generate token
	token := "test_token"

	err := user.Password.Set(payload.Password)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Users.CreateAndInvite(r.Context(), token, &user); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: token,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

}
