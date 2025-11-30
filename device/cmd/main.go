package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"device/internal/app"
)

func run(id int) {
	application := app.NewApp(id)
	application.Run()

	zap.L().Info("Device started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutdown initiated...")
	application.Stop()
	zap.L().Info("Shutdown completed")
}

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer func() {
		_ = logger.Sync()
	}()

	var wg sync.WaitGroup
	instances := 200

	for i := 0; i < instances; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			run(id)
		}(i + 1)
	}

	wg.Wait()
}
