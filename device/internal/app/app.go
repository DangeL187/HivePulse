package app

import (
	"context"
	"strconv"

	"device/internal/http"
	"device/internal/mqtt"
	"device/internal/shared/config"
	"device/internal/shared/tokens"
)

type App struct {
	config *config.Config
	tokens *tokens.Tokens

	authService    *http.AuthService
	metricsService *mqtt.MetricsService

	cancel context.CancelFunc
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	go a.metricsService.Run(ctx)
}

func (a *App) Stop() {
	a.cancel()

	a.metricsService.Stop()
}

func NewApp(id int) *App {
	cfg := config.NewConfig()
	cfg.DeviceID = "dev-" + strconv.Itoa(id)

	t := &tokens.Tokens{}
	authService := http.NewAuthService(cfg, t)

	return &App{
		config:         cfg,
		tokens:         t,
		metricsService: mqtt.NewMetricsService(cfg, t, authService),
	}
}
