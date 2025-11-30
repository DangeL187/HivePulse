package domain

type Repository interface {
	CreateDevice(deviceInputData DeviceInputData) (uint, error)
	DeviceExists(deviceID string) (bool, error)
	GetDeviceByDeviceID(deviceID string) (Device, error)
}
