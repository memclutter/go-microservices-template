package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/memclutter/go-microservices-template/internal/domain/user"
	"github.com/memclutter/go-microservices-template/internal/infrastructure/repository/sqlc"
)

// UserRepository implements user.Repository interface using PostgreSQL
type UserRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	params := sqlc.CreateUserParams{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	_, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user.User{
		ID:        row.ID,
		Email:     row.Email,
		Name:      row.Name,
		Password:  row.Password,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user.User{
		ID:        row.ID,
		Email:     row.Email,
		Name:      row.Name,
		Password:  row.Password,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	params := sqlc.UpdateUserParams{
		ID:        u.ID,
		Name:      u.Name,
		UpdatedAt: u.UpdatedAt,
	}

	_, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete removes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	rows, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*user.User, len(rows))
	for i, row := range rows {
		users[i] = &user.User{
			ID:        row.ID,
			Email:     row.Email,
			Name:      row.Name,
			Password:  row.Password,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	}

	return users, nil
}
