package model

import "database/sql"

type Models struct {
	Currencies CurrencyModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Currencies: CurrencyModel{db: db},
	}
}
