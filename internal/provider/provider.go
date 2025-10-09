package provider

import (
	"context"
	"database/sql"
	"log"

	"github.com/aglili/auction-app/internal/config"
	"github.com/aglili/auction-app/internal/events"
	"github.com/aglili/auction-app/internal/handlers"
	"github.com/aglili/auction-app/internal/repository"
	"github.com/aglili/auction-app/internal/scheduler"
	"github.com/aglili/auction-app/internal/service"
	"github.com/aglili/auction-app/internal/websocket"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Provider struct {
	DB             *sql.DB
	Validator      *validator.Validate
	UserHandler    *handlers.UserHandler
	HealthHandler  *handlers.HealthHandler
	AuctionHandler *handlers.AuctionHandler
	BidHandler     *handlers.BidHandler
	WsHandler      *handlers.WebSocketHandler
	Config         *config.Config
}

func NewProvider(config *config.Config, db *sql.DB, redis *redis.Client) *Provider {

	validator := validator.New()

	// resositories
	userRepository := repository.NewUserRepository(db)
	auctionRepository := repository.NewAuctionRepository(db)
	bidRepository := repository.NewBidRepository(db)

	wsConnManager := websocket.NewConnectionManager()

	publisher := events.NewEventPublisher(redis)
	subscriber := events.NewEventSubscriber(redis)

	// services
	paymentService := service.NewPaymentService(config)
	userService := service.NewUserService(userRepository)
	auctionService := service.NewAuctionService(auctionRepository)
	notificationService := service.NewNotificationService(userRepository, auctionRepository, wsConnManager,paymentService)
	bidService := service.NewBidService(bidRepository, auctionRepository, redis)
	

	// event handlers
	auctionEndedEventHandler := events.NewAuctionEventEndedHandler(notificationService, auctionRepository)
	outbidEventHandler := events.NewUserOutbidEventHandler(notificationService)

	// route handlers
	userHandler := handlers.NewUserHandler(userService, validator)
	auctionHandler := handlers.NewAuctionHandler(auctionService, validator)
	bidHandler := handlers.NewBidHandler(bidService, validator)
	wsHandler := handlers.NewWebSocketHandler(wsConnManager)
	healthHandler := handlers.NewHealthHandler()

	// scheduler
	scheduler := scheduler.NewAuctionScheduler(auctionRepository, redis, publisher)

	ctx := context.Background()
	if err := subscriber.Subscribe(ctx, events.EventAuctionEnded, events.EventUserOutbid); err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	go func() {
		handlers := map[string]events.EventHandler{
			events.EventAuctionEnded: auctionEndedEventHandler,
			events.EventUserOutbid:   outbidEventHandler,
		}
		if err := subscriber.Listen(ctx, handlers); err != nil {
			log.Printf("Event listener error: %v", err)
		}
	}()

	go scheduler.Start(ctx)


	return &Provider{
		HealthHandler:  healthHandler,
		Config:         config,
		UserHandler:    userHandler,
		DB:             db,
		AuctionHandler: auctionHandler,
		BidHandler:     bidHandler,
		WsHandler:      wsHandler,
	}
}
