package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aglili/auction-app/internal/domain"
)

type AuctionEventEndedHandler struct {
	notificationService NotificationService
	auctionRepo         domain.AuctionRepository
}

func NewAuctionEventEndedHandler(notificationService NotificationService, auctionRepo domain.AuctionRepository) *AuctionEventEndedHandler {
	return &AuctionEventEndedHandler{
		notificationService: notificationService,
		auctionRepo:         auctionRepo,
	}
}

func (h *AuctionEventEndedHandler) Handle(ctx context.Context, data []byte) error {
	var event AuctionEndedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal auction ended event: %w", err)
	}

	if err := h.notificationService.NotifyAuctionWon(ctx, event.WinnerID, event.AuctionID, event.FinalPrice); err != nil {
		return fmt.Errorf("failed to notify winner: %w", err)
	}

	if err := h.auctionRepo.CloseAuction(ctx, event.AuctionID); err != nil {
		return fmt.Errorf("failed to close auction: %w", err)
	}

	return nil
}
