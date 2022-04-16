package redisdao

import (
	"errors"
	goRedis "github.com/go-redis/redis/v8"
)

var RedisClientNotInit = errors.New("redis client not init")

var redisClient *goRedis.Client

func InitClient(cfg *goRedis.Options) {
	// goRedis
	redisClient = goRedis.NewClient(cfg)
	return
}

func GetClient() *goRedis.Client {
	return redisClient
}
