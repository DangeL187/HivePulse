package domain

type Repository interface {
	CreateUser(userInputData UserInputData) (uint, error)
	GetUserByEmail(email string) (User, error)
	UserExists(userID uint) (bool, error)
}
