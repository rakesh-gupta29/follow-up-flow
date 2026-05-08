package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shingo/server/database"
	"github.com/shingo/server/seed"
)

func main() {
	if err := godotenv.Load("../../.env.dev"); err != nil {
		log.Fatal("could not load env file")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := database.New(dbURL)
	if err != nil {
		log.Fatalf("seed: db connection failed: %v", err)
	}
	defer db.Close()

	seed.Run(db.Pool)
}
