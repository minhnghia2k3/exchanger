package main

import (
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
