package infra

import (
	"errors"
	"gorm.io/gorm"

	"github.com/DangeL187/erax"
	"github.com/jackc/pgx/v5/pgconn"

	"auth/internal/features/user/domain"
)

type UserRepo struct {
	db *gorm.DB
}

func (ur *UserRepo) CreateUser(userInputData domain.UserInputData) (uint, error) {
	user := domain.User{
		Email:        userInputData.Email,
		FullName:     userInputData.FullName,
		PasswordHash: userInputData.PasswordHash,
	}

	if err := ur.db.Create(&user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return 0, erax.WrapWithError(err, domain.ErrEmailAlreadyExists, "failed to insert value")
			}
		}

		return 0, erax.Wrap(err, "failed to insert value")
	}

	return user.ID, nil
}

func (ur *UserRepo) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User

	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, erax.WrapWithError(err, domain.ErrUserNotFound, "failed to query user by email")
		}

		return domain.User{}, erax.Wrap(err, "failed to query user by email")
	}

	return user, nil
}

func (ur *UserRepo) UserExists(userID uint) (bool, error) {
	var user domain.User
	err := ur.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, erax.Wrap(err, "failed to query user by user ID")
	}

	return true, nil
}

func NewUserRepo(db *gorm.DB) domain.Repository {
	return &UserRepo{db: db}
}
