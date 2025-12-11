package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"device/internal/app"
	m "device/metrics"
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
	instances := 600

	for i := 0; i < instances; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			run(id)
		}(i + 1)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("Metrics Per Second", m.Counter.Load())
				fmt.Println("Auths Total", m.AuthCounter.Load())
				m.Counter.Store(0)
			}
		}
	}()

	wg.Wait()
}
