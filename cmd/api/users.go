package main

import (
	"errors"
	"github.com/google/uuid"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"log"
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
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			input	body		RegisterUserPayload	true	"Register user payload"
//	@Success		201		{object}	UserWithToken
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		409		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users [post]
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

	// Generate invite activation token
	plainToken := uuid.Must(uuid.NewRandom()).String() // Return to user

	token := SHA256Hash(plainToken) // Store to the database

	if err = app.store.Users.CreateAndInvite(r.Context(), token, &user); err != nil {
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
		Token: plainToken,
	}

	go func(plainToken string) {
		dynamicData := map[string]any{
			"activationToken": plainToken,
		}

		err = app.mailer.Send(user.Email, "user_invitations.tmpl", dynamicData)
		if err != nil {
			log.Printf("Error sending invitation email %v\n", err)
		} else {
			log.Printf("Sending email to user %q successfully", user.Email)
		}
	}(plainToken)

	if err = app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Get user handler
//
//	@Summary		Get user
//	@Description	get user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"user ID"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(foundUserCtx).(*store.User)

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
	user := r.Context().Value(foundUserCtx).(*store.User)

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
