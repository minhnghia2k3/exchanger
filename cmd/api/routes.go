package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		//r.Route("/users", func(r chi.Router) {
		//	r.Post("/", app.createUserHandler)
		//
		//	r.Route("/{userID}", func(r chi.Router) {
		//		r.Get("/", app.getUserHandler)
		//		r.Patch("/", app.updateUserHandler)
		//		r.Delete("/", app.deleteUserHandler)
		//	})
		//})
		//
		r.Route("/currencies", func(r chi.Router) {
			r.Get("/", app.listCurrenciesHandler)
			//r.Post("/", app.addCurrencyHandler)
			//r.Route("/{currencyID}", func(r chi.Router) {
			//	r.Get("/", app.getCurrencyHandler)
			//	r.Patch("/", app.updateCurrencyHandler)
			//	r.Delete("/", app.deleteCurrencyHandler)
			//})
		})
		//
		//
		//r.Route("/rates", func(r chi.Router) {
		//	r.Route("/{base}/{target}", func(r chi.Router) {
		//		r.Get("/", app.getExchangeHandler)
		//		r.Post("/", app.createExchangeHandler)
		//		r.Patch("/", app.updateExchangeHandler)
		//		r.Delete("/", app.deleteExchangeHandler)
		//	})
		//})
		//
		//r.Route("/exchanges/pair", func(r chi.Router) {
		//	r.Get("/{base}/{target}", app.exchangePairHanlder)
		//})
	})

	return r
}
