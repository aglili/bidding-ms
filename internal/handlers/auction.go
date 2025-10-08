package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AuctionHandler struct {
	service   domain.AuctionService
	validator *validator.Validate
}

func NewAuctionHandler(service domain.AuctionService, validator *validator.Validate) *AuctionHandler {
	return &AuctionHandler{
		service:   service,
		validator: validator,
	}
}

type CreateAuctionRequest struct {
	Title         string   `json:"title" binding:"required"`
	Description   *string  `json:"description,omitempty"`
	StartingPrice float64  `json:"starting_price" binding:"required,gt=0"`
	StartTime     string   `json:"start_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime       string   `json:"end_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Images        []string `json:"images" binding:"required"`
}

func (h *AuctionHandler) CreateAuctionHandler(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.RespondWithError(ctx, err, "invalid session")
		return
	}

	var req CreateAuctionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondWithValidationError(ctx, err)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		utils.RespondWithError(ctx, err, "invalid start_time format")
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		utils.RespondWithError(ctx, err, "invalid end_time format")
		return
	}

	if !endTime.After(startTime) {
		utils.RespondWithError(ctx, nil, "end_time must be after start_time")
		return
	}
	if startTime.Before(time.Now()) {
		utils.RespondWithError(ctx, nil, "start_time cannot be in the past")
		return
	}

	auction := &domain.Auction{
		Title:         req.Title,
		Description:   req.Description,
		StartingPrice: req.StartingPrice,
		CurrentPrice:  req.StartingPrice, // initial = starting price
		Status:        "open",
		StartTime:     startTime,
		EndTime:       endTime,
	}

	createdAuction, err := h.service.CreateAuction(ctx.Request.Context(), auction, uid, req.Images)
	if err != nil {
		utils.RespondWithError(ctx, err, "failed to create auction")
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("auction created successfully", createdAuction))
}

func (h *AuctionHandler) GetAuction(ctx *gin.Context) {
	auctionID := ctx.Param("id")
	if auctionID == "" {
		utils.RespondWithError(ctx, nil, "auctionID is required")
		return
	}

	auctionUUID, err := uuid.Parse(auctionID)
	if err != nil {
		utils.RespondWithError(ctx, err, "invalid auctionID format")
		return
	}

	auction, err := h.service.GetAuction(ctx.Request.Context(), auctionUUID)
	if err != nil {
		utils.RespondWithError(ctx, err, "error fetching auction")
		return
	}

	auctionResponse := domain.AuctionResponse{
		ID:            auction.ID,
		Title:         auction.Title,
		Description:   auction.Description,
		StartingPrice: auction.StartingPrice,
		CurrentPrice:  auction.CurrentPrice,
		Status:        auction.Status,
		StartTime:     auction.StartTime,
		EndTime:       auction.EndTime,
		Images:        auction.Images,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("successfully fetched auction", auctionResponse))
}

func (h *AuctionHandler) GetUserAuctions(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.RespondWithError(ctx, errors.New("user not authenticated"), "unauthorized")
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.RespondWithError(ctx, err, "invalid user ID")
		return
	}

	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	auctions, total, err := h.service.GetUserAuctions(ctx.Request.Context(), uid, page, limit)
	if err != nil {
		utils.RespondWithError(ctx, err, "failed to fetch auctions")
		return
	}

	auctionResponse := make([]domain.Auction, 0, len(auctions))
	for _, auction := range auctions {
		auctionResponse = append(auctionResponse, domain.Auction{
			ID:            auction.ID,
			Title:         auction.Title,
			Description:   auction.Description,
			Status:        auction.Status,
			StartingPrice: auction.StartingPrice,
			CurrentPrice:  auction.CurrentPrice,
			StartTime:     auction.StartTime,
			EndTime:       auction.EndTime,
			CreatedAt:     auction.CreatedAt,
			Images:        auction.Images,
		})
	}
	response := utils.PaginatedResponse("successfuly fetched auctions", auctionResponse, page, limit, total)

	ctx.JSON(http.StatusOK, response)
}

func (h *AuctionHandler) GetOpenAuctions(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.RespondWithError(ctx, errors.New("user not authenticated"), "unauthorized")
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.RespondWithError(ctx, err, "invalid user ID")
		return
	}

	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	auctions, total, err := h.service.GetOpenAuctions(ctx.Request.Context(), uid, page, limit)
	if err != nil {
		utils.RespondWithError(ctx, err, "failed to fetch open auctions")
		return
	}

	auctionResponse := make([]domain.Auction, 0, len(auctions))
	for _, auction := range auctions {
		auctionResponse = append(auctionResponse, domain.Auction{
			ID:            auction.ID,
			Title:         auction.Title,
			Description:   auction.Description,
			Status:        auction.Status,
			StartingPrice: auction.StartingPrice,
			CurrentPrice:  auction.CurrentPrice,
			StartTime:     auction.StartTime,
			EndTime:       auction.EndTime,
			CreatedAt:     auction.CreatedAt,
			Images:        auction.Images,
		})
	}

	response := utils.PaginatedResponse("successfuly fetched auctions", auctionResponse, page, limit, total)

	ctx.JSON(http.StatusOK, response)

}
