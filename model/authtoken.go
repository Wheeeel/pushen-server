package model

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type AuthToken struct {
	ID       int64  `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	UserID   int64  `json:"userId"`
	DeviceID int64  `json:"deviceId"`
	Token    string `json:"token"`

	Timestamp
}

func AuthTokenByToken(db *gorm.DB, token string) (at AuthToken, err error) {
	err = db.Where("token = ?", token).First(&at).Error
	if err != nil {
		err = errors.Wrap(err, "auth token by token error")
		return
	}
	return
}

func AuthTokenCreate(db *gorm.DB, token *AuthToken) (err error) {
	err = db.Create(token).Error
	if err != nil {
		err = errors.Wrap(err, "auth token create error")
		return
	}
	return
}

func AuthTokenUpdateStatus(token *BindToken) (err error) {
	err = DefaultDB.Table("bind_token").Where("id = ?", token.ID).Update("status", token.Status).Error
	if err != nil {
		err = errors.Wrap(err, "bind token status update error")
		return
	}
	return
}
