package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
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
	switch conf.Database.Type {
	case "mysql":
		dsn := conf.Database.User + ":" + conf.Database.Password + "@tcp(" + conf.Database.Address + ")/" + conf.Database.DBname + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		break
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(conf.Database.DBname), &gorm.Config{})
	}
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}
