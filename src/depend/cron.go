package depend

import (
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/bootstrap"
	"RssReader/pkg/log"
	"context"
	"github.com/robfig/cron/v3"
	"time"
)

type Cron struct{}

var _ bootstrap.Component = (*Cron)(nil)

func (c *Cron) Init(ctx context.Context) error {
	tasks := cron.New()
	_, err := tasks.AddFunc("0 1 * * *", AutoRemovePost)
	if err != nil {
		return err
	}
	tasks.Run()
	return nil
}

// AutoRemovePost 自动删除半年前的POST
func AutoRemovePost() {
	log.Info("AutoRemovePost-Start")
	ctx := context.TODO()
	err := new(mongoModel.Post).AutoRemovePost(ctx, time.Now().Add(6*30*24*time.Hour).UnixMilli())
	if err != nil {
		log.Error("AutoRemovePost", log.Err(err))
		return
	}
	log.Info("AutoRemovePost-End")
}
