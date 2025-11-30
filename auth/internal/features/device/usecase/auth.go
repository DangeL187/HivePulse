package usecase

import (
	"errors"
	"time"

	"github.com/DangeL187/erax"

	"auth/internal/features/device/domain"
	"auth/internal/shared/auth"
	"auth/internal/shared/token"
)

type authRepo interface {
	CreateDevice(deviceInputData domain.DeviceInputData) (uint, error)
	GetDeviceByDeviceID(deviceID string) (domain.Device, error)
}

type AuthUseCase struct {
	repo         authRepo
	tokenManager token.Manager
}

type AuthInput struct {
	DeviceID string
	Password string
}

func (a AuthInput) Validate() error {
	if a.DeviceID == "" || a.Password == "" {
		return auth.ErrInvalidCredentials
	}

	return nil
}

func (a *AuthUseCase) Login(accessTokenTTL, refreshTokenTTL time.Duration, request AuthInput) (string, string, error) {
	if err := request.Validate(); err != nil {
		return "", "", erax.WrapWithError(err, auth.ErrInvalidCredentials, "failed to validate input")
	}

	device, err := a.repo.GetDeviceByDeviceID(request.DeviceID)
	if err != nil {
		return "", "", erax.WrapWithError(err, auth.ErrInvalidCredentials, "failed to get device by device ID")
	}

	isValid := auth.VerifyPassword(device.PasswordHash, request.Password)
	if !isValid {
		return "", "", erax.Wrap(auth.ErrInvalidCredentials, "failed to verify password")
	}

	accessToken, err := a.tokenManager.Generate(device.ID, "access", accessTokenTTL)
	if err != nil {
		return "", "", erax.Wrap(err, "failed to generate token")
	}

	refreshToken, err := a.tokenManager.Generate(device.ID, "refresh", refreshTokenTTL)
	if err != nil {
		return "", "", erax.Wrap(err, "failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func (a *AuthUseCase) Register(accessTokenTTL, refreshTokenTTL time.Duration, request AuthInput) (string, string, error) {
	if err := request.Validate(); err != nil {
		return "", "", erax.WrapWithError(err, auth.ErrInvalidCredentials, "failed to validate input")
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		return "", "", erax.Wrap(err, "failed to hash password")
	}

	deviceInputData := domain.DeviceInputData{
		DeviceID:     request.DeviceID,
		PasswordHash: hashedPassword,
	}

	deviceID, err := a.repo.CreateDevice(deviceInputData)
	if err != nil {
		return "", "", erax.WrapWithError(err, auth.ErrInvalidCredentials, "failed to create device")
	}

	accessToken, err := a.tokenManager.Generate(deviceID, "access", accessTokenTTL)
	if err != nil {
		return "", "", erax.Wrap(err, "failed to generate token")
	}

	refreshToken, err := a.tokenManager.Generate(deviceID, "refresh", refreshTokenTTL)
	if err != nil {
		return "", "", erax.Wrap(err, "failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func (a *AuthUseCase) Refresh(accessTokenTTL time.Duration, refreshToken string) (accessToken string, err error) {
	deviceID, tokenType, err := a.tokenManager.ParseToken(refreshToken)
	if err != nil {
		return "", erax.Wrap(err, "failed to parse token")
	}
	if tokenType != "refresh" {
		return "", errors.New("token is not a refresh token")
	}

	accessToken, err = a.tokenManager.Generate(deviceID, "access", accessTokenTTL)
	if err != nil {
		return "", erax.Wrap(err, "failed to generate access token")
	}

	return accessToken, nil
}

func NewAuthUseCase(repo authRepo, tokenManager token.Manager) *AuthUseCase {
	return &AuthUseCase{
		repo:         repo,
		tokenManager: tokenManager,
	}
}
