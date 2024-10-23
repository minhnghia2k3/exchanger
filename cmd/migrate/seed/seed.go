package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/minhnghia2k3/exchanger/internal/env"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"io"
	"log"
	"net/http"
)

type CurrencyData struct {
	Result         string     `json:"result"`
	SupportedCodes [][]string `json:"supported_codes"`
}

func getCurrencyData() *CurrencyData {
	var currencyData CurrencyData

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/codes", env.GetString("EXCHANGER_RATE_API", ""))
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &currencyData)
	if err != nil {
		log.Fatal(err)
	}

	return &currencyData
}

func seed(db *sql.DB, storage *store.Storage) {
	ctx := context.Background()
	currencyData := getCurrencyData()

	// Seeding currencies
	tx, _ := db.BeginTx(ctx, nil)

	for _, currency := range currencyData.SupportedCodes {
		c := &store.Currency{
			Code: currency[0],
			Name: currency[1],
		}

		if err := storage.Currencies.Insert(ctx, c); err != nil {
			_ = tx.Rollback()
			log.Println("failed to insert currency:", err)
			return
		}

		tx.Commit()
	}

	log.Println("Seeding successfully!")
}
