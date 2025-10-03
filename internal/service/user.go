package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/repository"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/google/uuid"
)

type UserService struct {
	repository domain.UserRepository
}

func NewUserService(repository domain.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := s.repository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, utils.NewAppError(err, "failed to create user", utils.ErrCodeDatabaseError, http.StatusInternalServerError)
		}
	}

	if existingUser != nil && existingUser.ID != uuid.Nil {
		return nil, utils.NewAppError(fmt.Errorf("user with email %s already exists", user.Email), "User already exists", utils.ErrCodeConflict, http.StatusConflict)
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, utils.NewAppError(err, "failed to create user", utils.ErrCodeInternal, http.StatusInternalServerError)
	}

	user.Password = hashedPassword

	created_user, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, utils.NewAppError(err, "failed to create user", utils.ErrCodeDatabaseError, http.StatusInternalServerError)
	}

	created_user.Password = ""

	return created_user, nil
}

func (s *UserService) LoginUser(ctx context.Context, payload *domain.UserLogin) (*domain.User, error) {
	existing_user, err := s.repository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, utils.NewAppError(err, "invalid credentials", utils.ErrCodeUnauthorized, http.StatusUnauthorized)
		}

		return nil, utils.NewAppError(err, "login failed", utils.ErrCodeDatabaseError, http.StatusInternalServerError)
	}

	validPassword := utils.ComparePassword(payload.Password, existing_user.Password)
	if !validPassword {
		return nil, utils.NewAppError(fmt.Errorf("invalid credentials"), "invalid credentials", utils.ErrCodeUnauthorized, http.StatusUnauthorized)
	}

	existing_user.Password = ""

	return existing_user, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, utils.NewAppError(err, "user not found", utils.ErrCodeNotFound, http.StatusNotFound)
		}
		return nil, utils.NewAppError(err, "failed to retrieve user", utils.ErrCodeDatabaseError, http.StatusInternalServerError)
	}

	user.Password = ""

	return user, nil
}
