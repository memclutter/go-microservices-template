package user

import (
	"context"
	"fmt"

	"github.com/memclutter/go-microservices-template/internal/domain/user"
	"github.com/memclutter/go-microservices-template/pkg/logger"
)

// GetUserInput represents input for getting a user
type GetUserInput struct {
	UserID string
}

// GetUserOutput represents user data
type GetUserOutput struct {
	ID    string
	Email string
	Name  string
}

// GetUserUseCase handles retrieving user data
type GetUserUseCase struct {
	repo   user.Repository
	logger *logger.Logger
}

// NewGetUserUseCase creates a new use case
func NewGetUserUseCase(repo user.Repository, logger *logger.Logger) *GetUserUseCase {
	return &GetUserUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute retrieves a user by ID
func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	uc.logger.WithField("user_id", input.UserID).Debug("Getting user")

	u, err := uc.repo.GetByID(ctx, input.UserID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, err
		}
		uc.logger.WithError(err).Error("Failed to get user from database")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &GetUserOutput{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}, nil
}
