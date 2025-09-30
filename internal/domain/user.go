package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"password" db:"-"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	Created_At time.Time `json:"created_at" db:"created_at"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	LoginUser(ctx context.Context, payload *UserLogin) (*User, error)
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*User, error)
}
