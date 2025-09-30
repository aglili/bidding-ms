package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/google/uuid"
)

type AuctionRepository struct {
	db *sql.DB
}

func NewAuctionRepository(db *sql.DB) *AuctionRepository {
	return &AuctionRepository{
		db: db,
	}
}

func (r *AuctionRepository) CreateAuction(ctx context.Context, auction *domain.Auction, sellerID uuid.UUID, imageURLs []string) (*domain.Auction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO auctions (
			seller_id, title, description, starting_price, current_price, status, start_time, end_time
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, seller_id, title, description, starting_price, current_price, status, start_time, end_time, created_at
	`

	createdAuction := &domain.Auction{}

	err = tx.QueryRowContext(
		ctx,
		query,
		sellerID,
		auction.Title,
		auction.Description,
		auction.StartingPrice,
		auction.CurrentPrice,
		auction.Status,
		auction.StartTime,
		auction.EndTime,
	).Scan(
		&createdAuction.ID,
		&createdAuction.SellerID,
		&createdAuction.Title,
		&createdAuction.Description,
		&createdAuction.StartingPrice,
		&createdAuction.CurrentPrice,
		&createdAuction.Status,
		&createdAuction.StartTime,
		&createdAuction.EndTime,
		&createdAuction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(imageURLs) > 0 {
		imageQuery := `INSERT INTO auction_images (auction_id, image_url) VALUES ($1, $2)`
		for _, url := range imageURLs {
			_, err := tx.ExecContext(ctx, imageQuery, createdAuction.ID, url)
			if err != nil {
				return nil, err
			}
			createdAuction.Images = append(createdAuction.Images, url)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdAuction, nil
}

func (r *AuctionRepository) GetAuction(ctx context.Context, auctionID uuid.UUID) (*domain.Auction, error) {
	query := `
		SELECT id, seller_id, title, description, starting_price, current_price, status, start_time, end_time, created_at
		FROM auctions
		WHERE id = $1
	`

	auction := &domain.Auction{}
	err := r.db.QueryRowContext(ctx, query, auctionID).Scan(
		&auction.ID,
		&auction.SellerID,
		&auction.Title,
		&auction.Description,
		&auction.StartingPrice,
		&auction.CurrentPrice,
		&auction.Status,
		&auction.StartTime,
		&auction.EndTime,
		&auction.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	imageQuery := `SELECT image_url FROM auction_images WHERE auction_id = $1`
	rows, err := r.db.QueryContext(ctx, imageQuery, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var img string
		if err := rows.Scan(&img); err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	auction.Images = images

	return auction, nil
}
