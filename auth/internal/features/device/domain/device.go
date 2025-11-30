package domain

type Device struct {
	ID           uint
	DeviceID     string
	PasswordHash string
}

type DeviceInputData struct {
	DeviceID     string
	PasswordHash string
}
