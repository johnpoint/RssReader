package config

import "time"

var Config = &ServiceConfig{}

type ServiceConfig struct {
	ServiceName      string         `mapstructure:"service_name"`
	HttpServerListen string         `mapstructure:"http_server_listen"`
	CORS             []string       `mapstructure:"cors"`
	Environment      string         `mapstructure:"environment"`
	MongoDBConfig    *MongoDBConfig `mapstructure:"mongo_db_config"`
	RedisConfig      *RedisConfig   `mapstructure:"redis_config"` // redis 配置
}

type MongoDBConfig struct {
	URL      string `mapstructure:"url"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Network            string        `mapstructure:"network"`
	Addr               string        `mapstructure:"addr"`
	Username           string        `mapstructure:"username"`
	Password           string        `mapstructure:"password"`
	DB                 int           `mapstructure:"db"`
	MaxRetries         int           `mapstructure:"max_retries"`
	MinRetryBackoff    time.Duration `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff    time.Duration `mapstructure:"max_retry_backoff"`
	DialTimeout        time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
	PoolFIFO           bool          `mapstructure:"pool_fifo"`
	PoolSize           int           `mapstructure:"pool_size"`
	MinIdleConns       int           `mapstructure:"min_idle_conns"`
	MaxConnAge         time.Duration `mapstructure:"max_conn_age"`
	PoolTimeout        time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout        time.Duration `mapstructure:"idle_timeout"`
	IdleCheckFrequency time.Duration `mapstructure:"idle_check_frequency"`
	readOnly           bool          `mapstructure:"read_only"`
}
