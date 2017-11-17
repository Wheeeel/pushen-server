package model

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type MessageStatus uint8

const (
	MessageStatusReceived   MessageStatus = 0
	MessageStatusSent       MessageStatus = 1
	MessageStatusSentFailed MessageStatus = 2
)

type Message struct {
	ID       int64         `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	DeviceID int64         `json:"deviceId"`
	AppName  string        `json:"appName"`
	AppIcon  string        `json:"appIcon"`
	Status   MessageStatus `json:"status"`
	Body     string        `json:"messageBody" gorm:"column:body"`

	Timestamp
}

func MessageCreate(db *gorm.DB, message *Message) (err error) {
	err = db.Create(message).Error
	if err != nil {
		err = errors.Wrap(err, "message create error")
		return
	}
	return
}

func MessageOrderByCreateTimestamp(db *gorm.DB) (ms []Message, err error) {
	err = db.Where("status = ?", MessageStatusReceived).Order("create_timestamp asc").Offset(10).Find(&ms).Error
	if err != nil {
		err = errors.Wrap(err, "message by create timestamp error")
		return
	}
	return
}

func MessageByDeviceIDOrderByCreateTimestamp(db *gorm.DB, deviceID int64) (ms []Message, err error) {
	err = db.Where("status = ? AND device_id = ?", MessageStatusReceived, deviceID).
		Order("create_timestamp desc").Offset(10).Find(&ms).Error
	if err != nil {
		err = errors.Wrap(err, "message by create timestamp error")
		return
	}
	return
}

func MessageUpdateStatus(db *gorm.DB, msg *Message) (err error) {
	err = db.Table("message").Where("id = ?", msg.ID).Update("status", msg.Status).Error
	if err != nil {
		err = errors.Wrap(err, "message status update error")
		return
	}
	return
}
