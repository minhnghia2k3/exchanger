package main

import (
	"github.com/joho/godotenv"
	"github.com/minhnghia2k3/exchanger/internal/database"
	"github.com/minhnghia2k3/exchanger/internal/env"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dsn := env.GetString("DATABASE_URL", "postgres://root:secret@localhost:5432/exchanger?sslmode=disable")

	db, err := database.ConnectDB(dsn, 3, 3, "15m")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storage := store.NewStorage(db)

	seed(db, storage)
}
