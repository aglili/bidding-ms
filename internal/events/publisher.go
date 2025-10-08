package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type EventPublisher struct {
	client *redis.Client
}

func NewEventPublisher(client *redis.Client) *EventPublisher {
	return &EventPublisher{
		client: client,
	}
}

func (p *EventPublisher) PublishAuctionEnded(ctx context.Context, event AuctionEndedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal auction ended event: %w", err)
	}

	err = p.client.Publish(ctx, EventAuctionEnded, data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish auction ended event: %w", err)
	}

	log.Printf("Published auction ended event for auction %s", event.AuctionID)
	return nil
}

func (p *EventPublisher) PublishPlayerOutbid(ctx context.Context, event UserOutbidEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal user outbid event: %w", err)
	}

	err = p.client.Publish(ctx, EventUserOutbid, data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish user outbid event: %w", err)
	}

	log.Printf("Published outbid event for user %s on auction %s", event.OutbidUserID, event.AuctionID)
	return nil

}
