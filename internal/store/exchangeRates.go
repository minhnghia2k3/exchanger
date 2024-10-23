package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type IExchangeRates interface {
	GetByPair(ctx context.Context, base, target string) (*ExchangeRate, error)
	Save(ctx context.Context, rate *ExchangeRate) error
	Update(ctx context.Context, rate *ExchangeRate) error
}

type ExchangeRate struct {
	ID         int64     `json:"id"`
	LastUpdate time.Time `json:"last_update"`
	NextUpdate time.Time `json:"next_update"`
	BaseCode   string    `json:"base_code"`
	TargetCode string    `json:"target_code"`
	Rate       float64   `json:"rate"`
}

type ExchangeRateStorage struct {
	db *sql.DB
}

func (s *ExchangeRateStorage) GetByPair(ctx context.Context, base, target string) (*ExchangeRate, error) {
	var rate ExchangeRate
	query := `
	SELECT er.id, rate, last_update, next_update, base.code AS base_code, target.code AS target_code
	FROM exchange_rates er
	INNER JOIN currencies base ON base.code = er.base_code
	INNER JOIN currencies target ON target.code = er.target_code
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

func (s *ExchangeRateStorage) Save(ctx context.Context, rate *ExchangeRate) error {
	query := `
	INSERT INTO exchange_rates(base_code, target_code, rate, last_update, next_update)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	args := []any{rate.BaseCode, rate.TargetCode, rate.Rate, rate.LastUpdate, rate.NextUpdate}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&rate.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExchangeRateStorage) Update(ctx context.Context, rate *ExchangeRate) error {
	query := `
	UPDATE exchange_rates r
	SET	rate = $1
	WHERE r.base_code = $2 AND r.target_code = $3
`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, rate.Rate, rate.BaseCode, rate.TargetCode)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
