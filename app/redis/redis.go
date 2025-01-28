package redis

import (
	"github.com/google/logger"
	"github.com/redis/go-redis/v9"
	"github.com/saadrupai/go-message-broker/app/config"
)

func ConnectRedis() *redis.Client {
	redisPort := config.LocalConfig.RedisPort

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + redisPort,
		Password: "",
		DB:       0,
		Protocol: 2, // connection protocol
	})

	logger.Info("successfully connected to redis db")

	return client

}
