package depend

import (
	"RssReader/config"
	"RssReader/dao/redisdao"
	"RssReader/pkg/bootstrap"
	"context"
	goRedis "github.com/go-redis/redis/v8"
)

// Redis 初始化 Redis 客户端
type Redis struct{}

var _ bootstrap.Component = (*Redis)(nil)

func (d *Redis) Init(ctx context.Context) error {
	redisdao.InitClient(&goRedis.Options{
		Network:            config.Config.RedisConfig.Network,
		Addr:               config.Config.RedisConfig.Addr,
		Username:           config.Config.RedisConfig.Username,
		Password:           config.Config.RedisConfig.Password,
		DB:                 config.Config.RedisConfig.DB,
		MaxRetries:         config.Config.RedisConfig.MaxRetries,
		MinRetryBackoff:    config.Config.RedisConfig.MinRetryBackoff,
		MaxRetryBackoff:    config.Config.RedisConfig.MaxRetryBackoff,
		DialTimeout:        config.Config.RedisConfig.DialTimeout,
		ReadTimeout:        config.Config.RedisConfig.ReadTimeout,
		WriteTimeout:       config.Config.RedisConfig.WriteTimeout,
		PoolFIFO:           config.Config.RedisConfig.PoolFIFO,
		PoolSize:           config.Config.RedisConfig.PoolSize,
		MinIdleConns:       config.Config.RedisConfig.MinIdleConns,
		MaxConnAge:         config.Config.RedisConfig.MaxConnAge,
		PoolTimeout:        config.Config.RedisConfig.PoolTimeout,
		IdleTimeout:        config.Config.RedisConfig.IdleTimeout,
		IdleCheckFrequency: config.Config.RedisConfig.IdleCheckFrequency,
	})
	return nil
}
