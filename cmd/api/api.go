package main

import (
	"context"
	"fmt"
	"github.com/minhnghia2k3/exchanger/internal/mail"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"log/slog"
	"net/http"
	"time"
)

type application struct {
	config config
	store  *store.Storage
	mailer *mail.Mailer
	logger *slog.Logger
}

type config struct {
	port       int
	env        string
	dbConfig   dbConfig
	mailConfig mailConfig
	jwtConfig  jwtConfig
}

type jwtConfig struct {
	issuer        string
	secret        string
	expiry        string
	refreshExpiry string
}

type dbConfig struct {
	dsn         string
	maxIdleConn int
	maxOpenConn int
	maxIdleTime string
}

type mailConfig struct {
	sender   string
	host     string
	port     int
	username string
	password string
}

func (app *application) serve() error {
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.LogAttrs(context.Background(),
		slog.LevelInfo,
		"Server is running",
		slog.String("application", "exchanger"),
		slog.Int("port", app.config.port),
	)
	return srv.ListenAndServe()
}
