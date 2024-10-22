package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	currencyCtx  = "currency"
	userCtx      = "user"
	foundUserCtx = "foundUser"
)

func (app *application) currencyContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currencyID, err := strconv.ParseInt(chi.URLParam(r, "currencyID"), 10, 64)

		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		currency, err := app.store.Currencies.Get(r.Context(), currencyID)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), currencyCtx, currency)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) findUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		user, err := app.store.Users.GetByID(r.Context(), userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), foundUserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) validateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestToken := r.Header.Get("Authorization")

		if requestToken == "" {
			app.unauthorizedResponse(w, r, ErrMissingJWT)
			return
		}

		if !strings.HasPrefix(requestToken, "Bearer ") {
			app.unauthorizedResponse(w, r, ErrInvalidJWT)
			return
		}

		splitToken := strings.Split(requestToken, "Bearer ")

		if len(splitToken) != 2 {
			app.unauthorizedResponse(w, r, ErrInvalidJWT)
			return
		}

		claims, err := app.verifyToken(splitToken[1])
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		// Get claims
		expiry, err := claims.GetExpirationTime()
		if err != nil {
			app.unauthorizedResponse(w, r, fmt.Errorf("%w: expiry time", ErrClaimsMissing))
			return
		}

		if expiry.Unix() < time.Now().Unix() {
			app.unauthorizedResponse(w, r, ErrExpiredJWT)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			app.unauthorizedResponse(w, r, fmt.Errorf("%w: user id", ErrClaimsMissing))
			return
		}

		user, err := app.store.Users.GetByID(context.Background(), int64(userID))
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.unauthorizedResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) adminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value(userCtx).(*store.User)

		if currentUser.Role.Level < 3 {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
