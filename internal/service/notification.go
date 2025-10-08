package service

import (
	"context"
	"fmt"
	"log"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/websocket"
	"github.com/google/uuid"
)

type NotificationService struct {
	userRepo    domain.UserRepository
	auctionRepo domain.AuctionRepository
	connManager *websocket.ConnectionManager
}

func NewNotificationService(userRepo domain.UserRepository, auctionRepo domain.AuctionRepository, connManager *websocket.ConnectionManager) *NotificationService {
	return &NotificationService{
		userRepo:    userRepo,
		connManager: connManager,
		auctionRepo: auctionRepo,
	}
}

func (s *NotificationService) NotifyAuctionWon(ctx context.Context, userID, auctionID uuid.UUID, price float64) error {
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	auction, err := s.auctionRepo.GetAuction(ctx, auctionID)
	if err != nil {
		return fmt.Errorf("failed to get auction: %w", err)
	}

	message := websocket.NotificationMessage{
		Type: "auction_won",
		Payload: map[string]any{
			"auction_id":  auction.ID,
			"title":       auction.Title,
			"description": auction.Description,
			"price":       price,
			"message":     fmt.Sprintf("Congratulations! You won the auction for $%.2f", price),
		},
	}

	if err := s.connManager.SendToUser(userID, message); err != nil {
		log.Printf("Failed to send WebSocket notification: %v", err)

	}

	return nil
}

func (s *NotificationService) NotifyOutbid(ctx context.Context, userID, auctionID uuid.UUID, newPrice float64) error {
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	auction, err := s.auctionRepo.GetAuction(ctx, auctionID)
	if err != nil {
		return fmt.Errorf("failed to get auction: %w", err)
	}

	message := websocket.NotificationMessage{
		Type: "outbid",
		Payload: map[string]interface{}{
			"auction_id": auctionID,
			"title":      auction.Title,
			"new_bid":    newPrice,
			"message":    fmt.Sprintf("You've been outbid! Current bid is $%.2f", newPrice),
		},
	}

	if err := s.connManager.SendToUser(userID, message); err != nil {
		log.Printf("Failed to send WebSocket notification: %v", err)
	}

	return nil
}
