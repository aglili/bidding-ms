package provider

import (
	"database/sql"

	"github.com/aglili/auction-app/internal/config"
	"github.com/aglili/auction-app/internal/handlers"
	"github.com/aglili/auction-app/internal/repository"
	"github.com/aglili/auction-app/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Provider struct {
	DB             *sql.DB
	Validator   *validator.Validate
	UserHandler    *handlers.UserHandler
	HealthHandler  *handlers.HealthHandler
	AuctionHandler *handlers.AuctionHandler
	BidHandler *handlers.BidHandler
	Config         *config.Config
}

func NewProvider(config *config.Config, db *sql.DB, redis *redis.Client) *Provider {

	validator := validator.New()


	// resositories
	userRepository := repository.NewUserRepository(db)
	auctionRepository := repository.NewAuctionRepository(db)
	bidRepository := repository.NewBidRepository(db)

	// services
	userService := service.NewUserService(userRepository)
	auctionService := service.NewAuctionService(auctionRepository)
	bidService := service.NewBidService(bidRepository,auctionRepository,redis)

	userHandler := handlers.NewUserHandler(userService,validator)
	auctionHandler := handlers.NewAuctionHandler(auctionService,validator)
	bidHandler := handlers.NewBidHandler(bidService,validator)
	healthHandler := handlers.NewHealthHandler()

	return &Provider{
		HealthHandler:  healthHandler,
		Config:         config,
		UserHandler:    userHandler,
		DB:             db,
		AuctionHandler: auctionHandler,
		BidHandler: bidHandler,
	}
}
