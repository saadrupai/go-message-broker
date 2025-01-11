package config

import (
	"log"
	"net"

	"github.com/spf13/viper"
)

type Config struct {
	Port             string `json:"port"`
	SockerServerPort string `json:socket_server_port`
	Listener         net.Listener
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to load env variables")
	}

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("failed to load env variables")
	}

	return config
}

var LocalConfig *Config

func SetConfig() {
	LocalConfig = LoadConfig()
}
