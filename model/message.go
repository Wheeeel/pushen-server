package model

import "github.com/pkg/errors"

type MessageStatus uint8

const (
	MessageStatusReceived   MessageStatus = 0
	MessageStatusSendt      MessageStatus = 1
	MessageStatusSendFailed MessageStatus = 2
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

func MessageCreate(message *Message) (err error) {
	err = DefaultDB.Create(message).Error
	if err != nil {
		err = errors.Wrap(err, "message create error")
		return
	}
	return
}

func MessageByCreateTimestamp() (ms []Message, err error) {
	err = DefaultDB.Where("status = ?", MessageStatusReceived).
		Order("create_timestamp desc").Offset(10).Find(&ms).Error
	if err != nil {
		err = errors.Wrap(err, "message by create timestamp error")
		return
	}
	return
}

func MessageUpdateStatus(msg *Message) (err error) {
	err = DefaultDB.Table("message").Where("id = ?", msg.ID).Update("status", msg.Status).Error
	if err != nil {
		err = errors.Wrap(err, "message status update error")
		return
	}
	return
}
