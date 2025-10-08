package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/repository"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/google/uuid"
)

type AuctionService struct {
	repository domain.AuctionRepository
}

func NewAuctionService(repository domain.AuctionRepository) *AuctionService {
	return &AuctionService{
		repository: repository,
	}
}

func (s *AuctionService) CreateAuction(ctx context.Context, auction *domain.Auction, sellerID uuid.UUID, imageURLs []string) (*domain.Auction, error) {

	auction, err := s.repository.CreateAuction(ctx, auction, sellerID, imageURLs)
	if err != nil {
		return nil, err
	}

	return auction, nil
}

func (s *AuctionService) GetAuction(ctx context.Context, auctionID uuid.UUID) (*domain.Auction, error) {
	auction, err := s.repository.GetAuction(ctx, auctionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, utils.NewAppError(err, "auction not found", utils.ErrCodeNotFound, http.StatusNotFound)
		}

		return nil, utils.NewAppError(err, "failed to fetch auction", utils.ErrCodeInternal, http.StatusInternalServerError)
	}

	return auction, nil
}

func (s *AuctionService) GetUserAuctions(ctx context.Context, userID uuid.UUID, page, limit int) ([]*domain.Auction, int, error) {
	return s.repository.GetUserAuctions(ctx, userID, page, limit)
}

func (s *AuctionService) GetOpenAuctions(ctx context.Context, userID uuid.UUID, page, limit int) ([]*domain.Auction, int, error) {
	return s.repository.GetOpenAuctions(ctx, userID, page, limit)
}
