package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/events"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuctionScheduler struct {
	auctionRepo domain.AuctionRepository
	cache       *redis.Client
	publisher   *events.EventPublisher
}

func NewAuctionScheduler(auctionRepo domain.AuctionRepository, cache *redis.Client, publisher *events.EventPublisher) *AuctionScheduler {
	return &AuctionScheduler{
		auctionRepo: auctionRepo,
		cache:       cache,
		publisher:   publisher,
	}
}

func (s *AuctionScheduler) Start(ctx context.Context) {
	log.Println("Starting auction scheduler")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	s.checkAndCloseAuctions(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping auction scheduler")
			return
		case <-ticker.C:
			s.checkAndCloseAuctions(ctx)
		}
	}

}

func (s *AuctionScheduler) checkAndCloseAuctions(ctx context.Context) {
	// Get all active auctions that have ended
	auctions, err := s.auctionRepo.GetEndedActiveAuctions(ctx, time.Now())
	if err != nil {
		log.Printf("Error fetching ended auctions: %v", err)
		return
	}

	log.Printf("Found %d auctions to close", len(auctions))

	for _, auction := range auctions {
		if err := s.closeAuction(ctx, auction); err != nil {
			log.Printf("Error closing auction %s: %v", auction.ID, err)
		}
	}
}

func (s *AuctionScheduler) closeAuction(ctx context.Context, auction *domain.Auction) error {
	log.Printf("Closing auction %s", auction.ID)

	// Get winner from Redis
	bidderKey := fmt.Sprintf("auction:%s:highest_bidder", auction.ID.String())
	priceKey := fmt.Sprintf("auction:%s:highest_bid", auction.ID.String())

	winnerIDStr, err := s.cache.Get(ctx, bidderKey).Result()
	if err == redis.Nil {
		return s.auctionRepo.CloseAuction(ctx, auction.ID)
	} else if err != nil {
		return fmt.Errorf("failed to get winner: %w", err)
	}

	winnerID, err := uuid.Parse(winnerIDStr)
	if err != nil {
		return fmt.Errorf("invalid winner ID: %w", err)
	}

	finalPriceStr, err := s.cache.Get(ctx, priceKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get final price: %w", err)
	}

	var finalPrice float64
	fmt.Sscanf(finalPriceStr, "%f", &finalPrice)

	if err := s.auctionRepo.CloseAuction(ctx, auction.ID); err != nil {
		return fmt.Errorf("failed to update auction status: %w", err)
	}

	event := events.AuctionEndedEvent{
		AuctionID:  auction.ID,
		WinnerID:   winnerID,
		FinalPrice: finalPrice,
		EndedAt:    time.Now(),
	}

	if err := s.publisher.PublishAuctionEnded(ctx, event); err != nil {
		log.Printf("Failed to publish auction ended event: %v", err)
	}

	s.cache.Del(ctx, bidderKey, priceKey)

	log.Printf("Auction %s closed. Winner: %s, Price: %.2f", auction.ID, winnerID, finalPrice)
	return nil
}
