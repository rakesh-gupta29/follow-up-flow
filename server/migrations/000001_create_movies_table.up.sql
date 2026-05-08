-- migrations/000001_create_movies_table.up.sql
CREATE TABLE IF NOT EXISTS movies (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    year VARCHAR(10),
    rated VARCHAR(20),
    released VARCHAR(50),
    runtime VARCHAR(20),
    genre VARCHAR(255),
    director VARCHAR(255),
    writer TEXT,
    actors TEXT,
    plot TEXT,
    language VARCHAR(100),
    country VARCHAR(100),
    awards TEXT,
    poster TEXT,
    metascore VARCHAR(10),
    imdb_rating VARCHAR(10),
    imdb_votes VARCHAR(20),
    imdb_id VARCHAR(20) UNIQUE,
    type VARCHAR(20),
    images TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_movies_imdb_id ON movies(imdb_id);
CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_year ON movies(year);