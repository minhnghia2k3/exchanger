package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := CurrencyModel{db}

	testCases := []struct {
		name          string
		id            int64
		mockResponse  func()
		expectedError error
		expectedData  *Currency
	}{
		{
			name: "should return a currency successfully",
			id:   1,
			mockResponse: func() {
				const id = 1

				// Add a row to mock database
				rows := sqlmock.NewRows([]string{"id", "code", "name", "symbol_url"}).
					AddRow(id, "USD", "US Dollar", nil)

				mock.ExpectQuery(`SELECT id, code, name, symbol_url FROM currencies WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedError: nil,
			expectedData:  &Currency{ID: 1, Code: "USD", Name: "US Dollar", SymbolUrl: nil},
		},
		{
			name: "should return an error not found",
			id:   3,
			mockResponse: func() {
				const id = 3
				mock.ExpectQuery(`SELECT id, code, name, symbol_url FROM currencies WHERE id = \$1`).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: fmt.Errorf("%w by id %d", ErrNotFound, 3),
			expectedData:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockResponse()

			ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
			defer cancel()

			currency, err := model.Get(ctx, tc.id)

			// error: not found
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedData, currency)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCurrencyModel_List(t *testing.T) {
	// Initialize mock database and sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize CurrencyModel with mock db
	model := CurrencyModel{db: db}

	// Define test cases
	tests := []struct {
		name           string
		searchTerm     string
		mockRows       *sqlmock.Rows
		mockError      error
		expectedResult []Currency
		expectedMeta   Metadata
		expectError    bool
	}{
		{
			name:       "Found multiple currencies",
			searchTerm: "USD",
			mockRows: sqlmock.NewRows([]string{"count", "id", "code", "name", "symbol_url"}).
				AddRow(2, 1, "USD", "US Dollar", "https://example.com/usd-symbol.png").
				AddRow(2, 2, "EUR", "Euro", "https://example.com/eur-symbol.png"),
			expectedResult: []Currency{
				{ID: 1, Code: "USD", Name: "US Dollar", SymbolUrl: ptr("https://example.com/usd-symbol.png")},
				{ID: 2, Code: "EUR", Name: "Euro", SymbolUrl: ptr("https://example.com/eur-symbol.png")},
			},
			expectedMeta: Metadata{
				CurrentPage: 1,
				PageSize:    20,
				FirstPage:   1,
				LastPage:    1,
				TotalRecord: 2,
			},
			expectError: false,
		},
		{
			name:           "No currencies found",
			searchTerm:     "JPY",
			mockRows:       sqlmock.NewRows([]string{"count", "id", "code", "name", "symbol_url"}), // No results
			expectedResult: []Currency(nil),
			expectedMeta: Metadata{
				CurrentPage: 0,
				PageSize:    0,
				FirstPage:   0,
				LastPage:    0,
				TotalRecord: 0,
			},
			expectError: false,
		},
		{
			name:        "Database error",
			searchTerm:  "USD",
			mockRows:    nil,
			mockError:   assert.AnError,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock query expectations
			query := `SELECT COUNT\(\*\) OVER\(\), id, code, name, symbol_url FROM currencies
	WHERE \(to_tsvector\('simple', name\) @@ plainto_tsquery\('simple', \$1\) 
	OR code = \$1 OR \$1 = ''\)
	ORDER BY id ASC, id ASC
	LIMIT \$2 OFFSET \$3`

			// Mock rows and error handling
			if tc.mockError == nil {
				mock.ExpectQuery(query).
					WithArgs(tc.searchTerm, 20, 0). // Simulate page size and offset
					WillReturnRows(tc.mockRows)
			} else {
				mock.ExpectQuery(query).
					WithArgs(tc.searchTerm, 20, 0).
					WillReturnError(tc.mockError)
			}

			// Set up context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			// Call the List method
			filter := Filter{
				Page:         1,
				PageSize:     20,
				Sort:         "",
				SortSafeList: nil,
				Search:       tc.searchTerm,
			}
			currencies, meta, err := model.List(ctx, filter)

			// Assert the result
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, currencies)
				assert.Equal(t, tc.expectedMeta, meta)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInsertCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := CurrencyModel{db}

	// Define the currency you want to insert
	currency := &Currency{
		Code:      "USD",
		Name:      "US Dollar",
		SymbolUrl: ptr("https://example.com/usd-symbol.png"),
	}

	// Mock successful query
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO currencies`).
		WithArgs(currency.Code, currency.Name, currency.SymbolUrl).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.Background()

	err = model.Insert(ctx, currency)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := CurrencyModel{db}

	// Define the currency you want to update
	currency := &Currency{
		Code:      "USD",
		Name:      "US Dollar updated",
		SymbolUrl: ptr("https://example.com/usd-symbol.png"),
	}

	id := int64(1)

	// Mock successful update
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE currencies SET").
		WithArgs(currency.Code, currency.Name, currency.SymbolUrl, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.Background()

	// Call the update method
	err = model.Update(ctx, id, currency)

	// Assert no errors
	assert.NoError(t, err)

	// Check that expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCurrencyNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := CurrencyModel{db}

	// Define the currency to update
	currency := &Currency{
		Code:      "USD",
		Name:      "US Dollar",
		SymbolUrl: ptr("https://example.com/usd-symbol.png"),
	}

	id := int64(1)

	// Mock transaction and query execution
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE currencies SET").
		WithArgs(currency.Code, currency.Name, currency.SymbolUrl, id).
		WillReturnError(ErrNotFound) // 0 rows affected
	mock.ExpectRollback() // Expect a rollback since no rows were found

	ctx := context.Background()

	// Call the update method
	err = model.Update(ctx, id, currency)

	// Assert that ErrNotFound is returned when no rows are affected
	assert.ErrorIs(t, err, ErrNotFound)

	// Check that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := CurrencyModel{db}

	id := int64(1)

	// Mock successful delete
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM currencies").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.Background()

	// Call the delete method
	err = model.Delete(ctx, id)

	// Assert no errors
	assert.NoError(t, err)

	// Check that expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCurrencyNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := CurrencyModel{db}

	id := int64(1)

	// Mock delete failure with ErrNoRows
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM currencies").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	ctx := context.Background()

	// Call the delete method
	err = model.Delete(ctx, id)

	// Assert the expected error (ErrNotFound)
	assert.ErrorIs(t, err, ErrNotFound)

	// Check that expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Helper function to return a pointer to a string
func ptr(s string) *string {
	return &s
}
