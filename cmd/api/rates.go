package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/minhnghia2k3/exchanger/internal/env"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
	"time"
)

var (
	errInvalidCurrencyCode     = errors.New("invalid currency code")
	errPairAlreadyExists       = errors.New("the pair of exchange rate is already exists")
	errUnsupportedCurrencyCode = errors.New("unsupported currency code")
	errUnsupportedCurrencyFmt  = "%s code not supported"
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

	if !validCurrencyCode(base, target) {
		app.badRequestResponse(w, r, errors.New("invalid currency code"))
		return
	}

	rate, err := app.store.Rates.GetByPair(r.Context(), base, target)
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

type ExchangeRateData struct {
	Result             string  `json:"result"`
	TimeLastUpdateUnix int     `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string  `json:"time_last_update_utc"`
	TimeNextUpdateUnix int     `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string  `json:"time_next_update_utc"`
	BaseCode           string  `json:"base_code"`
	TargetCode         string  `json:"target_code"`
	ConversionRate     float64 `json:"conversion_rate"`
}

// Add exchange rate of pair
//
//	@Summary		Add exchange rate
//	@Description	add exchange rate by code
//	@Tags			Exchange Rates
//	@Accept			json
//	@Produce		json
//	@Param			base	path	string	true	"Base currency code"
//	@Param			target	path	string	true	"Target currency code"
//	@Security		ApiKeyAuth
//	@Success		201	{object}	store.ExchangeRate
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		403	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/rates/{base}/{target} [post]
func (app *application) addExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	base := chi.URLParam(r, "base")
	target := chi.URLParam(r, "target")

	if !validCurrencyCode(base, target) {
		app.badRequestResponse(w, r, errInvalidCurrencyCode)
		return
	}

	// 1. Validate code pair
	ok, err := isPairCodeExists(r.Context(), app.store.Rates, base, target)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		app.internalServerError(w, r, err)
		return
	}

	if ok {
		app.conflictErrorResponse(w, r, errPairAlreadyExists)
		return
	}

	// 2. Validate supported currency codes (for both base and target)
	for _, code := range []string{base, target} {
		if ok, err = isSupportedCode(r.Context(), app.store.Currencies, code); err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, fmt.Errorf(errUnsupportedCurrencyFmt, code))
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
	}

	// 3. Get exchange rate from api
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s",
		env.GetString("EXCHANGER_RATE_API", ""),
		base,
		target,
	)

	var data ExchangeRateData
	err = GetAPI(url, &data)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// 3. Parse date and store exchange rate to db
	exchangeRate, err := buildExchangeRate(data, base, target)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Rates.Save(r.Context(), exchangeRate)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusCreated, exchangeRate); err != nil {
		app.internalServerError(w, r, err)
	}
}

func buildExchangeRate(data ExchangeRateData, base, target string) (*store.ExchangeRate, error) {
	lastUpdate, err := time.Parse(time.RFC1123, data.TimeLastUpdateUtc)
	if err != nil {
		return nil, err
	}

	nextUpdate, err := time.Parse(time.RFC1123, data.TimeNextUpdateUtc)
	if err != nil {
		return nil, err
	}

	return &store.ExchangeRate{
		NextUpdate: nextUpdate,
		BaseCode:   base,
		TargetCode: target,
		LastUpdate: lastUpdate,
		Rate:       data.ConversionRate,
	}, nil
}

func validCurrencyCode(base, target string) bool {
	return len(base) == 3 && len(target) == 3
}

func isPairCodeExists(ctx context.Context, rates store.IExchangeRates, base, target string) (bool, error) {
	rate, err := rates.GetByPair(ctx, base, target)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	return rate != nil, nil
}

func isSupportedCode(ctx context.Context, currencies store.ICurrencies, code string) (bool, error) {
	_, err := currencies.GetByCode(ctx, code)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Update exchange rate conversion
//
//	@Summary		Update exchange rate conversion
//	@Description	update exchange rate conversion
//	@Tags			Exchange Rates
//	@Accept			json
//	@Produce		json
//	@Param			base	path	string	true	"Base currency code"
//	@Param			target	path	string	true	"Target currency code"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	store.ExchangeRate
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		403	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/rates/{base}/{target} [patch]
func (app *application) updateExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	base := chi.URLParam(r, "base")
	target := chi.URLParam(r, "target")

	if !validCurrencyCode(base, target) {
		app.badRequestResponse(w, r, errInvalidCurrencyCode)
		return
	}

	rate, err := app.store.Rates.GetByPair(r.Context(), base, target)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s",
		env.GetString("EXCHANGER_RATE_API", ""),
		base,
		target,
	)

	var data ExchangeRateData
	err = GetAPI(url, &data)
	if err != nil {
		switch {
		case errors.Is(err, errUnsupportedCurrencyCode):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	lastUpdate, err := time.Parse(time.RFC1123, data.TimeLastUpdateUtc)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	nextUpdate, err := time.Parse(time.RFC1123, data.TimeNextUpdateUtc)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	rate.LastUpdate = lastUpdate
	rate.NextUpdate = nextUpdate
	rate.Rate = data.ConversionRate

	err = app.store.Rates.Update(r.Context(), rate)
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
