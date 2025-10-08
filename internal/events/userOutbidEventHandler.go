package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type UserOutbidHandler struct {
	notificationService NotificationService
}

func NewUserOutbidEventHandler(notificationService NotificationService) *UserOutbidHandler {
	return &UserOutbidHandler{
		notificationService: notificationService,
	}
}

func (h *UserOutbidHandler) Handle(ctx context.Context, data []byte) error {
	var event UserOutbidEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user outbid event: %w", err)
	}

	log.Printf("User %s outbid on auction %s: %.2f -> %.2f",
		event.OutbidUserID, event.AuctionID, event.OldBid, event.NewBid)

	// Send notification to outbid user
	if err := h.notificationService.NotifyOutbid(ctx, event.OutbidUserID, event.AuctionID, event.NewBid); err != nil {
		return fmt.Errorf("failed to notify outbid user: %w", err)
	}

	return nil
}
