package logic

import (
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/log"
	"RssReader/pkg/utils"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/mmcdole/gofeed"
	"time"
)

var Feed = &FeedL{}

type FeedL struct{}

func (f *FeedL) GetFeed(ctx context.Context, url string) (*mongoModel.Feed, error) {
	if len(url) == 0 {
		return nil, errors.New("url can not be empty")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	feed, err := f.Fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	var feedM = mongoModel.Feed{
		ID:       utils.Sha256(url),
		Title:    feed.Title,
		Url:      url,
		Md5:      utils.Md5(feed.String()),
		CreateAt: time.Now().UnixMilli(),
		UpdateAt: time.Now().UnixMilli(),
	}
	err = feedM.InsertOne(ctx)
	if err != nil {
		return nil, err
	}
	f.UpdatePost(ctx, feed, 0, feedM.ID)
	return &feedM, nil
}

func (f *FeedL) UpdateFeed(ctx context.Context, feedM *mongoModel.Feed) error {
	if len(feedM.Url) == 0 {
		return errors.New("url can not be empty")
	}
	log.Info("FeedL.UpdateFeed", log.String("u", feedM.Url))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	feed, err := f.Fetch(ctx, feedM.Url)
	if err != nil {
		return err
	}
	f.UpdatePost(ctx, feed, feedM.UpdateAt, feedM.ID)
	return err
}

var ignoreCharMap = map[string]struct{}{
	"U+001B": {},
	"U+0008": {},
}

func (f *FeedL) Fetch(ctx context.Context, url string) (*gofeed.Feed, error) {
	var b string
	err := gout.GET(url).BindBody(&b).Do()
	if err != nil {
		return nil, err
	}
	var n bytes.Buffer
	for _, v := range b {
		if _, has := ignoreCharMap[fmt.Sprintf("%U", v)]; !has {
			n.Write([]byte(fmt.Sprintf("%c", v)))
		}
	}
	parser := gofeed.NewParser()
	feed, err := parser.Parse(&n)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

var parseTimeFmt = []string{
	time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123,
	time.RFC1123Z, time.RFC3339, time.RFC3339Nano,
}

func (f *FeedL) UpdatePost(ctx context.Context, feed *gofeed.Feed, newerThan int64, feedID string) {
	var posts []interface{}
	for i := range feed.Items {
		var t time.Time
		for _, v := range parseTimeFmt {
			var err error
			t, err = time.Parse(v, feed.Items[i].Published)
			if err == nil {
				break
			}
			t = time.Now()
		}
		if t.UnixMilli() > newerThan {
			var p mongoModel.Post
			p.Url = feed.Items[i].Link
			p.ID = utils.Md5(p.Url)
			p.FID = feedID
			p.Title = feed.Items[i].Title
			p.Content = feed.Items[i].Content
			p.Description = feed.Items[i].Description
			p.Published = t.UnixMilli()
			posts = append(posts, &p)
		}
	}
	err := new(mongoModel.Post).InsertMany(ctx, posts)
	if err != nil {
		return
	}

	var feedM = mongoModel.Feed{
		ID: feedID,
	}
	err = feedM.UpdateUpdateAtByID(ctx, time.Now().UnixMilli())
	if err != nil {
		return
	}
}
