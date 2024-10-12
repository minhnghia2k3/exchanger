package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found record")
	ErrConflict = errors.New("record already exists")
)

var (
	QueryContextTimeout = 3 * time.Second
)

type Storage struct {
	Currencies ICurrencies
	Users      IUser
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Currencies: &CurrencyModel{db: db},
		Users:      &UserModel{db: db},
	}
}

func withTx(ctx context.Context, db *sql.DB, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
