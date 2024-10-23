package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
)

// Get exchange rate by code
//
//	@Summary		Get exchange rate
//	@Description	get exchange rate by code
//	@Tags			Exchange Rates
//	@Accept			json
//	@Produce		json
//	@Param			base	path	string	true	"Base currency code"
//	@Param			target	path	string	true	"Target currency code"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	store.ExchangeRate
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/rates/{base}/{target} [get]
func (app *application) getExchangeRatesHandler(w http.ResponseWriter, r *http.Request) {
	base := chi.URLParam(r, "base")
	target := chi.URLParam(r, "target")

	if len(base) != 3 || len(target) != 3 {
		app.badRequestResponse(w, r, errors.New("invalid currency code"))
		return
	}

	rate, err := app.store.Rates.FindByCode(r.Context(), base, target)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, rate); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) addExchangeRateHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) updateExchangeRateHandler(w http.ResponseWriter, r *http.Request) {}
