package main

import (
	"go.uber.org/zap"

	"auth/internal/app"
	grpc "auth/internal/infra/grpc/server"
	http "auth/internal/infra/http/server"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		zap.S().Fatalf("Failed to create App:\n%f", err)
	}

	grpcServer := grpc.NewServer(application)
	go func() {
		err = grpcServer.Run("0.0.0.0:50051")
		if err != nil {
			zap.S().Fatalf("Failed to run gRPC server:\n%f", err)
		}
	}()

	httpServer := http.NewServer(application)
	err = httpServer.Run("0.0.0.0:8000")
	if err != nil {
		zap.S().Fatalf("Failed to run HTTP server:\n%f", err)
	}
}
