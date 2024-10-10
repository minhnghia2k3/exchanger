package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found record")
)

var (
	QueryContextTimeout = 3 * time.Second
)

type Storage struct {
	Currencies ICurrencies
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Currencies: &CurrencyModel{db: db},
	}
}

func withTx(ctx context.Context, db *sql.DB, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		// Rollback in case anything fails in the closure
		err = tx.Rollback()
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
