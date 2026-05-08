package bootstrap

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/config"
	"github.com/shingo/server/database"
	"github.com/shingo/server/handlers/api"
	"github.com/shingo/server/repository"
	apiRoutes "github.com/shingo/server/routes/api"
)

type Application struct {
	Server *fiber.App
	DB     *database.DB
	Config config.Config
}

func New() *Application {
	server := fiber.New()
	cfg := config.AppConfig

	db, err := database.New(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	movieHandler := api.NewMovieHandler(repository.NewMovieRepository(db.Pool))
	apiRoutes.RegisterMovieRoutes(server, movieHandler)

	appHandler := api.NewAppHandler()
	apiRoutes.RegisterAppRoutes(server, appHandler)

	return &Application{
		Server: server,
		DB:     db,
		Config: cfg,
	}
}
