package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func connectDB(cfg dbConfig) (*sql.DB, error) {
	pool, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	pool.SetMaxOpenConns(cfg.maxOpenConn)
	pool.SetMaxIdleConns(cfg.maxIdleConn)
	duration, err := time.ParseDuration(cfg.maxIdleTime)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = pool.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to database")
	return pool, nil
}
