package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type User struct {
	ID                 int64     `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	Password           string    `json:"-"`
	AvatarURL          string    `json:"avatarUrl"`
	LastLoginTimestamp time.Time `json:"lastLoginTimestamp"`

	Timestamp
}

func UserByID(db *gorm.DB, id int64) (user User, err error) {
	err = db.Where("id = ?", id).First(&user).Error
	if err != nil {
		err = errors.Wrap(err, "user by id error")
		return
	}
	return
}

func UserByEmail(db *gorm.DB, email string) (user User, err error) {
	err = db.Where("email = ?", email).First(&user).Error
	if err != nil {
		err = errors.Wrap(err, "user by email error")
		return
	}
	return
}

func UserCreate(db *gorm.DB, user *User) (err error) {
	err = db.Create(user).Error
	if err != nil {
		err = errors.Wrap(err, "user create error")
		return
	}
	return
}

func UserValidate(db *gorm.DB, email, password string) (ok bool, err error) {
	users := make([]User, 0, 1)
	err = db.Where("email = ? AND password = ?", email, password).Find(&users).Error
	if err != nil {
		err = errors.Wrap(err, "user validate error")
		return
	}
	if len(users) == 1 {
		ok = true
	}
	return
}
