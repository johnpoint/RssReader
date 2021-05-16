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

var Db *gorm.DB

func InitGorm() {
	if Db == nil {
		var err error
		Db, err = initDatabase()
		if err != nil {
			log.Println(err.Error())
		} else {
			if err := Db.AutoMigrate(&Feed{}, &Post{}, &User{}, &Read{}, &subscribe{}, Log{}); err != nil {
				panic(err)
			}
		}
	}
}

func initDatabase() (*gorm.DB, error) {
	conf := Cfg
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
	var err error
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
