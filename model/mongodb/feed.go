package mongodb

import (
	"RssReader/dao/mongoDao"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type Feed struct {
	ID    string `json:"_id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Url   string `json:"url" bson:"url"`
}

func (m *Feed) CollectionName() string {
	return "feed"
}

func (m *Feed) FindByIDs(ctx context.Context, ids []string, feeds []Feed) error {
	if len(ids) == 0 {
		return nil
	}
	return mongoDao.Client(m.CollectionName()).FindOne(ctx, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}).Decode(feeds)
}
