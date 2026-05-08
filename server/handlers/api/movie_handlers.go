package api

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/shingo/server/repository"
	"github.com/shingo/server/utils/response"
)

type MovieHandler struct {
	repo *repository.MovieRepository
}

func NewMovieHandler(repo *repository.MovieRepository) *MovieHandler {
	return &MovieHandler{repo: repo}
}

func (h *MovieHandler) HealthCheck(c fiber.Ctx) error {
	return response.Send(c, "Movies working fine")
}

func (h *MovieHandler) GetAllMovies(c fiber.Ctx) error {
	movies, err := h.repo.GetAll(c.Context())
	if err != nil {
		return response.InternalError(c, "failed to fetch movies")
	}
	return response.Send(c, movies)
}

func (h *MovieHandler) GetMovieByID(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	movie, err := h.repo.GetByID(c.Context(), id)
	if err == pgx.ErrNoRows {
		return response.NotFound(c, "movie not found")
	}
	if err != nil {
		return response.InternalError(c, "failed to fetch movie")
	}

	return response.Send(c, movie)
}
