package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

func Initdatabase() (*gorm.DB, error) {
	conf := Config{}
	err := conf.Load()
	if err != nil {
		return nil, err
	}
	var db *gorm.DB
	db, err = gorm.Open("sqlite3", conf.Database)
	if conf.Debug {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}
