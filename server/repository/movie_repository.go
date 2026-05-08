package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shingo/server/models"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]models.Movies, error) {
	q := `SELECT id, title, year, rated, released, runtime, genre, director,
	             writer, actors, plot, language, country, awards, poster,
	             metascore, imdb_rating, imdb_votes, imdb_id, type, images,
	             created_at, updated_at
	      FROM movies`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movies
	for rows.Next() {
		var m models.Movies
		if err := rows.Scan(
			&m.ID, &m.Title, &m.Year, &m.Rated, &m.Released,
			&m.Runtime, &m.Genre, &m.Director, &m.Writer, &m.Actors,
			&m.Plot, &m.Language, &m.Country, &m.Awards, &m.Poster,
			&m.Metascore, &m.ImdbRating, &m.ImdbVotes, &m.ImdbID,
			&m.Type, &m.Images, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, rows.Err()
}

func (r *MovieRepository) GetByID(ctx context.Context, id int64) (*models.Movies, error) {
	q := `SELECT id, title, year, rated, released, runtime, genre, director,
	             writer, actors, plot, language, country, awards, poster,
	             metascore, imdb_rating, imdb_votes, imdb_id, type, images,
	             created_at, updated_at
	      FROM movies WHERE id = $1`

	m := &models.Movies{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&m.ID, &m.Title, &m.Year, &m.Rated, &m.Released,
		&m.Runtime, &m.Genre, &m.Director, &m.Writer, &m.Actors,
		&m.Plot, &m.Language, &m.Country, &m.Awards, &m.Poster,
		&m.Metascore, &m.ImdbRating, &m.ImdbVotes, &m.ImdbID,
		&m.Type, &m.Images, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MovieRepository) Create(ctx context.Context, m *models.Movies) error {
	q := `INSERT INTO movies
		      (title, year, rated, released, runtime, genre, director, writer,
		       actors, plot, language, country, awards, poster, metascore,
		       imdb_rating, imdb_votes, imdb_id, type, images)
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
		  RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, q,
		m.Title, m.Year, m.Rated, m.Released, m.Runtime,
		m.Genre, m.Director, m.Writer, m.Actors, m.Plot,
		m.Language, m.Country, m.Awards, m.Poster, m.Metascore,
		m.ImdbRating, m.ImdbVotes, m.ImdbID, m.Type, m.Images,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}
