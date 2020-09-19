package model

import (
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
)

type Feed struct {
	ID    int64
	Title string
	Url   string
	Feed  string
	Num   int64
}

func (f *Feed) Get() error {
	if f.ID == 0 && f.Url == "" {
		return errors.New("Url and ID can not be empty")
	}
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	_ = db.AutoMigrate(&User{})
	Feeds := []Feed{}
	db.Where(Feed{Url: f.Url, ID: f.ID}).Find(&Feeds)
	if len(Feeds) == 0 {
		return errors.New("Not Found")
	}
	f.ID = Feeds[0].ID
	f.Title = Feeds[0].Title
	f.Url = Feeds[0].Url
	f.Feed = Feeds[0].Feed
	f.Num = Feeds[0].Num
	return nil
}

func (f *Feed) New() error {
	if f.Url == "" {
		return errors.New("Url can not be empty")
	}
	err := f.Get()
	if err == nil {
		return errors.New("Feed already exists")
	}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(f.Url)
	if err != nil {
		return err
	}
	f.Title = feed.Title
	f.Feed = feed.String()
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	_ = tx.AutoMigrate(&Feed{})
	if err := tx.Create(f).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	f.Update()
	return nil
}

func (f *Feed) Update() error {
	if f.Url == "" && f.ID == 0 {
		return errors.New("Url and ID can not be empty")
	}
	err := f.Get()
	if err != nil {
		err := f.New()
		return err
	}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(f.Url)
	if err != nil {
		return err
	}
	p := Post{}
	conf := Config{}
	err = conf.Load()
	if err != nil {
		return err
	}
	for _, i := range feed.Items {
		p.ID = 0
		p.FID = f.ID
		p.Title = i.Title
		p.Content = i.Content
		p.Description = i.Description
		p.Url = i.Link
		p.Published = i.Published
		p.New()
	}
	f.Feed = feed.String()
	err = f.save()
	return err
}

func (f *Feed) save() error {
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	_ = tx.AutoMigrate(&Feed{})
	where := Feed{ID: f.ID}
	if err := tx.Model(&where).Where(where).Update(f).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (f *Feed) Post() []Post {
	db, err := Initdatabase()
	if err != nil {
		return []Post{}
	}
	defer db.Close()
	_ = db.AutoMigrate(&Post{})
	Posts := []Post{}
	db.Where(Post{FID: f.ID}).Find(&Posts)
	return Posts
}

func (f *Feed) Detele() error {
	if f.Num != -1 {
		return errors.New("Feed can not be delete")
	}
	if f.ID == 0 {
		return errors.New("ID can not be empty")
	}
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	_ = tx.AutoMigrate(&Feed{})
	where := Feed{ID: f.ID}
	if err := tx.Where(where).Delete(f).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	p := f.Post()
	for _, i := range p {
		_ = i.Delete()
	}
	return nil
}
