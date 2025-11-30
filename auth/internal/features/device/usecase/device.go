package usecase

type deviceExistChecker interface {
	DeviceExists(deviceID string) (bool, error)
}

type DeviceUseCase struct {
	repo deviceExistChecker
}

func (d *DeviceUseCase) DeviceExists(deviceID string) (bool, error) {
	exists, err := d.repo.DeviceExists(deviceID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func NewDeviceUseCase(repo deviceExistChecker) *DeviceUseCase {
	return &DeviceUseCase{
		repo: repo,
	}
}
