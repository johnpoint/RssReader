package model

import (
	"encoding/json"
	"os"
)

type Config struct {
	AllowOrigins []string
	Port         string
	Debug        bool
	Database      string
	Salt         string
	TLS          bool
	CERTPath     string
	KEYPath      string
}

func (c *Config) Load() error {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Config{}
	err := decoder.Decode(&conf)
	if err != nil {
		return err
	}
	*c = conf
	return nil
}
