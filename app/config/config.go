package config

import (
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string `json:"port"`
	RedisPort   string `json:"redis_port"`
	RedisCLient *redis.Client
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

func SetRedisClient(client *redis.Client) {
	LocalConfig.RedisCLient = client
}
