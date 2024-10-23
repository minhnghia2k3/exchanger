package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func readInt(r *http.Request, key string, fallback int) int {
	val := r.URL.Query().Get(key)

	if val == "" {
		return fallback
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valInt
}

func readString(r *http.Request, key string, fallback string) string {
	val := r.URL.Query().Get(key)

	if val == "" {
		return fallback
	}

	return val
}

func SHA256Hash(text string) string {
	h := sha256.Sum256([]byte(text))

	return hex.EncodeToString(h[:])
}

func GetAPI(url string, data any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		return fmt.Errorf("%w, response failed with status code: %d and\nbody: %s\n", errUnsupportedCurrencyCode, resp.StatusCode, body)
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	return nil
}
