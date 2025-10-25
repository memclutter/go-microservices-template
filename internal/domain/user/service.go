package user

import (
	"context"
	"errors"
)

// Service defines domain operations that don't naturally fit in entities
// These are domain services, not application services
type Service interface {
	// IsEmailUnique checks if email is not already taken
	IsEmailUnique(ctx context.Context, email string) (bool, error)

	// CanUserBeDeleted checks business rules for user deletion
	CanUserBeDeleted(ctx context.Context, userID string) (bool, error)
}

type service struct {
	repo Repository
}

// NewService creates a new domain service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	_, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == ErrUserNotFound {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (s *service) CanUserBeDeleted(ctx context.Context, userID string) (bool, error) {
	// Business rule: Check if user has any dependencies
	// For example, active subscriptions, pending orders, etc.
	// This is just an example, implement your actual business rules
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Example business rule
	if user.Email == "admin@example.com" {
		return false, errors.New("admin user cannot be deleted")
	}

	return true, nil
}
