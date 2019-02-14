package util

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime/debug"
)

func Go(function func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf(
					"Panic: %s\n%s",
					err, debug.Stack(),
				)
				errf, err := os.OpenFile("posam.err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
				if err != nil {
					log.Println("Error when catching panic:", err)
					return
				}
				errf.WriteString(msg)
				errf.Sync()
				errf.Close()
				log.Fatal(msg)
			}
		}()
		function()
	}()
}

const CONFIG_FILE = "config.yml"

var configFileLoaded = false

func LoadConfig() {
	if configFileLoaded {
		return
	}
	fmt.Println("loading config...")
	f, err := os.Open(CONFIG_FILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(f); err != nil {
		fmt.Println(err)
	}
	f.Close()
	viper.SetConfigFile(CONFIG_FILE)
	configFileLoaded = true
}

func GetBool(key string) bool {
	LoadConfig()
	value := viper.GetBool(key)
	if value {
		fmt.Printf("ENABLE: %s\n", key)
	}
	return value
}

func Set(key string, value interface{}) error {
	LoadConfig()
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func GetString(key string) string {
	LoadConfig()
	return viper.GetString(key)
}

func GetInt(key string) int {
	LoadConfig()
	value := viper.GetInt(key)
	fmt.Printf("ENABLE: %s = %v\n", key, value)
	return value
}
