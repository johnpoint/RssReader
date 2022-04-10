package depend

import (
	"RssReader/config"
	"RssReader/dao/mongoDao"
	"RssReader/pkg/bootstrap"
	"context"
)

type MongoDB struct{}

var _ bootstrap.Component = (*MongoDB)(nil)

func (r *MongoDB) Init(ctx context.Context) error {
	mongoDao.InitMongoClient(config.Config.MongoDBConfig)
	err := mongoDao.MongoClient.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
