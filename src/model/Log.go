package model

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Log struct {
	ID      int64 `gorm:"autoIncrement"`
	Level   int64
	Message string
	Type    string
	Created int64 `gorm:"autoCreateTime"`
}

func (l *Log) New() error {
	if db == nil {
		file := strconv.FormatInt(time.Now().Unix(), 10) + ".log"
		fp, err := os.Create("log/" + file)
		if err != nil {
			log.Println("Error report could not be saved")
			log.Println(l.Message)
			return errors.New("Error report could not be saved")
		}
		logdata, _ := json.Marshal(l)
		_, _ = io.WriteString(fp, string(logdata))
		log.Println("The error report cannot be saved to the database, but it has been saved to a file: " + file)
		return errors.New("The error report cannot be saved to the database, but it has been saved to a file: " + file)
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	_ = tx.AutoMigrate(&Log{})
	if err := tx.Create(&l).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
