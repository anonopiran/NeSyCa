package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	log "github.com/sirupsen/logrus"
)

var validate *validator.Validate
var lock = &sync.Mutex{}
var settingsInstance *SettingsType

// func Describe() {
// 	var cfg SettingsType
// 	help, err := cleanenv.GetDescription(&cfg, nil)
// 	if err != nil {
// 		log.WithError(err).Panic("can not generate description")
// 	}
// 	log.Println(help)
// }

func Config() *SettingsType {
	if settingsInstance == nil {
		settingsInstance = &SettingsType{}
		lock.Lock()
		defer lock.Unlock()
		k := koanf.New(".")
		_, error := os.Stat(".env")
		if !os.IsNotExist(error) {
			log.Warn("found .env file")
			if err := godotenv.Load(); err != nil {
				panic(err)
			}
		}
		setDefaultValues(k)
		if err := k.Load(env.Provider("NESYCA__", ".", envNormalizer), nil); err != nil {
			panic(err)
		}
		if err := k.Unmarshal("", &settingsInstance); err != nil {
			panic(err)
		}
		validate = validator.New()
		if err := validate.Struct(settingsInstance); err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", *settingsInstance)
	}
	return settingsInstance
}
func envNormalizer(s string) string {
	return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "NESYCA__")), "__", ".", -1)
}

func setDefaultValues(k *koanf.Koanf) {
	def := map[string]interface{}{
		"batch_size_min": "1k",
		"batch_size_max": "10k",
		"file_size_min":  "10m",
		"file_size_max":  "50m",
		"log_level":      "warning",
		"timeout":        "1",
	}
	k.Load(confmap.Provider(def, "."), nil)
}
