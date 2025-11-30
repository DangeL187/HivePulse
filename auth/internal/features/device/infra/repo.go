package infra

import (
	"errors"
	"gorm.io/gorm"

	"github.com/DangeL187/erax"
	"github.com/jackc/pgx/v5/pgconn"

	"auth/internal/features/device/domain"
)

type DeviceRepo struct {
	db *gorm.DB
}

func (dr *DeviceRepo) CreateDevice(deviceInputData domain.DeviceInputData) (uint, error) {
	device := domain.Device{
		DeviceID:     deviceInputData.DeviceID,
		PasswordHash: deviceInputData.PasswordHash,
	}

	if err := dr.db.Create(&device).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "devices_device_id_key" {
				return 0, erax.WrapWithError(err, domain.ErrDeviceAlreadyExists, "failed to insert value")
			}
		}

		return 0, erax.Wrap(err, "failed to insert value")
	}

	return device.ID, nil
}

func (dr *DeviceRepo) DeviceExists(deviceID string) (bool, error) {
	var device domain.Device
	err := dr.db.First(&device, deviceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, erax.Wrap(err, "failed to query device by ID")
	}

	return true, nil
}

func (dr *DeviceRepo) GetDeviceByDeviceID(deviceID string) (domain.Device, error) {
	var device domain.Device

	if err := dr.db.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Device{}, erax.WrapWithError(err, domain.ErrDeviceNotFound, "failed to query device by device ID")
		}

		return domain.Device{}, erax.Wrap(err, "failed to query device by device ID")
	}

	return device, nil
}

func NewDeviceRepo(db *gorm.DB) *DeviceRepo {
	return &DeviceRepo{db: db}
}
