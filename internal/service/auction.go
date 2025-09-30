package service

import (
	"context"
	"fmt"

	"github.com/aglili/auction-app/internal/domain"
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
		return nil, err
	}

	// TODO: Handle this properly 
	if auction == nil {
		return nil, fmt.Errorf("auction not found")
	}

	return auction, nil
}
