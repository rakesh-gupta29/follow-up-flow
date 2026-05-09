package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodHead,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{"*"},
	}))
	cfg := config.AppConfig

	db, err := database.New(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 3. NEW: Initialize Auth Logic
	authRepo := repository.NewAdminRepository(db.Client)
	authHandler := api.NewAuthHandler(authRepo)
	apiRoutes.RegisterAdminRoutes(server, authHandler)

	contactsRepo := repository.NewContactsRepository(db.Client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := contactsRepo.EnsureCollection(ctx); err != nil {
		cancel()
		log.Fatalf("failed to initialize contacts collection: %v", err)
	}
	cancel()
	contactsHandler := api.NewContactsHandler(contactsRepo)
	apiRoutes.RegisterContactsRoutes(server, contactsHandler)

	campaignsRepo := repository.NewCampaignsRepository(db.Client)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err := campaignsRepo.EnsureCollection(ctx); err != nil {
		cancel()
		log.Fatalf("failed to initialize campaigns collection: %v", err)
	}
	cancel()
	campaignsHandler := api.NewCampaignsHandler(campaignsRepo)
	apiRoutes.RegisterCampaignRoutes(server, campaignsHandler)

	// 4. Initialize App General Logic
	appHandler := api.NewAppHandler()
	apiRoutes.RegisterAppRoutes(server, appHandler)

	return &Application{
		Server: server,
		DB:     db,
		Config: cfg,
	}
}
