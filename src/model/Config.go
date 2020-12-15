package model

import (
	"encoding/json"
	"fmt"
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

func (c *Config) Load(config string) error {
	_, err := os.Stat(config) //os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			fmt.Println("The configuration file does not exist")
			os.Exit(1)
		}
	}
	file, _ := os.Open(config)
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
