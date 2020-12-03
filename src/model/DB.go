package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func init() {
	if db == nil {
		var err error
		db, err = InitDatabase()
		if err != nil {
			log.Println(err.Error())
		} else {
			if err := db.AutoMigrate(&Feed{}, &Post{}, &User{}, &Read{}, &subscribe{}, Log{}); err != nil {
				panic(err)
			}
		}
	}
}

func InitDatabase() (*gorm.DB, error) {
	conf := Config{}
	err := conf.Load()
	if err != nil {
		return nil, err
	}
	var db *gorm.DB
	var logInfo logger.Interface
	if conf.Debug {
		logInfo = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // Log level
				Colorful:      false,       // 禁用彩色打印
			},
		)
	} else {
		logInfo = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // 禁用彩色打印
			},
		)
	}
	switch conf.Database.Type {
	case "mysql":
		dsn := conf.Database.User + ":" + conf.Database.Password + "@tcp(" + conf.Database.Address + ")/" + conf.Database.DBname + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logInfo,
		})
		break
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(conf.Database.DBname), &gorm.Config{
			Logger: logInfo,
		})
	}
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}
