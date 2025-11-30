package usecase

import (
	"errors"
	"fmt"
	"time"

	"auth/internal/features/user/domain"
)

type roleGranter interface {
	GrantRole(userID uint, role string) (bool, error)
}

type userGetter interface {
	GetUserIDByEmail(email string) (uint, error)
}

type userRegistrator interface {
	Register(accessTokenTTL time.Duration, request RegisterInput) (string, error)
}

type SeedUseCase struct {
	roleGranter     roleGranter
	userGetter      userGetter
	userRegistrator userRegistrator
}

func (s *SeedUseCase) SeedAdmin() error {
	adminUserID, err := s.userGetter.GetUserIDByEmail("admin")
	if errors.Is(err, domain.ErrUserNotFound) {
		request := RegisterInput{
			Email:    "admin",
			FullName: "admin",
			Password: "admin",
		}
		_, err = s.userRegistrator.Register(0, request)
		if err != nil {
			return fmt.Errorf("failed to register admin: %w", err)
		}

		adminUserID, err = s.userGetter.GetUserIDByEmail("admin")
		if err != nil {
			return fmt.Errorf("failed to get admin user ID after registration: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to get admin user ID: %w", err)
	}

	if _, err = s.roleGranter.GrantRole(adminUserID, "admin"); err != nil {
		return fmt.Errorf("failed to grant role: %w", err)
	}

	return nil
}

func NewSeedUseCase(userRegistrator userRegistrator, roleGranter roleGranter, userGetter userGetter) *SeedUseCase {
	return &SeedUseCase{
		userRegistrator: userRegistrator,
		roleGranter:     roleGranter,
		userGetter:      userGetter,
	}
}
