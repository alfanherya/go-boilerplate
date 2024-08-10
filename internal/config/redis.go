package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(config *viper.Viper) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetString("redis.host"),
		Password: config.GetString("redis.password"),
		DB:       config.GetInt("redis.db"),
	})

	return client
}
