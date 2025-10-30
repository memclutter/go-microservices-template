package grpc

import (
	"context"
	"time"

	"github.com/memclutter/go-microservices-template/api/gen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	userUseCase "github.com/memclutter/go-microservices-template/internal/usecase/user"
	"github.com/memclutter/go-microservices-template/pkg/logger"
	"github.com/memclutter/go-microservices-template/pkg/metrics"
)

// UserServiceServer implements the gRPC UserService
type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	createUserUC *userUseCase.CreateUserUseCase
	getUserUC    *userUseCase.GetUserUseCase
	logger       *logger.Logger
	metrics      *metrics.Metrics
}

// NewUserServiceServer creates a new gRPC user service server
func NewUserServiceServer(
	createUserUC *userUseCase.CreateUserUseCase,
	getUserUC *userUseCase.GetUserUseCase,
	log *logger.Logger,
	metrics *metrics.Metrics,
) *UserServiceServer {
	return &UserServiceServer{
		createUserUC: createUserUC,
		getUserUC:    getUserUC,
		logger:       log,
		metrics:      metrics,
	}
}

// CreateUser creates a new user
func (s *UserServiceServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		s.metrics.GRPCRequestDuration.WithLabelValues("CreateUser").Observe(duration)
	}()

	s.logger.WithFields(map[string]any{
		"email": req.Email,
		"name":  req.Name,
	}).Info("CreateUser gRPC request")

	// Validate input
	if req.Email == "" {
		s.metrics.GRPCRequestsTotal.WithLabelValues("CreateUser", "invalid_argument").Inc()
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Name == "" {
		s.metrics.GRPCRequestsTotal.WithLabelValues("CreateUser", "invalid_argument").Inc()
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Password == "" {
		s.metrics.GRPCRequestsTotal.WithLabelValues("CreateUser", "invalid_argument").Inc()
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Execute use case
	input := userUseCase.CreateUserInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	output, err := s.createUserUC.Execute(ctx, input)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create user")
		s.metrics.GRPCRequestsTotal.WithLabelValues("CreateUser", "internal_error").Inc()
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	s.metrics.GRPCRequestsTotal.WithLabelValues("CreateUser", "ok").Inc()

	// Build response
	return &user.CreateUserResponse{
		User: &user.User{
			Id:    output.UserID,
			Email: output.Email,
			Name:  output.Name,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		s.metrics.GRPCRequestDuration.WithLabelValues("GetUser").Observe(duration)
	}()

	s.logger.WithField("user_id", req.UserId).Info("GetUser gRPC request")

	// Validate input
	if req.UserId == "" {
		s.metrics.GRPCRequestsTotal.WithLabelValues("GetUser", "invalid_argument").Inc()
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Execute use case
	input := userUseCase.GetUserInput{
		UserID: req.UserId,
	}

	output, err := s.getUserUC.Execute(ctx, input)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user")
		s.metrics.GRPCRequestsTotal.WithLabelValues("GetUser", "not_found").Inc()
		return nil, status.Error(codes.NotFound, "user not found")
	}

	s.metrics.GRPCRequestsTotal.WithLabelValues("GetUser", "ok").Inc()

	// Build response
	return &user.GetUserResponse{
		User: &user.User{
			Id:    output.ID,
			Email: output.Email,
			Name:  output.Name,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}, nil
}

// UpdateUser updates an existing user
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	// TODO: Implement UpdateUser use case
	s.metrics.GRPCRequestsTotal.WithLabelValues("UpdateUser", "unimplemented").Inc()
	return nil, status.Error(codes.Unimplemented, "UpdateUser not implemented yet")
}

// DeleteUser deletes a user
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	// TODO: Implement DeleteUser use case
	s.metrics.GRPCRequestsTotal.WithLabelValues("DeleteUser", "unimplemented").Inc()
	return nil, status.Error(codes.Unimplemented, "DeleteUser not implemented yet")
}

// ListUsers retrieves a list of users
func (s *UserServiceServer) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	// TODO: Implement ListUsers use case
	s.metrics.GRPCRequestsTotal.WithLabelValues("ListUsers", "unimplemented").Inc()
	return nil, status.Error(codes.Unimplemented, "ListUsers not implemented yet")
}
