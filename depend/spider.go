package depend

import (
	"RssReader/pkg/bootstrap"
	"context"
)

type Spider struct {
	Path string
}

var _ bootstrap.Component = (*Spider)(nil)

func (d *Spider) Init(ctx context.Context) error {
	return nil
}
