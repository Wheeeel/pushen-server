package model

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type BindStatus uint8

const (
	BindStatusNotBinded BindStatus = 0
	BindStatusBinded    BindStatus = 1
)

type BindToken struct {
	ID       int64      `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	UserID   int64      `json:"userId"`
	DeviceID int64      `json:"deviceId"`
	Token    string     `json:"token"`
	Status   BindStatus `json:"status"`

	Timestamp
}

func BindTokenByToken(db *gorm.DB, token string) (ba BindToken, err error) {
	err = db.Where("token = ?", token).First(&ba).Error
	if err != nil {
		err = errors.Wrap(err, "bind token by token error")
		return
	}
	return
}

func BindTokenCreate(db *gorm.DB, token *BindToken) (err error) {
	err = db.Create(token).Error
	if err != nil {
		err = errors.Wrap(err, "bind token create error")
		return
	}
	return
}

func BindTokenUpdateStatus(db *gorm.DB, token *BindToken) (err error) {
	err = db.Table("bind_token").Where("id = ?", token.ID).Update("status", token.Status).Error
	if err != nil {
		err = errors.Wrap(err, "bind token status update error")
		return
	}
	return
}
