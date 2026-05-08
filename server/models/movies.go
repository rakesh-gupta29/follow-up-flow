package models

import (
	"time"
)

type Movies struct {
	ID         int64     `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	Year       string    `json:"year" db:"year"`
	Rated      string    `json:"rated" db:"rated"`
	Released   string    `json:"released" db:"released"`
	Runtime    string    `json:"runtime" db:"runtime"`
	Genre      string    `json:"genre" db:"genre"`
	Director   string    `json:"director" db:"director"`
	Writer     string    `json:"writer" db:"writer"`
	Actors     string    `json:"actors" db:"actors"`
	Plot       string    `json:"plot" db:"plot"`
	Language   string    `json:"language" db:"language"`
	Country    string    `json:"country" db:"country"`
	Awards     string    `json:"awards" db:"awards"`
	Poster     string    `json:"poster" db:"poster"`
	Metascore  string    `json:"metascore" db:"metascore"`
	ImdbRating string    `json:"imdbRating" db:"imdb_rating"`
	ImdbVotes  string    `json:"imdbVotes" db:"imdb_votes"`
	ImdbID     string    `json:"imdbID" db:"imdb_id"`
	Type       string    `json:"type" db:"type"`
	Images     []string  `json:"images" db:"images"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}
