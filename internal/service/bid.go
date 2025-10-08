package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/repository"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type BidService struct {
	bidRepo     domain.BidRepository
	auctionRepo domain.AuctionRepository
	cache       *redis.Client
}

func NewBidService(bidRepo domain.BidRepository, auctionRepo domain.AuctionRepository, cache *redis.Client) *BidService {
	return &BidService{
		bidRepo:     bidRepo,
		auctionRepo: auctionRepo,
		cache:       cache,
	}
}

func (s *BidService) CreateBid(ctx context.Context, auctionID, userID uuid.UUID, amount float64) error {
	auction, err := s.auctionRepo.GetAuction(ctx, auctionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return utils.NewAppError(err, "auction not found", utils.ErrCodeNotFound, http.StatusNotFound)
		}
		return utils.NewAppError(err, "failed to fetch auction", utils.ErrCodeInternal, http.StatusInternalServerError)
	}

	now := time.Now()
	if now.Before(auction.StartTime) {
		return utils.NewAppError(nil, "auction hasn't started yet", utils.ErrCodeForbidden, http.StatusForbidden)
	}
	if now.After(auction.EndTime) {
		return utils.NewAppError(nil, "auction has ended", utils.ErrCodeForbidden, http.StatusForbidden)
	}

	if auction.SellerID == userID {
		return utils.NewAppError(nil, "cannot bid on own auction", utils.ErrCodeForbidden, http.StatusForbidden)
	}

	if auction.Status == "closed" {
		return utils.NewAppError(nil, "auction has ended", utils.ErrCodeForbidden, http.StatusForbidden)
	}

	key := fmt.Sprintf("auction:%s:highest_bid", auctionID.String())
	bidderKey := fmt.Sprintf("auction:%s:highest_bidder", auctionID.String())

	// Use Redis WATCH for optimistic locking
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err = s.cache.Watch(ctx, func(tx *redis.Tx) error {
			highestBidStr, err := tx.Get(ctx, key).Result()
			var highestBid float64

			if err == redis.Nil {
				highestBid = auction.StartingPrice
			} else if err != nil {
				return fmt.Errorf("failed to get highest bid: %w", err)
			} else {
				highestBid, err = strconv.ParseFloat(highestBidStr, 64)
				if err != nil {
					return fmt.Errorf("invalid highest bid format: %w", err)
				}
			}

			// Simple validation: bid must be higher than current highest
			if amount <= highestBid {
				return utils.NewAppError(nil,
					fmt.Sprintf("bid must be higher than current highest bid of %.2f", highestBid),
					utils.ErrCodeNotAllowed, http.StatusBadRequest)
			}

			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, key, amount, 0)
				pipe.Set(ctx, bidderKey, userID.String(), 0)
				return nil
			})
			return err
		}, key)

		if err == nil {
			break
		}
		if err == redis.TxFailedErr {
			continue
		}
		return utils.NewAppError(err, "failed to update cache", utils.ErrCodeInternal, http.StatusInternalServerError)
	}

	if err != nil {
		return utils.NewAppError(err, "failed to place bid after retries", utils.ErrCodeInternal, http.StatusInternalServerError)
	}

	if err := s.bidRepo.CreateBid(ctx, auctionID, userID, amount); err != nil {
		return utils.NewAppError(err, "failed to save bid", utils.ErrCodeInternal, http.StatusInternalServerError)
	}

	if err := s.auctionRepo.UpdateCurrentPrice(ctx, auctionID, amount); err != nil {
		return utils.NewAppError(err, "failed to update auction price", utils.ErrCodeInternal, http.StatusInternalServerError)
	}

	return nil
}
