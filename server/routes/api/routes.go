package apiRoutes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/handlers/api"
	"github.com/shingo/server/utils/response"
)

func RegisterMovieRoutes(app fiber.Router, movieHandler *api.MovieHandler) {
	apiGroup := app.Group("/api/v1/movies")
	apiGroup.Get("/health", movieHandler.HealthCheck)

	apiGroup.Get("/", movieHandler.GetAllMovies)
	apiGroup.Get("/:id", movieHandler.GetMovieByID)

	apiGroup.Use(response.APINotFound)
}

func RegisterAppRoutes(app fiber.Router, appHandler *api.AppHandler) {
	apiGroup := app.Group("/api/v1")
	apiGroup.Get("/health", appHandler.HealthCheck)

	apiGroup.Use(response.APINotFound)
}
