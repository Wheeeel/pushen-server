package model

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type DeviceStatus uint8

const (
	DeviceStatusInit     DeviceStatus = 0
	DeviceStatusBinded   DeviceStatus = 1
	DeviceStatusUnbinded DeviceStatus = 2
)

type Device struct {
	ID     int64        `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	UserID int64        `json:"userId"`
	Type   string       `json:"type"`
	UUID   string       `json:"uuid"`
	Status DeviceStatus `json:"status"`

	Timestamp
}

func DeviceCreate(db *gorm.DB, device *Device) (err error) {
	err = db.Create(device).Error
	if err != nil {
		err = errors.Wrap(err, "device create error")
		return
	}
	return
}

func DevicesByUserID(db *gorm.DB, id int64) (devices []Device, err error) {
	err = db.Where("id = ?", id).Find(&devices).Error
	if err != nil {
		err = errors.Wrap(err, "device by user id error")
		return
	}
	return
}

func DeviceByUUID(db *gorm.DB, uuid string) (device Device, err error) {
	err = db.Where("uuid = ?", uuid).First(&device).Error
	if err != nil {
		err = errors.Wrap(err, "device by uuid error")
		return
	}
	return
}

func DeviceUpdateStatus(db *gorm.DB, device *Device) (err error) {
	err = db.Table("device").Where("uuid = ?", device.UUID).Update("status", device.Status).Error
	if err != nil {
		err = errors.Wrap(err, "device status update error")
		return
	}
	return
}
