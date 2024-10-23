package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type IExchangeRates interface {
	FindByCode(ctx context.Context, base, target string) (*ExchangeRate, error)
}

type ExchangeRate struct {
	ID               int64     `json:"id"`
	BaseCurrencyID   string    `json:"base_currency_id"`
	TargetCurrencyID string    `json:"target_currency_id"`
	Rate             float64   `json:"rate"`
	LastUpdate       time.Time `json:"last_update"`
	NextUpdate       time.Time `json:"next_update"`
	BaseCode         string    `json:"base_code"`
	TargetCode       string    `json:"target_code"`
}

type ExchangeRateStorage struct {
	db *sql.DB
}

func (s *ExchangeRateStorage) FindByCode(ctx context.Context, base, target string) (*ExchangeRate, error) {
	var rate ExchangeRate
	query := `
	SELECT er.id, base_currency_id, target_currency_id, rate, last_update, next_update, base.code AS base_code, target.code AS target_code
	FROM exchange_rates er
	INNER JOIN currencies base ON base.id = er.base_currency_id
	INNER JOIN currencies target ON target.id = er.target_currency_id
	WHERE base.code = $1 AND target.code = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, base, target).Scan(
		&rate.ID,
		&rate.BaseCurrencyID,
		&rate.TargetCurrencyID,
		&rate.Rate,
		&rate.LastUpdate,
		&rate.NextUpdate,
		&rate.BaseCode,
		&rate.TargetCode,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &rate, nil
}
