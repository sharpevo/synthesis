package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

const CONFIG_FILE = "config.yml"

var configFileLoaded = false

func load() {
	if configFileLoaded {
		return
	}
	log.Println("loading config...")
	f, err := os.OpenFile(CONFIG_FILE, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
		return
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(f); err != nil {
		log.Println(err)
	}
	f.Close()
	viper.SetConfigFile(CONFIG_FILE)
	configFileLoaded = true
}

func Set(key string, value interface{}) error {
	load()
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		log.Println(err)
	}
	return err
}

func SetDefault(key string, value interface{}) {
	load()
	viper.SetDefault(key, value)
}

func GetBool(key string) bool {
	load()
	value := viper.GetBool(key)
	if value {
		log.Printf("ENABLE: %s\n", key)
	}
	return value
}

func GetString(key string) string {
	load()
	return viper.GetString(key)
}

func GetInt(key string) int {
	load()
	value := viper.GetInt(key)
	log.Printf("ENABLE: %s = %v\n", key, value)
	return value
}

func GetFloat(key string) float64 {
	load()
	value := viper.GetFloat64(key)
	log.Printf("ENABLE: %s = %v\n", key, value)
	return value
}
