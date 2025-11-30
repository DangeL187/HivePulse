package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"consumer/internal/app"
	"consumer/internal/infra/metrics"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer func() {
		_ = logger.Sync()
	}()

	metrics.RegisterAll()

	var err error
	application, err := app.NewApp()
	if err != nil {
		zap.S().Fatalf("Failed to create application:\n%f", err)
	}

	// === RUN ===

	application.Run()

	metricsServer := metrics.NewServer("0.0.0.0:2112")
	go func() {
		err = metricsServer.Run()
		if err != nil {
			zap.S().Fatalf("failed to run metrics server:\n%f", err)
		}
	}()

	// === STOP ===

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutdown initiated...")

	application.Stop()

	err = metricsServer.Stop()
	if err != nil {
		zap.S().Fatalf("failed to stop metrics server:\n%f", err)
	}

	zap.L().Info("Shutdown completed")
}
