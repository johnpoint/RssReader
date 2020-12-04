package model

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"
)

type Feed struct {
	ID     int64 `gorm:"autoIncrement"`
	Title  string
	Url    string
	Feed   string
	Num    int64
	Status int64 `gorm:"default:0"` // 0 OK   10 ERROR
	Md5    string
	Posts  []Post `gorm:"foreignKey:FID;constraint:OnDelete:CASCADE;"`
}

func (f *Feed) Get(selects []string) error {
	if f.ID == 0 && f.Url == "" {
		return errors.New("url and ID can not be empty")
	}
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	var Feeds []Feed
	if selects == nil {
		db.Where(Feed{Url: f.Url, ID: f.ID}).Find(&Feeds)
	} else {
		db.Select(selects).Where(Feed{Url: f.Url, ID: f.ID}).Find(&Feeds)
	}
	if len(Feeds) == 0 {
		return errors.New("not Found")
	}
	*f = Feeds[0]
	return nil
}

func (f *Feed) New() error {
	if f.Url == "" {
		return errors.New("url can not be empty")
	}
	err := f.Get([]string{"id"})
	if err == nil {
		return errors.New("feed already exists")
	}
	err = nil
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(f.Url, ctx)
	if err != nil {
		return err
	}
	f.Title = feed.Title
	f.Feed = feed.String()
	r := md5.Sum([]byte(f.Feed))
	f.Md5 = hex.EncodeToString(r[:])
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

	if err := tx.Create(&f).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	_ = f.Update()
	return nil
}

func (f *Feed) Update() error {
	if f.Url == "" && f.ID == 0 {
		return errors.New("url and ID can not be empty")
	}
	err := f.Get([]string{"id", "url", "md5"})
	if err != nil {
		err := f.New()
		if err != nil {
			return err
		}
	}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(f.Url)
	if err != nil {
		return err
	}
	conf := Config{}
	err = conf.Load()
	if err != nil {
		return err
	}
	r := md5.Sum([]byte(feed.String()))
	if f.Md5 == hex.EncodeToString(r[:]) {
		return nil
	}
	f.Md5 = hex.EncodeToString(r[:])
	for _, i := range feed.Items {
		p := Post{}
		p.Url = i.Link
		errGet := p.Get([]string{"id"})
		if errGet == nil {
			continue
		}
		/*
				RFC822      = "02 Jan 06 15:04 MST"
			    RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
			    RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
			    RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
			    RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
			    RFC3339     = "2006-01-02T15:04:05Z07:00"
			    RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
		*/
		t, err := time.Parse(time.RFC822, i.Published)
		if err != nil {
			err = nil
			t, err = time.Parse(time.RFC822Z, i.Published)
		}
		if err != nil {
			err = nil
			t, err = time.Parse(time.RFC850, i.Published)
		}
		if err != nil {
			err = nil
			t, err = time.Parse(time.RFC1123, i.Published)
		}
		if err != nil {
			err = nil
			t, err = time.Parse(time.RFC1123Z, i.Published)
		}
		if err != nil {
			err = nil
			t, err = time.Parse(time.RFC3339, i.Published)
		}
		if err != nil {
			err = nil
			t, err = time.Parse(time.RFC3339Nano, i.Published)
		}
		if err != nil {
			err = nil
			t = time.Now()
		}
		p.FID = f.ID
		p.Title = i.Title
		p.Content = i.Content
		p.Description = i.Description
		p.Published = strconv.FormatInt(t.Unix(), 10)
		if err := p.New(); err != nil {
			continue
		}
	}
	f.Feed = feed.String()
	err = f.save()
	return err
}

func (f *Feed) Save() error {
	err := f.save()
	if err != nil {
		return err
	}
	return nil
}

func (f *Feed) save() error {
	if db == nil {
		return errors.New("database connection failed")
	}
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
	where := Feed{ID: f.ID}
	if err := tx.Model(&where).Where(where).Updates(f).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (f *Feed) Post(limit int) []Post {
	if db == nil {
		return []Post{}
	}
	var Posts []Post
	db.Where(Post{FID: f.ID}).Find(&Posts).Order("published").Limit(limit)
	return Posts
}

func (f *Feed) Delete() error {
	if f.Num != 0 {
		return errors.New("feed can not be delete")
	}
	if f.ID == 0 {
		return errors.New("ID can not be empty")
	}
	if db == nil {
		return errors.New("database connection failed")
	}
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
	if err := tx.Delete(f).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
