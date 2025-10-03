package handlers

import (
	"errors"
	"net/http"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserHandler struct {
	service   domain.UserService
	validator *validator.Validate
}

func NewUserHandler(service domain.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		service:   service,
		validator: validator,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) CreateUserHandler(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.BindJSON(&req); err != nil {
		utils.RespondWithValidationError(ctx, err)
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := h.service.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		utils.RespondWithError(ctx, err, "failed to create user")
		return
	}

	if err := h.createUserSession(ctx, createdUser.ID.String()); err != nil {
		utils.RespondWithError(ctx, err, "failed to create session")
		return
	}

	userResponse := domain.UserResponse{
		ID:         createdUser.ID,
		Email:      createdUser.Email,
		IsVerified: createdUser.IsVerified,
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse("sign up successful", userResponse))
}

func (h *UserHandler) createUserSession(ctx *gin.Context, userID string) error {
	session := sessions.Default(ctx)
	session.Set("user_id", userID)

	return session.Save()
}

type LoginRequest struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBind(&req); err != nil {
		utils.RespondWithValidationError(ctx, err)
		return
	}

	loginData := &domain.UserLogin{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.service.LoginUser(ctx.Request.Context(), loginData)
	if err != nil {
		utils.RespondWithError(ctx, err, "Login failed")
		return
	}

	if err := h.createUserSession(ctx, user.ID.String()); err != nil {
		utils.RespondWithError(ctx, err, "Failed to create session")
		return
	}

	userResponse := domain.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("login success", userResponse))
}

func (h *UserHandler) GetUserProfile(ctx *gin.Context) {
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

	user, err := h.service.GetUserProfile(ctx.Request.Context(), uid)
	if err != nil {
		utils.RespondWithError(ctx, err, "Failed to fetch user profile")
		return
	}

	userResponse := domain.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("Successfully fetched user profile", userResponse))
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()

	if err := session.Save(); err != nil {
		utils.RespondWithError(ctx, err, "failed to logout")
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("logout success", nil))
}
