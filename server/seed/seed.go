package seed

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shingo/server/models"
	"github.com/shingo/server/repository"
)

func Run(pool *pgxpool.Pool) {
	seedMovies(pool)
}

func seedMovies(pool *pgxpool.Pool) {
	var seeds []models.Movies
	readJSON("../../seed/data/movies.json", &seeds)

	repo := repository.NewMovieRepository(pool)
	for _, m := range seeds {
		if err := repo.Create(context.Background(), &m); err != nil {
			log.Printf("seed movies: skipping %s — %v", m.Title, err)
			continue
		}
		log.Printf("seed movies: inserted %s", m.Title)
	}
}

func readJSON(path string, dest any) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("seed: could not read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		log.Fatalf("seed: could not parse %s: %v", path, err)
	}
}
