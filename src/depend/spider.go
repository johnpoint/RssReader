package depend

import (
	"RssReader/app/logic"
	"RssReader/pkg/bootstrap"
	"RssReader/pkg/log"
	"context"
	"github.com/robfig/cron/v3"
)

type Spider struct{}

var _ bootstrap.Component = (*Spider)(nil)

func (d *Spider) Init(ctx context.Context) error {
	tasks := cron.New()
	_, err := tasks.AddFunc("@every 30m", AutoRefreshFeed)
	if err != nil {
		return err
	}
	tasks.Run()
	return nil
}

// AutoRefreshFeed 刷新所有 Feed
func AutoRefreshFeed() {
	log.Info("AutoRefreshFeed-Start")
	ctx := context.TODO()
	err := logic.NewFeedSpider(ctx).Catch()
	if err != nil {
		log.Error("AutoRefreshFeed", log.Err(err))
		return
	}
	log.Info("AutoRefreshFeed-End")
}
