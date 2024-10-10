package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ICurrencies interface {
	Get(ctx context.Context, id int64) (*Currency, error)
	List(ctx context.Context, code, name string, filter Filter) ([]Currency, Metadata, error)
	Insert(ctx context.Context, currency *Currency) error
	Update(ctx context.Context, id int64, currency *Currency) error
	Delete(ctx context.Context, id int64) error
}

type Currency struct {
	ID        int64   `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	SymbolUrl *string `json:"symbol_url"`
}

type CurrencyModel struct {
	db *sql.DB
}

func (m *CurrencyModel) Get(ctx context.Context, id int64) (*Currency, error) {
	query := `SELECT id, code, name, symbol_url FROM currencies WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	var currency Currency

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&currency.ID,
		&currency.Code,
		&currency.Name,
		&currency.SymbolUrl,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &currency, nil
}
func (m *CurrencyModel) List(ctx context.Context, code, name string, filter Filter) ([]Currency, Metadata, error) {
	var currencies []Currency
	var totalRecord int

	query := fmt.Sprintf(`SELECT COUNT(*) OVER(), id, code, name, symbol_url FROM currencies
	WHERE (to_tsvector('simple', name) @@ to_tsquery('simple', $1) OR $1 = '') OR (code = $2 OR $2 = '')
	ORDER BY %s %s, id ASC
    LIMIT $3 OFFSET $4`, filter.sortColumn(), filter.sortDirection())

	fmt.Println("query:", query)

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	args := []any{code, name, filter.limit(), filter.calculateOffset()}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var currency Currency

		err = rows.Scan(&totalRecord, &currency.ID, &currency.Code, &currency.Name, &currency.SymbolUrl)
		if err != nil {
			return nil, Metadata{}, err
		}

		currencies = append(currencies, currency)
	}

	metadata := filter.calculateMetadata(totalRecord)

	return currencies, metadata, nil
}
func (m *CurrencyModel) Insert(ctx context.Context, currency *Currency) error {
	query := `INSERT INTO currencies(code, name, symbol_url) VALUES($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, currency.Code, currency.Name, currency.SymbolUrl)

	if err != nil {
		return err
	}

	return nil
}

func (m *CurrencyModel) Update(ctx context.Context, id int64, currency *Currency) error {
	query := `UPDATE currencies SET code = $1, name = $2, symbol_url = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, currency.Code, currency.Name, currency.SymbolUrl, id)

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
func (m *CurrencyModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM currencies WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)

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