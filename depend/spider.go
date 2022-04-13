package depend

import (
	"RssReader/app/logic"
	"RssReader/pkg/bootstrap"
	"context"
)

type Spider struct{}

var _ bootstrap.Component = (*Spider)(nil)

func (d *Spider) Init(ctx context.Context) error {
	go logic.NewFeedSpider(ctx).Loop()
	return nil
}
