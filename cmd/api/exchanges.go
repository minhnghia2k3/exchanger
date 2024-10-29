package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/minhnghia2k3/exchanger/internal/store"
)

var (
	ErrInvalidAmount = errors.New("invalid converted amount")
)

// Exchange rate handler
//
//	@Summary		exchange rate
//	@Description	exchange rate
//	@Tags			Exchanges
//	@Accept			json
//	@Produce		json
//	@Param			base	path	string	true	"Base currency code"
//	@Param			target	path	string	true	"Target currency code"
//	@Param			amount	path	string	true	"Amount to convert"
//	@Success		200	{object}	float64
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/exchanges/pair/{base}/{target}/{amount} [post]
func (app *application) exchangePairHandler(w http.ResponseWriter, r *http.Request) {
	// Get base
	base := chi.URLParam(r, "base")

	// Get target
	target := chi.URLParam(r, "target")

	// Get amount
	amountParam := chi.URLParam(r, "amount")
	amount, err := strconv.ParseFloat(amountParam, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if amount < 1 {
		app.badRequestResponse(w, r, ErrInvalidAmount)
		return
	}

	// Validate input
	if !validCurrencyCode(base, target) {
		app.badRequestResponse(w, r, errInvalidCurrencyCode)
		return
	}

	// Get rates by pair
	rates, err := app.store.Rates.GetByPair(r.Context(), base, target)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			// If there is no rates in database => fetch and store
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Get rates and calculate the target result.
	result := amount * rates.Rate

	// Store data to transaction history
	transaction := &store.Transaction{
		UserID:          nil,
		BaseCode:        rates.BaseCode,
		TargetCode:      rates.TargetCode,
		ConvertedAmount: amount,
		ConvertedRate:   rates.Rate,
		Result:          result,
	}

	err = app.store.Transactions.Save(r.Context(), transaction)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Return result
	if err = app.jsonResponse(w, http.StatusOK, result); err != nil {
		app.internalServerError(w, r, err)
	}
}
