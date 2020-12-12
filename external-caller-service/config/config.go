package config

import (
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	ServiceName string        `mapstructure:"service_name"`
	Restful     RestfulConfig `mapstructure:"restful"`
	Botnoi      BotnoiConfig  `mapstructure:"botnoi"`
}

// RestfulConfig ...
type RestfulConfig struct {
	Port int `mapstructure:"port"`
}

// BotnoiConfig ...
type BotnoiConfig struct {
	Address string `mapstructure:"address"`
	Token   string `mapstructure:"token"`
}

var config *Config
var once sync.Once

// InitConfig just reat it's name
func InitialConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}

		viper.WatchConfig()

		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}

	})
	log.Println("Initialize Config")
	return GetConfig()
}

// GetConfig ...
func GetConfig() *Config {
	return config
}
