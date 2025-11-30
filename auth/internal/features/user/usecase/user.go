package usecase

import (
	"github.com/DangeL187/erax"

	"auth/internal/features/user/domain"
)

type UserRepo interface {
	GetUserByEmail(email string) (domain.User, error)
	UserExists(userID uint) (bool, error)
}

type UserUseCase struct {
	repo UserRepo
}

func (u *UserUseCase) UserExists(userID uint) (bool, error) {
	exists, err := u.repo.UserExists(userID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (u *UserUseCase) GetUserIDByEmail(email string) (uint, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return 0, erax.Wrap(err, "failed to get user by email")
	}

	return user.ID, nil
}

func NewUserUseCase(repo domain.Repository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}
