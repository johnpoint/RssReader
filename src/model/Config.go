package model

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	AllowOrigins []string
	Port         string
	Debug        bool
	Database     Database
	Salt         string
	TLS          bool
	CERTPath     string
	KEYPath      string
}

type Database struct {
	Type     string
	Address  string
	User     string
	Password string
	DBname   string
}

func (c *Config) Load() error {
	_, err := os.Stat("config.json") //os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			return errors.New("config not found")
		}
	}
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		return err
	}
	*c = conf
	return nil
}
