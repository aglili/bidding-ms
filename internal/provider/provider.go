package provider

import (
	"database/sql"

	"github.com/aglili/auction-app/internal/config"
	"github.com/aglili/auction-app/internal/handlers"
	"github.com/aglili/auction-app/internal/repository"
	"github.com/aglili/auction-app/internal/service"
)

type Provider struct {
	DB             *sql.DB
	UserHandler    *handlers.UserHandler
	HealthHandler  *handlers.HealthHandler
	AuctionHandler *handlers.AuctionHandler
	Config         *config.Config
}

func NewProvider(config *config.Config, db *sql.DB) *Provider {
	// resositories
	userRepository := repository.NewUserRepository(db)
	auctionRepository := repository.NewAuctionRepository(db)

	// services
	userService := service.NewUserService(userRepository)
	auctionService := service.NewAuctionService(auctionRepository)

	userHandler := handlers.NewUserHandler(userService)
	auctionHandler := handlers.NewAuctionHandler(auctionService)
	healthHandler := handlers.NewHealthHandler()

	return &Provider{
		HealthHandler:  healthHandler,
		Config:         config,
		UserHandler:    userHandler,
		DB:             db,
		AuctionHandler: auctionHandler,
	}
}
