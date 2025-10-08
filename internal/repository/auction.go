package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
			return nil, ErrNotFound
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

func (r *AuctionRepository) GetUserAuctions(ctx context.Context, userID uuid.UUID, page, limit int) ([]*domain.Auction, int, error) {
	offset := (page - 1) * limit

	query := `
		SELECT 
			a.id,
			a.seller_id,
			a.title,
			a.description,
			a.starting_price,
			a.current_price,
			a.status,
			a.start_time,
			a.end_time,
			a.created_at,
			COALESCE(ARRAY_AGG(ai.image_url) FILTER (WHERE ai.image_url IS NOT NULL), '{}') AS images
		FROM auctions a
		LEFT JOIN auction_images ai ON ai.auction_id = a.id
		WHERE a.seller_id = $1
		GROUP BY a.id
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var auctions []*domain.Auction
	for rows.Next() {
		auction := &domain.Auction{}
		err := rows.Scan(
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
			pq.Array(&auction.Images), // scan array of images
		)
		if err != nil {
			return nil, 0, err
		}
		auctions = append(auctions, auction)
	}

	// total count
	countQuery := `SELECT COUNT(*) FROM auctions WHERE seller_id = $1`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	return auctions, total, nil
}

func (r *AuctionRepository) UpdateCurrentPrice(ctx context.Context, auctionID uuid.UUID, amount float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE auctions SET current_price = $1 WHERE id = $2`,
		amount, auctionID,
	)
	return err
}

func (r *AuctionRepository) CloseAuction(ctx context.Context, auctionID uuid.UUID) error {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx,
		`UPDATE auctions SET status = 'closed' WHERE id = $1 RETURNING id`,
		auctionID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return fmt.Errorf("auction with id %s not found", auctionID)
	}

	return err
}

func (r *AuctionRepository) GetEndedActiveAuctions(ctx context.Context, currentTime time.Time) ([]*domain.Auction, error) {
	var auctions []*domain.Auction

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, seller_id, title, description, starting_price, current_price, status, start_time, end_time, created_at 
         FROM auctions 
         WHERE end_time <= $1 AND status = 'open'`,
		currentTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var auction domain.Auction
		err := rows.Scan(
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
			return nil, err
		}
		auctions = append(auctions, &auction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return auctions, nil
}

func (r *AuctionRepository) GetOpenAuctions(ctx context.Context, userID uuid.UUID, page, limit int) ([]*domain.Auction, int, error) {
	offset := (page - 1) * limit

	query := `
		SELECT 
			a.id,
			a.title,
			a.description,
			a.starting_price,
			a.current_price,
			a.status,
			a.start_time,
			a.end_time,
			a.created_at,
			COALESCE(ARRAY_AGG(ai.image_url) FILTER (WHERE ai.image_url IS NOT NULL), '{}') AS images
		FROM auctions a
		LEFT JOIN auction_images ai ON ai.auction_id = a.id
		WHERE a.status = 'open' AND a.seller_id != $1
		GROUP BY a.id
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var auctions []*domain.Auction
	for rows.Next() {
		auction := &domain.Auction{}
		err := rows.Scan(
			&auction.ID,
			&auction.Title,
			&auction.Description,
			&auction.StartingPrice,
			&auction.CurrentPrice,
			&auction.Status,
			&auction.StartTime,
			&auction.EndTime,
			&auction.CreatedAt,
			pq.Array(&auction.Images), // scan array of images
		)
		if err != nil {
			return nil, 0, err
		}
		auctions = append(auctions, auction)
	}

	countQuery := `SELECT COUNT(*) FROM auctions WHERE status = 'open' AND seller_id != $1`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	return auctions, total, nil
}
