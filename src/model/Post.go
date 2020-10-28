package model

import (
	"errors"
)

type Post struct {
	ID          int64 `gorm:"autoIncrement"`
	FID         int64 `gorm:"primaryKey"`
	Title       string
	Content     string
	Url         string `gorm:"primaryKey"`
	Description string
	Published   string
}

func (p *Post) Get() error {
	if p.FID == 0 && p.ID == 0 && p.Url == "" {
		return errors.New("Incomplete parameters")
	}
	if db == nil {
		return errors.New("Database connection failed")
	}
	// defer db.Close()
	_ = db.AutoMigrate(&Post{})
	Posts := []Post{}
	db.Where(Post{ID: p.ID, FID: p.FID, Url: p.Url}).Find(&Posts)
	if len(Posts) == 0 {
		return errors.New("Not Found")
	}
	p.Content = Posts[0].Content
	p.ID = Posts[0].ID
	p.FID = Posts[0].FID
	p.Published = Posts[0].Published
	p.Description = Posts[0].Description
	p.Title = Posts[0].Title
	p.Url = Posts[0].Url
	return nil
}

func (p *Post) FeedPost(where []int64, limit int) []Post {
	if db == nil {
		return []Post{}
	}
	// defer db.Close()
	_ = db.AutoMigrate(&Post{})
	Posts := []Post{}
	db.Where(map[string]interface{}{"f_id": where}).Order("published desc").Limit(limit).Find(&Posts)
	if len(Posts) == 0 {
		return []Post{}
	}
	return Posts
}

func (p *Post) New() error {
	if p.FID == 0 || p.Url == "" || p.Title == "" || p.Published == "" {
		return errors.New("Incomplete parameters")
	}
	err := p.Get()
	if err == nil {
		return errors.New("Post already exist")
	}
	if db == nil {
		return errors.New("Database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}

	_ = tx.AutoMigrate(&Post{})
	if err := tx.Create(&p).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (p *Post) Delete() error {
	if p.ID == 0 {
		return errors.New("Incomplete parameters")
	}
	if db == nil {
		return errors.New("Database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}
	_ = tx.AutoMigrate(&Post{})
	if err := tx.Where(Post{ID: p.ID}).Delete(Post{}).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (p *Post) save() error {
	if db == nil {
		return errors.New("Database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}
	_ = tx.AutoMigrate(&Post{})
	where := Post{ID: p.ID}
	if err := tx.Model(&where).Where(where).Updates(p).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
