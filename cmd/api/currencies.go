package main

import (
	"errors"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"net/http"
)

type AddCurrencyInput struct {
	Code      string  `json:"code" validate:"required,len=3"`
	Name      string  `json:"name" validate:"required,len=50"`
	SymbolUrl *string `json:"symbol_url" validate:"omitempty,url"`
}

type UpdateCurrencyInput struct {
	Code      string  `json:"code" validate:"omitempty,len=3"`
	Name      string  `json:"name" validate:"omitempty,len=50"`
	SymbolUrl *string `json:"symbol_url" validate:"omitempty,url"`
}

// List currencies
//
//	@Summary		List currencies
//	@Description	get all currencies
//	@Tags			currencies
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"Current page"
//	@Param			page_size	query	int		false	"Page size"
//	@Param			sort		query	string	false	"Sort"
//	@Param			search		query	string	false	"Search"
//	@Success		200
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/currencies [get]
func (app *application) listCurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	input := store.Filter{
		Page:         readInt(r, "page", 1),
		PageSize:     readInt(r, "page_size", 10),
		Sort:         readString(r, "sort", "id"),
		Search:       readString(r, "search", ""),
		SortSafeList: []string{"id", "code", "name"},
	}

	err := Validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	list, metadata, err := app.store.Currencies.List(r.Context(), input)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = writeJSON(w, http.StatusOK, envelop{"metadata": metadata, "data": list}); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Add currency
//
//	@Summary		Add currency
//	@Description	add currency detail
//	@Tags			currencies
//	@Accept			json
//	@Produce		json
//	@Param			input	body	AddCurrencyInput	true	"Add currency"
//	@Success		200
//	@Failure		400	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Router			/currencies [post]
func (app *application) addCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	var input AddCurrencyInput

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	currency := store.Currency{
		Code:      input.Code,
		Name:      input.Name,
		SymbolUrl: input.SymbolUrl,
	}

	if err := app.store.Currencies.Insert(r.Context(), &currency); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, currency); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Get currency
//
//	@Summary		Get currency
//	@Description	get currency by id
//	@Tags			currencies
//	@Accept			json
//	@Produce		json
//	@Param			currencyID	path	int	false	"Currency ID"
//	@Success		200
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/currencies/{currencyID} [get]
func (app *application) getCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	currency := r.Context().Value(currencyCtx).(*store.Currency)

	if err := app.jsonResponse(w, http.StatusOK, currency); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Update currency
//
//	@Summary		Update currency
//	@Description	update currency by id
//	@Tags			currencies
//	@Accept			json
//	@Produce		json
//	@Param			currencyID	path	int					true	"currency ID"
//	@Param			input		body	UpdateCurrencyInput	true	"Update currency payload"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/currencies/{currencyID} [patch]
func (app *application) updateCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	var input UpdateCurrencyInput
	currency := r.Context().Value(currencyCtx).(*store.Currency)

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Currencies.Update(r.Context(), currency.ID, currency); err != nil {
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

// Delete currency
//
//	@Summary		Delete currency
//	@Description	delete currency by id
//	@Tags			currencies
//	@Accept			json
//	@Produce		json
//	@Param			currencyID	path	int	true	"currency ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Router			/currencies/{currencyID} [delete]
func (app *application) deleteCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	currency := r.Context().Value(currencyCtx).(*store.Currency)

	if err := app.store.Currencies.Update(r.Context(), currency.ID, currency); err != nil {
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
