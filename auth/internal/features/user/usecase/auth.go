package usecase

import (
	"time"

	"github.com/DangeL187/erax"

	"auth/internal/features/user/domain"
	"auth/internal/shared/auth"
	"auth/internal/shared/token"
)

type authRepo interface {
	CreateUser(userInputData domain.UserInputData) (uint, error)
	GetUserByEmail(email string) (domain.User, error)
}

type AuthUseCase struct {
	repo           authRepo
	tokenGenerator token.Generator
}

type LoginInput struct {
	Email    string
	Password string
}

func (l LoginInput) Validate() error {
	if l.Email == "" || l.Password == "" {
		return auth.ErrInvalidCredentials
	}

	return nil
}

func (a *AuthUseCase) Login(accessTokenTTL time.Duration, request LoginInput) (string, error) {
	if err := request.Validate(); err != nil {
		return "", erax.WrapWithError(err, auth.ErrInvalidCredentials, "failed to validate input")
	}

	user, err := a.repo.GetUserByEmail(request.Email)
	if err != nil {
		return "", erax.WrapWithError(err, auth.ErrInvalidCredentials, "failed to get user by email")
	}

	isValid := auth.VerifyPassword(user.PasswordHash, request.Password)
	if !isValid {
		return "", erax.Wrap(auth.ErrInvalidCredentials, "failed to verify password")
	}

	accessToken, err := a.tokenGenerator.Generate(user.ID, "access", accessTokenTTL)
	if err != nil {
		return "", erax.Wrap(err, "failed to generate access token")
	}

	return accessToken, nil
}

type RegisterInput struct {
	Email    string
	FullName string
	Password string
}

func (r RegisterInput) Validate() error {
	if r.Email == "" || r.FullName == "" || r.Password == "" {
		return auth.ErrInvalidCredentials
	}

	return nil
}

func (a *AuthUseCase) Register(accessTokenTTL time.Duration, request RegisterInput) (string, error) {
	if err := request.Validate(); err != nil {
		return "", erax.WrapWithError(err, auth.ErrInvalidCredentials, "failed to validate input")
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		return "", erax.Wrap(err, "failed to hash password")
	}

	userInputData := domain.UserInputData{
		Email:        request.Email,
		FullName:     request.FullName,
		PasswordHash: hashedPassword,
	}

	userID, err := a.repo.CreateUser(userInputData)
	if err != nil {
		return "", erax.WrapWithError(err, auth.ErrInvalidCredentials, "failed to create user")
	}

	accessToken, err := a.tokenGenerator.Generate(userID, "access", accessTokenTTL)
	if err != nil {
		return "", erax.Wrap(err, "failed to generate access token")
	}

	return accessToken, nil
}

func NewAuthUseCase(repo authRepo, tokenGenerator token.Generator) *AuthUseCase {
	return &AuthUseCase{
		repo:           repo,
		tokenGenerator: tokenGenerator,
	}
}
