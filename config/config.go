package config

var Config = &ServiceConfig{}

type ServiceConfig struct {
	ServiceName      string         `mapstructure:"service_name"`
	HttpServerListen string         `mapstructure:"http_server_listen"`
	CORS             []string       `mapstructure:"cors"`
	Environment      string         `mapstructure:"environment"`
	MongoDBConfig    *MongoDBConfig `mapstructure:"mongo_db_config"`
}

type MongoDBConfig struct {
	URL      string `mapstructure:"url"`
	Database string `mapstructure:"database"`
}
