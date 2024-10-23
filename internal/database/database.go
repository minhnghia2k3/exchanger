package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func ConnectDB(dsn string, maxIdleConn, maxOpenConn int, maxIdleTime string) (*sql.DB, error) {
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	pool.SetMaxOpenConns(maxIdleConn)
	pool.SetMaxIdleConns(maxOpenConn)

	duration, err := time.ParseDuration(maxIdleTime)
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
