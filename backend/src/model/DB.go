package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

func Initdatabase() *gorm.DB {
	conf := Config{}
	conf.Load()
	var err error
	var db *gorm.DB
	db, err = gorm.Open("sqlite3", conf.Database)
	if conf.Debug {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return db
}
