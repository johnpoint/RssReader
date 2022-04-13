package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Feed struct {
	ID         string `json:"_id" bson:"_id"`
	Title      string `json:"title" bson:"title"`
	Url        string `json:"url" bson:"url"`
	Md5        string `json:"md5" bson:"md5"`
	Subscriber int64  `json:"subscriber" bson:"subscriber"`
	UpdateAt   int64  `json:"update_at" bson:"update_at"`
	CreateAt   int64  `json:"create_at" bson:"create_at"`
}

func (m *Feed) CollectionName() string {
	return "feed"
}

func (m *Feed) GetFeeds(ctx context.Context, skip, limit int64) ([]*Feed, error) {
	var feeds []*Feed
	cur, err := DB(m).Find(ctx, bson.M{}, &options.FindOptions{
		Projection: bson.M{
			"_id":       1,
			"url":       1,
			"md5":       1,
			"update_at": 1,
		},
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (m *Feed) UpdateUpdateAtByID(ctx context.Context, updateAt int64) error {
	_, err := DB(m).UpdateOne(ctx, bson.M{
		"_id": m.ID,
	}, bson.M{
		"$set": bson.M{
			"update_at": updateAt,
		},
	})
	return err
}

func (m *Feed) FindByIDs(ctx context.Context, ids []string) ([]*Feed, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	cur, err := DB(m).Find(ctx, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return nil, err
	}

	var feeds []*Feed
	err = cur.All(ctx, &feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (m *Feed) FindByUrl(ctx context.Context, url string) error {
	if len(url) == 0 {
		return nil
	}

	return DB(m).FindOne(ctx, bson.M{
		"url": url,
	}).Decode(m)
}

func (m *Feed) FindByUrls(ctx context.Context, url ...string) ([]*Feed, error) {
	if len(url) == 0 {
		return nil, nil
	}

	cur, err := DB(m).Find(ctx, bson.M{
		"url": bson.M{
			"$in": url,
		},
	})
	if err != nil {
		return nil, err
	}
	var feeds = make([]*Feed, 0)
	err = cur.All(ctx, &feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (m *Feed) InsertOne(ctx context.Context) error {
	_, err := DB(m).InsertOne(ctx, m)
	return err
}
