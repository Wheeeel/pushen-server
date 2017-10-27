package config

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

const defaultConfigFile = "config.toml"

var DefaultAppConfig App

type App struct {
	Addr string `toml:"addr"`
	DSN  string `toml:"dsn"`
}

func init() {
	var err error
	DefaultAppConfig, err = Parse(defaultConfigFile)
	if err != nil {
		log.Println("load default config file error")
		return
	}
}

func Parse(file string) (app App, err error) {
	_, err = toml.DecodeFile(file, &app)
	if err != nil {
		err = errors.Wrap(err, "parse config error")
		return
	}
	return
}
