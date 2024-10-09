package model

import "database/sql"

type Currency struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type CurrencyModel struct {
	db *sql.DB
}

func (m *CurrencyModel) Get(id int64) (*Currency, error) {
	return nil, nil
}
func (m *CurrencyModel) List() ([]Currency, error)                 { return nil, nil }
func (m *CurrencyModel) Insert(currency *Currency) error           { return nil }
func (m *CurrencyModel) Update(id int64, currency *Currency) error { return nil }
func (m *CurrencyModel) Delete(id int64) error                     { return nil }
