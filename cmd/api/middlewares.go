package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
	"strconv"
)

const currencyCtx = "currency"

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
