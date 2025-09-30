package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aglili/auction-app/internal/domain"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id,email,password,is_verified,created_at FROM users WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.IsVerified, &user.Created_At)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// user does not exist
			return nil, ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id,email,password,is_verified,created_at FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.IsVerified, &user.Created_At)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `INSERT into users (email,password)
	VALUES ($1,$2)
	RETURNING id,email,password,is_verified,created_at
	`

	row := r.db.QueryRowContext(ctx, query, user.Email, user.Password)

	createdUser := &domain.User{}

	err := row.Scan(&createdUser.ID, &createdUser.Email, &createdUser.Password, &createdUser.IsVerified, &createdUser.Created_At)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `UPDATE users
	SET email = $1,password = $2, is_verified = $3
	WHERE id = $4
	RETURNING id, email, password, is_verified, created_at
	`

	row := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.IsVerified, user.ID)

	updatedUser := &domain.User{}

	err := row.Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.Password, &updatedUser.IsVerified, &updatedUser.Created_At)
	if err != nil {
		return &domain.User{}, err
	}

	return updatedUser, nil
}
