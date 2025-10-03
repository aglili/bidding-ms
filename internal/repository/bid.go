package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type BidRepository struct {
	db *sql.DB
}

func NewBidRepository(db *sql.DB) *BidRepository {
	return &BidRepository{
		db: db,
	}
}

func (r *BidRepository) CreateBid(ctx context.Context, auctionID, userID uuid.UUID, amount float64) error {
	query := `INSERT INTO bids
	(auction_id,bidder_id,amount) 
	VALUES ($1,$2,$3)
	`

	_, err := r.db.ExecContext(ctx, query, auctionID, userID, amount)

	return err

}
