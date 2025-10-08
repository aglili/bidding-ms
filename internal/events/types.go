package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	EventAuctionEnded = "auction:ended"
	EventUserOutbid   = "user:outbid"
)

type AuctionEndedEvent struct {
	AuctionID  uuid.UUID `json:"auction_id"`
	WinnerID   uuid.UUID `json:"winner_id"`
	FinalPrice float64   `json:"final_price"`
	EndedAt    time.Time `json:"ended_at"`
}

type UserOutbidEvent struct {
	AuctionID    uuid.UUID `json:"auction_id"`
	OutbidUserID uuid.UUID `json:"outbid_user_id"`
	OldBid       float64   `json:"old_bid"`
	NewBid       float64   `json:"new_bid"`
	NewBidderID  uuid.UUID `json:"new_bidder_id"`
	OutbidAt     time.Time `json:"outbid_at"`
}

type EventHandler interface {
	Handle(ctx context.Context, data []byte) error
}

type NotificationService interface {
	NotifyAuctionWon(ctx context.Context, userID, auctionID uuid.UUID, price float64) error
	NotifyOutbid(ctx context.Context, userID, auctionID uuid.UUID, newBid float64) error
}
