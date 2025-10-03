package domain

import (
	"context"

	"github.com/google/uuid"
)

type Bid struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AuctionID uuid.UUID `json:"auction_id" db:"auction_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Amount    float64   `json:"amount" db:"amount"`
}

type BidRepository interface {
	CreateBid(ctx context.Context, auctionID, userID uuid.UUID, amount float64) error
}

type BidService interface {
	CreateBid(ctx context.Context, auctionID, userID uuid.UUID, amount float64) error
}
