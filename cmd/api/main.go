package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	user2 "github.com/memclutter/go-microservices-template/api/gen/user"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/memclutter/go-microservices-template/internal/domain/user"
	"github.com/memclutter/go-microservices-template/internal/infrastructure/database"
	grpcHandler "github.com/memclutter/go-microservices-template/internal/infrastructure/grpc"
	"github.com/memclutter/go-microservices-template/internal/infrastructure/messaging/rabbitmq"
	"github.com/memclutter/go-microservices-template/internal/infrastructure/repository/postgres"
	userUseCase "github.com/memclutter/go-microservices-template/internal/usecase/user"
	"github.com/memclutter/go-microservices-template/pkg/config"
	"github.com/memclutter/go-microservices-template/pkg/logger"
	"github.com/memclutter/go-microservices-template/pkg/metrics"
)

func main() {
	// Load configuration
	cfg, err := config.Load(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.App.Env)
	log.Info("Starting microservices application",
		"env", cfg.App.Env,
		"app", cfg.App.Name,
	)

	// Initialize metrics
	appMetrics := metrics.NewMetrics("microservices")

	// Initialize database connection
	ctx := context.Background()
	dbPool, err := database.NewPostgresPool(ctx, &cfg.Database, log)
	if err != nil {
		log.WithError(err).Error("Failed to connect to database")
		os.Exit(1)
	}
	defer database.ClosePostgresPool(dbPool, log)

	// Initialize RabbitMQ publisher
	eventPublisher, err := rabbitmq.NewPublisher(cfg.RabbitMQ.GetRabbitMQURL(), log)
	if err != nil {
		log.WithError(err).Error("Failed to create RabbitMQ publisher")
		os.Exit(1)
	}
	defer func() {
		if err := eventPublisher.Close(); err != nil {
			log.WithError(err).Warn("Failed to close RabbitMQ publisher")
		}
	}()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(dbPool)

	// Initialize domain services
	userDomainService := user.NewService(userRepo)

	// Initialize use cases
	createUserUC := userUseCase.NewCreateUserUseCase(userRepo, userDomainService, eventPublisher, log)
	getUserUC := userUseCase.NewGetUserUseCase(userRepo, log)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	userGRPCService := grpcHandler.NewUserServiceServer(createUserUC, getUserUC, log, appMetrics)
	user2.RegisterUserServiceServer(grpcServer, userGRPCService)

	// Enable gRPC reflection for tools like grpcurl
	reflection.Register(grpcServer)

	// Start gRPC server in goroutine
	grpcAddr := fmt.Sprintf(":%d", cfg.GRPC.Port)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.WithError(err).Error("Failed to create gRPC listener")
		os.Exit(1)
	}

	go func() {
		log.Info("Starting gRPC server", "address", grpcAddr)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.WithError(err).Error("gRPC server failed")
		}
	}()

	// Initialize HTTP gateway
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register gRPC-gateway
	err = user2.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		log.WithError(err).Error("Failed to register gateway")
		os.Exit(1)
	}

	// Create HTTP mux
	mux := http.NewServeMux()

	// Register gateway routes
	mux.Handle("/v1/", gwmux)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Readiness check endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if err := dbPool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	})

	// Metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Start HTTP server
	httpAddr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	httpServer := &http.Server{
		Addr:         httpAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Starting HTTP server", "address", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down servers...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("HTTP server shutdown failed")
	}

	// Stop gRPC server
	grpcServer.GracefulStop()

	log.Info("Servers stopped gracefully")
}
