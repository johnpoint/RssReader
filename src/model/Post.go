package model

import (
	"errors"
)

type Post struct {
	ID          int64 `gorm:"autoIncrement"`
	FID         int64
	Title       string
	Content     string
	Url         string
	Description string
	Published   string
	Reads       []Read `gorm:"foreignKey:PID;constraint:OnDelete:CASCADE;"`
}

func (p *Post) Get(selects []string) error {
	if p.FID == 0 && p.ID == 0 && p.Url == "" {
		return errors.New("incomplete parameters")
	}
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	var Posts []Post
	if selects == nil {
		db.Where(Post{ID: p.ID, FID: p.FID, Url: p.Url}).Find(&Posts)
	} else {
		db.Select(selects).Where(Post{ID: p.ID, FID: p.FID, Url: p.Url}).Find(&Posts)
	}

	if len(Posts) == 0 {
		return errors.New("not Found")
	}
	*p = Posts[0]
	return nil
}

func (p *Post) FeedPost(where []int64, limit int) []Post {
	if db == nil {
		return []Post{}
	}
	// defer db.Close()
	var Posts []Post
	db.Where(map[string]interface{}{"f_id": where}).Order("published desc").Limit(limit).Find(&Posts)
	if len(Posts) == 0 {
		return []Post{}
	}
	return Posts
}

func (p *Post) New() error {
	if p.FID == 0 || p.Url == "" || p.Title == "" || p.Published == "" {
		return errors.New("incomplete parameters")
	}
	err := p.Get([]string{"id"})
	if err == nil {
		return errors.New("post already exist")
	}
	if db == nil {
		return errors.New("database connection failed")
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
		return errors.New("incomplete parameters")
	}
	if db == nil {
		return errors.New("database connection failed")
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
		return errors.New("database connection failed")
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
	where := Post{ID: p.ID}
	if err := tx.Model(&where).Where(where).Updates(&p).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
