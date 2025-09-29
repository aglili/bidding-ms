package provider

import (
	"database/sql"

	"github.com/aglili/auction-app/internal/config"
	"github.com/aglili/auction-app/internal/handlers"
	"github.com/aglili/auction-app/internal/repository"
	"github.com/aglili/auction-app/internal/service"
)



type Provider struct{
	DB *sql.DB
	UserHandler *handlers.UserHandler
	HealthHandler *handlers.HealthHandler
	Config *config.Config
}



func NewProvider(config *config.Config, db *sql.DB) *Provider {
	// resositories
	userRepository := repository.NewUserRepository(db)


	// services
	userService := service.NewUserService(userRepository)

	userHandler := handlers.NewUserHandler(userService)
	healthHandler := handlers.NewHealthHandler()

	return &Provider{
		HealthHandler: healthHandler,
		Config: config,
		UserHandler: userHandler,
		DB: db,
	}
}