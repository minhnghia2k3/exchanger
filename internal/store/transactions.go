package store

import (
	"context"
	"database/sql"
	"time"
)

type ITransaction interface {
	Save(ctx context.Context, transaction *Transaction) error
}

type Transaction struct {
	ID              int64     `json:"id"`
	UserID          *int64    `json:"user_id,omitempty"`
	BaseCode        string    `json:"base_code"`
	TargetCode      string    `json:"target_code"`
	ConvertedAmount float64   `json:"converted_amount"`
	ConvertedRate   float64   `json:"converted_rate"`
	Result          float64   `json:"result"`
	CreatedAt       time.Time `json:"created_at"`
}

type TransactionStorage struct {
	db *sql.DB
}

func (s *TransactionStorage) Save(ctx context.Context, transaction *Transaction) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		query := `
	INSERT INTO transactions(user_id, base_code, target_code, converted_amount, converted_rate, result)
	VALUES($1, $2, $3, $4, $5, $6)`

		ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
		defer cancel()

		args := []any{transaction.UserID, transaction.BaseCode, transaction.TargetCode, transaction.ConvertedAmount,
			transaction.ConvertedAmount, transaction.ConvertedRate, transaction.Result}

		err := tx.QueryRowContext(ctx, query, args...).Scan(&transaction.ID, &transaction.CreatedAt)

		if err != nil {
			return err
		}

		return nil
	})
}
