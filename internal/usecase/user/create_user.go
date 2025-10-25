package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/memclutter/go-microservices-template/internal/domain/user"
	"github.com/memclutter/go-microservices-template/pkg/logger"
)

// CreateUserInput represents input data for creating a user
type CreateUserInput struct {
	Email    string
	Name     string
	Password string
}

// CreateUserOutput represents the result of user creation
type CreateUserOutput struct {
	UserID string
	Email  string
	Name   string
}

// CreateUserUseCase handles user creation business flow
type CreateUserUseCase struct {
	repo          user.Repository
	domainService user.Service
	eventPub      EventPublisher
	logger        *logger.Logger
}

// NewCreateUserUseCase creates a new use case instance
func NewCreateUserUseCase(
	repo user.Repository,
	domainService user.Service,
	eventPub EventPublisher,
	logger *logger.Logger,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo:          repo,
		domainService: domainService,
		eventPub:      eventPub,
		logger:        logger,
	}
}

// Execute creates a new user
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// Log use case execution
	uc.logger.WithFields(map[string]any{
		"email": input.Email,
		"name":  input.Name,
	}).Info("Creating new user")

	// 1. Check if email is unique (domain service)
	isUnique, err := uc.domainService.IsEmailUnique(ctx, input.Email)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to check email uniqueness")
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if !isUnique {
		return nil, user.ErrUserAlreadyExists
	}

	// 2. Create domain entity (with validation)
	newUser, err := user.NewUser(input.Email, input.Name, input.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	// 3. Generate ID
	newUser.ID = uuid.New().String()

	// 4. Save to repository
	if err := uc.repo.Create(ctx, newUser); err != nil {
		uc.logger.WithError(err).Error("Failed to create user in database")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 5. Publish domain event
	event := user.UserCreatedEvent{
		UserID:    newUser.ID,
		Email:     newUser.Email,
		Name:      newUser.Name,
		CreatedAt: newUser.CreatedAt,
	}
	if err := uc.eventPub.Publish(ctx, user.EventTypeUserCreated, event); err != nil {
		// Don't fail the use case, just log the error
		uc.logger.WithError(err).Warn("Failed to publish user created event")
	}

	uc.logger.WithField("user_id", newUser.ID).Info("User created successfully")

	return &CreateUserOutput{
		UserID: newUser.ID,
		Email:  newUser.Email,
		Name:   newUser.Name,
	}, nil
}
