package module

import (
	"auth/internal/features/device/domain"
	"auth/internal/features/device/usecase"
	"auth/internal/shared/token"
)

type Module struct {
	Auth   *usecase.AuthUseCase
	Device *usecase.DeviceUseCase
}

func NewModule(repo domain.Repository, tokenGenerator token.Manager) *Module {
	return &Module{
		Auth:   usecase.NewAuthUseCase(repo, tokenGenerator),
		Device: usecase.NewDeviceUseCase(repo),
	}
}
