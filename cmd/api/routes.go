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

//	@title			Exchanger API
//	@version		1.0
//	@description	Exchanger Open API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1
func (app *application) routes() http.Handler {
	r := chi.NewRouter()

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
