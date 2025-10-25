package user

import (
	"context"
	"testing"

	"github.com/memclutter/go-microservices-template/internal/domain/user"
	"github.com/memclutter/go-microservices-template/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

type MockDomainService struct {
	mock.Mock
}

func (m *MockDomainService) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockDomainService) CanUserBeDeleted(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(ctx context.Context, eventType string, payload interface{}) error {
	args := m.Called(ctx, eventType, payload)
	return args.Error(0)
}

func TestCreateUserUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateUserInput
		setup   func(*MockRepository, *MockDomainService, *MockEventPublisher)
		wantErr error
	}{
		{
			name: "successful user creation",
			input: CreateUserInput{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			setup: func(repo *MockRepository, ds *MockDomainService, pub *MockEventPublisher) {
				ds.On("IsEmailUnique", mock.Anything, "test@example.com").Return(true, nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				pub.On("Publish", mock.Anything, user.EventTypeUserCreated, mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "email already exists",
			input: CreateUserInput{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			setup: func(repo *MockRepository, ds *MockDomainService, pub *MockEventPublisher) {
				ds.On("IsEmailUnique", mock.Anything, "test@example.com").Return(false, nil)
			},
			wantErr: user.ErrUserAlreadyExists,
		},
		{
			name: "invalid email",
			input: CreateUserInput{
				Email:    "",
				Name:     "Test User",
				Password: "password123",
			},
			setup:   func(repo *MockRepository, ds *MockDomainService, pub *MockEventPublisher) {},
			wantErr: user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			repo := new(MockRepository)
			domainService := new(MockDomainService)
			eventPub := new(MockEventPublisher)
			log := logger.New("test")

			tt.setup(repo, domainService, eventPub)

			// Create use case
			uc := NewCreateUserUseCase(repo, domainService, eventPub, log)

			// Execute
			result, err := uc.Execute(context.Background(), tt.input)

			// Assert
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.UserID)
				assert.Equal(t, tt.input.Email, result.Email)
			}

			// Verify mocks
			repo.AssertExpectations(t)
			domainService.AssertExpectations(t)
			eventPub.AssertExpectations(t)
		})
	}
}
