// migrations/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	direction := flag.String("direction", "up", "up | down")
	flag.Parse()

	if err := godotenv.Load(".env.dev"); err != nil {
		log.Fatal("could not load env file")
	}

	dbURL := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("failed to init migrate: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate up failed: %v", err)
		}
		log.Println("migrations applied")
	case "down":
		if err = m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate down failed: %v", err)
		}
		log.Println("migrations rolled back")
	default:
		log.Fatalf("unknown direction: %s", *direction)
	}
}