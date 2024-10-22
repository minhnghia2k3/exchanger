package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minhnghia2k3/exchanger/docs"
	"github.com/minhnghia2k3/exchanger/internal/env"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// TODO: HANDLE NOT FOUND ROUTES, METHOD NOT ALLOW
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/v1"

	r.Route("/v1", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("%s/v1/swagger/doc.json", env.GetString("ADDR", "http://localhost:8080"))),
		))

		r.Get("/healthcheck", app.healthcheckHandler)

		r.Route("/tokens", func(r chi.Router) {
			r.Post("/authentication", app.createTokenHandler)
			r.Put("/activate", app.activateTokenHandler)
			r.Post("/refresh", app.refreshTokenHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.registerUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.validateAccessToken, app.findUserContext)

				r.Get("/", app.getUserHandler)
				r.Delete("/", app.deleteUserHandler)
			})
		})

		r.Route("/currencies", func(r chi.Router) {
			r.Get("/", app.listCurrenciesHandler)
			r.Post("/", app.addCurrencyHandler)
			r.Route("/{currencyID}", func(r chi.Router) {
				r.Use(app.currencyContext)

				r.Get("/", app.getCurrencyHandler)
				r.Patch("/", app.updateCurrencyHandler)
				r.Delete("/", app.deleteCurrencyHandler)
			})
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
