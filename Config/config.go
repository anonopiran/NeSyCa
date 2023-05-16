package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var lock = &sync.Mutex{}
var settingsInstance *SettingsType

func Config() *SettingsType {
	if settingsInstance == nil {
		settingsInstance = &SettingsType{}
		lock.Lock()
		defer lock.Unlock()

		if _, error := os.Stat(".env"); !os.IsNotExist(error) {
			logrus.Warn("found .env file")
			if err := godotenv.Load(); err != nil {
				logrus.Panicf("%+v", err)
			}
		} else {
			logrus.Warn("no .env file found")
		}
		if err := env.Parse(settingsInstance); err != nil {
			logrus.Panicf("%+v", err)
		}
		logrus.WithField("config", fmt.Sprintf("%+v", *settingsInstance)).Warn()
	}
	return settingsInstance
}
