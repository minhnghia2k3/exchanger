package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func readInt(r *http.Request, key string, fallback int) int {
	val := chi.URLParam(r, key)

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
	val := chi.URLParam(r, key)

	if val == "" {
		return fallback
	}

	return val
}
