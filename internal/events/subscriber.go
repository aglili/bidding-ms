package events

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type EventSubscriber struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func NewEventSubscriber(client *redis.Client) *EventSubscriber {
	return &EventSubscriber{
		client: client,
	}
}

func (s *EventSubscriber) Subscribe(ctx context.Context, channels ...string) error {
	s.pubsub = s.client.Subscribe(ctx, channels...)

	_, err := s.pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	log.Printf("Subscribed to channels : %v", channels)
	return nil
}

func (s *EventSubscriber) Listen(ctx context.Context, handlers map[string]EventHandler) error {
	if s.pubsub == nil {
		return fmt.Errorf("not subscribed to any channels")
	}

	ch := s.pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping event listener")
			return s.pubsub.Close()
		case msg := <-ch:
			handler, exists := handlers[msg.Channel]
			if !exists {
				log.Printf("No handler for channel: %s", msg.Channel)
				continue
			}

			go func(m *redis.Message) {
				if err := handler.Handle(ctx, []byte(m.Payload)); err != nil {
					log.Printf("Error handling event on channel %s: %v", m.Channel, err)
				}
			}(msg)
		}
	}
}
