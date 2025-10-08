package handlers

import (
	"errors"
	"net/http"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BidHandler struct {
	bidService domain.BidService
	validator  *validator.Validate
}

func NewBidHandler(service domain.BidService, validator *validator.Validate) *BidHandler {
	return &BidHandler{
		bidService: service,
		validator:  validator,
	}
}

type CreateBidRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *BidHandler) CreateBid(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		response := utils.ErrorResponse("Unauthorized", errors.New("user not authenticated"))
		if response.Error != nil {
			response.Error.Code = utils.ErrCodeUnauthorized
		}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		response := utils.ErrorResponse("Invalid user ID", err)
		if response.Error != nil {
			response.Error.Code = utils.ErrCodeInvalidInput
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	auctionIDString := utils.GetParamStr(ctx, "id", "")

	auctionID, err := uuid.Parse(auctionIDString)
	if err != nil {
		response := utils.ErrorResponse("Invalid auction ID", err)
		if response.Error != nil {
			response.Error.Code = utils.ErrCodeInvalidInput
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var req CreateBidRequest

	if err := ctx.BindJSON(&req); err != nil {
		utils.RespondWithValidationError(ctx, err)
		return
	}

	err = h.bidService.CreateBid(ctx.Request.Context(), auctionID, uid, req.Amount)
	if err != nil {
		utils.RespondWithError(ctx, err, "failed to create bid")
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("bid created successfully", nil))
}
