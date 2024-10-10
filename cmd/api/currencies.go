package main

import (
	"github.com/minhnghia2k3/CoinScraper/internal/store"
	"net/http"
)

type ListCurrenciesInput struct {
	Code string `json:"code" validate:"omitempty,lte=255"`
	Name string `json:"name" validate:"omitempty,lte=255"`
	store.Filter
}

func (app *application) listCurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	input := ListCurrenciesInput{
		Code: readString(r, "code", ""),
		Name: readString(r, "name", ""),
		Filter: store.Filter{
			Page:         readInt(r, "page", 1),
			PageSize:     readInt(r, "page_size", 10),
			Sort:         readString(r, "sort", "id"),
			SortSafeList: []string{"id", "code", "name"},
		},
	}

	err := Validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	list, metadata, err := app.store.Currencies.List(r.Context(), input.Code, input.Name, input.Filter)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = writeJSON(w, http.StatusOK, envelop{"metadata": metadata, "data": list}); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *application) addCurrencyHandler(w http.ResponseWriter, r *http.Request) {

}
func (app *application) getCurrencyHandler(w http.ResponseWriter, r *http.Request) {

}
func (app *application) updateCurrencyHandler(w http.ResponseWriter, r *http.Request) {

}
func (app *application) deleteCurrencyHandler(w http.ResponseWriter, r *http.Request) {

}
