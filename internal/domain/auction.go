package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuctionImage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AuctionID uuid.UUID `json:"auction_id" db:"auction_id"`
	ImageURL  string    `json:"image_url" db:"image_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Auction struct {
	ID            uuid.UUID `json:"id" db:"id"`
	SellerID      uuid.UUID `json:"seller_id,omitempty" db:"seller_id"`
	Title         string    `json:"title" db:"title"`
	Description   *string   `json:"description,omitempty" db:"description"`
	StartingPrice float64   `json:"starting_price" db:"starting_price"`
	CurrentPrice  float64   `json:"current_price" db:"current_price"`
	Status        string    `json:"status" db:"status"` // open || closed || cancelled
	StartTime     time.Time `json:"start_time" db:"start_time"`
	EndTime       time.Time `json:"end_time" db:"end_time"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	Images        []string  `json:"images,omitempty" db:"-"`
}

type AuctionResponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   *string   `json:"description,omitempty"`
	StartingPrice float64   `json:"starting_price"`
	CurrentPrice  float64   `json:"current_price"`
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Images        []string  `json:"images,omitempty"`
}

type AuctionRepository interface {
	CreateAuction(ctx context.Context, auction *Auction, sellerID uuid.UUID, imageURLs []string) (*Auction, error)
	GetAuction(ctx context.Context, auctionID uuid.UUID) (*Auction, error)
	GetUserAuctions(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Auction, int, error)
	UpdateCurrentPrice(ctx context.Context, auctionID uuid.UUID, amount float64) error
	CloseAuction(ctx context.Context, auctionID uuid.UUID) error
	GetEndedActiveAuctions(ctx context.Context, currentTime time.Time) ([]*Auction, error)
	GetOpenAuctions(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Auction, int, error)
}

type AuctionService interface {
	CreateAuction(ctx context.Context, auction *Auction, sellerID uuid.UUID, imageURLs []string) (*Auction, error)
	GetAuction(ctx context.Context, auctionID uuid.UUID) (*Auction, error)
	GetUserAuctions(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Auction, int, error)
	GetOpenAuctions(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Auction, int, error)
}
