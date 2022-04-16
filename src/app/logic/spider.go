package logic

import (
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/log"
	"context"
	"time"
)

type FeedSpider struct {
	ctx context.Context
}

func NewFeedSpider(ctx context.Context) *FeedSpider {
	return &FeedSpider{
		ctx: ctx,
	}
}

func (f *FeedSpider) Loop() {
	for {
		log.Info("Loop.start")
		err := f.Catch()
		if err != nil {
			log.Error("FeedSpider.Loop", log.Err(err))
		}
		log.Info("Loop.finish")
		time.Sleep(30 * time.Minute)
	}
}

func (f *FeedSpider) Catch() error {
	var page int64
	var pageSize int64 = 5
	for {
		feeds, err := new(mongoModel.Feed).GetFeeds(f.ctx, (page-1)*pageSize, pageSize)
		page++
		if err != nil {
			continue
		}
		if len(feeds) == 0 {
			return nil
		}
		for _, v := range feeds {
			err := new(FeedL).UpdateFeed(f.ctx, v)
			if err != nil {
				continue
			}
		}
	}

}
