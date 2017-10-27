package model

import (
	"time"

	"github.com/pkg/errors"
)

type User struct {
	ID                 int64     `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	Password           string    `json:"password"`
	AvatarURL          string    `json:"avatarUrl"`
	LastLoginTimestamp time.Time `json:"lastLoginTimestamp"`

	Timestamp
}

func UserByEmail(email string) (user User, err error) {
	err = DefaultDB.Where("email = ?", email).First(&user).Error
	if err != nil {
		err = errors.Wrap(err, "user by email error")
		return
	}
	return
}

func UserCreate(user *User) (err error) {
	err = DefaultDB.Create(user).Error
	if err != nil {
		err = errors.Wrap(err, "user create error")
		return
	}
	return
}

func UserValidate(email, password string) (ok bool, err error) {
	users := make([]User, 0, 1)
	err = DefaultDB.Where("email = ? AND password = ?", email, password).Find(&users).Error
	if err != nil {
		err = errors.Wrap(err, "user validate error")
		return
	}
	if len(users) == 1 {
		ok = true
	}
	return
}
