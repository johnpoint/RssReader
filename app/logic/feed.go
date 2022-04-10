package logic

import (
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/utils"
	"context"
	"errors"
	"github.com/mmcdole/gofeed"
	"time"
)

var Feed = &FeedL{}

type FeedL struct{}

func (l *FeedL) GetFeed(ctx context.Context, url string) (*mongoModel.Feed, error) {
	if len(url) == 0 {
		return nil, errors.New("url can not be empty")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	parser := gofeed.NewParser()
	feed, err := parser.ParseURLWithContext(url, ctx)
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
	return &feedM, nil
}
