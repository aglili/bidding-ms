package handlers

import (
	"net/http"
	"time"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuctionHandler struct {
	service domain.AuctionService
}

func NewAuctionHandler(service domain.AuctionService) *AuctionHandler {
	return &AuctionHandler{
		service: service,
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
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("invalid session", err))
		return
	}

	var req CreateAuctionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid payload", err))
		return
	}

	// TODO: Move the validation to the services and use playground validate for the validating requests

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid start_time format", err))
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid end_time format", err))
		return
	}

	if !endTime.After(startTime) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("end_time must be after start_time", nil))
		return
	}

	if startTime.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("start_time cannot be in the past", nil))
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
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to create auction", err))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("auction created successfully", createdAuction))
}

func (h *AuctionHandler) GetAuction(ctx *gin.Context) {
	auctionID := ctx.Param("id")

	if auctionID == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("auctionID cannot is required", nil))
		return
	}

	auctionUUID, err := uuid.Parse(auctionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("auctionID cannot is invalid", err))
		return
	}

	auction, err := h.service.GetAuction(ctx.Request.Context(), auctionUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("error fetching auction", err))
		return
	}

	if auction == nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("auction not found", nil))
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

	ctx.JSON(http.StatusOK, utils.SuccessResponse("successfuly fetched auction", auctionResponse))

}
