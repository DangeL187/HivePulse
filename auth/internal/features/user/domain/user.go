package domain

type User struct {
	ID           uint
	Email        string
	PasswordHash string
	FullName     string
}

type UserInputData struct {
	Email        string
	FullName     string
	PasswordHash string
}
