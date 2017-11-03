package model

import (
	"log"

	"github.com/Wheeeel/pushen-server/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var DefaultDB *gorm.DB

func init() {
	if config.DefaultAppConfig.DSN != "" {
		var err error
		DefaultDB, err = New(config.DefaultAppConfig.DSN)
		if err != nil {
			log.Fatalf("load default db error: %v", err)
		}
	}
}

func New(dsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		err = errors.Wrap(err, "new db error")
		return
	}
	db.DB().SetMaxIdleConns(1024)
	db.DB().SetMaxOpenConns(1024)

	err = db.DB().Ping()
	if err != nil {
		err = errors.Wrap(err, "new db error")
		return
	}

	db.SingularTable(true)
	return
}
