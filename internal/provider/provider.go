package provider

import (
	"database/sql"

	"github.com/aglili/auction-app/internal/config"
	"github.com/aglili/auction-app/internal/handlers"
	"github.com/aglili/auction-app/internal/repository"
	"github.com/aglili/auction-app/internal/service"
	"github.com/go-playground/validator/v10"
)

type Provider struct {
	DB             *sql.DB
	Validator   *validator.Validate
	UserHandler    *handlers.UserHandler
	HealthHandler  *handlers.HealthHandler
	AuctionHandler *handlers.AuctionHandler
	Config         *config.Config
}

func NewProvider(config *config.Config, db *sql.DB) *Provider {

	validator := validator.New()


	// resositories
	userRepository := repository.NewUserRepository(db)
	auctionRepository := repository.NewAuctionRepository(db)

	// services
	userService := service.NewUserService(userRepository)
	auctionService := service.NewAuctionService(auctionRepository)

	userHandler := handlers.NewUserHandler(userService,validator)
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
