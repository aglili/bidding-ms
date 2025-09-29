package service

import (
	"context"
	"fmt"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/google/uuid"
)




type UserService struct {
	repository domain.UserRepository
}



func NewUserService(repository domain.UserRepository) *UserService{
	return &UserService{
		repository: repository,
	}
}



func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User,error) {
	existingUser, err := s.repository.GetUserByEmail(ctx,user.Email)
	if err != nil{
		return nil,err
	}

	if existingUser != nil && existingUser.ID != uuid.Nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword,err := utils.HashPassword(user.Password)
	if err != nil {
		return nil,err
	}

	user.Password = hashedPassword


	created_user,err := s.repository.CreateUser(ctx,user)
	if err != nil {
		return  nil, err
	}

	return  created_user,nil
}



func (s *UserService) LoginUser(ctx context.Context, payload *domain.UserLogin) (*domain.User, error){
	existing_user,err := s.repository.GetUserByEmail(ctx,payload.Email)
	if err != nil{
		return nil,err
	}

	if existing_user == nil{
		return nil,fmt.Errorf("invalid email or password")
	}


	validPassword := utils.ComparePassword(payload.Password,existing_user.Password)
	if !validPassword{
		return nil, fmt.Errorf("invalid email or password")
	}

	return  existing_user,nil
}



func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.User,error)  {
	return s.repository.GetUserByID(ctx,userID)
}


