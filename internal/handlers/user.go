package handlers

import (
	"net/http"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



type UserHandler struct {
	service domain.UserService
}



func NewUserHandler(service domain.UserService) *UserHandler  {
	return  &UserHandler{
		service: service,
	}
}



type CreateUserRequest struct {
	Email string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}


func (h *UserHandler) CreateUserHandler(ctx *gin.Context){
	var req CreateUserRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,utils.ErrorResponse("invalid payload",err))
		return
	}

	user := &domain.User{
		Email: req.Email,
		Password: req.Password,
	}

	createdUser, err := h.service.CreateUser(ctx.Request.Context(),user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,utils.ErrorResponse("sign up failed",err))
		return
	}


	if err := h.createUserSession(ctx,createdUser.ID.String()); err != nil {
		ctx.JSON(http.StatusInternalServerError,utils.ErrorResponse("failed to create session",err))
		return
	}

	userResponse := domain.UserResponse{
		ID: createdUser.ID,
		Email: createdUser.Email,
		IsVerified: createdUser.IsVerified,
	}


	ctx.JSON(http.StatusCreated,utils.SuccessResponse("sign up successful",userResponse))
}


func (h *UserHandler) createUserSession(ctx *gin.Context, userID string) error{
	session := sessions.Default(ctx)
	session.Set("user_id",userID)

	return session.Save()
}


type LoginRequest struct {
	Email string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}


func (h *UserHandler) LoginUser(ctx *gin.Context)  {
	var req LoginRequest
	
	if err := ctx.ShouldBind(&req);err != nil{
		ctx.JSON(http.StatusBadRequest,utils.ErrorResponse("invalid payload",err))
		return
	}

	loginData := &domain.UserLogin{
		Email: req.Email,
		Password: req.Password,
	}

	user, err := h.service.LoginUser(ctx.Request.Context(),loginData)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,utils.ErrorResponse("login failed",err))
		return
	}


	if err := h.createUserSession(ctx, user.ID.String()); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to create session", err))
		return
	}


	userResponse := domain.UserResponse{
		ID: user.ID,
		Email: user.Email,
		IsVerified: user.IsVerified,
	}

	ctx.JSON(http.StatusOK,utils.SuccessResponse("login success",userResponse))
}



func (h *UserHandler) GetUserProfile(ctx *gin.Context){
	userID := ctx.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("invalid session", err))
		return
	}


	user, err := h.service.GetUserProfile(ctx.Request.Context(),uid)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,utils.ErrorResponse("failed to fetch user profile",err))
		return
	}

	if user == nil {
		ctx.JSON(http.StatusBadRequest,utils.ErrorResponse("failed to fetch user profile",err))
		return
	}

	userResponse := domain.UserResponse{
		ID: user.ID,
		Email: user.Email,
		IsVerified: user.IsVerified,
	}


	ctx.JSON(http.StatusOK,utils.SuccessResponse("successfuly fetched user profile",userResponse))
}