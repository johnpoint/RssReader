package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID          string `json:"_id" bson:"_id"`
	FID         string `json:"fid" bson:"fid"`
	Title       string `json:"title" bson:"title"`
	Content     string `json:"content" bson:"content"`
	Url         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"description"`
	Published   int64  `json:"published" bson:"published"`
}

func (m *Post) CollectionName() string {
	return "post"
}

func (m *Post) AutoRemovePost(ctx context.Context, before int64) error {
	_, err := DB(m).DeleteMany(ctx, bson.M{
		"published": bson.M{
			"$lte": before,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *Post) FindPostsByFeed(ctx context.Context, feedIDs []string, limit int64) ([]*Post, error) {
	if len(feedIDs) == 0 {
		return []*Post{}, nil
	}
	cur, err := DB(m).Find(ctx, bson.M{
		"fid": bson.M{
			"$in": feedIDs,
		},
	}, &options.FindOptions{
		Limit: &limit,
		Sort:  bson.M{"published": -1},
		Projection: bson.M{
			"_id":       1,
			"title":     1,
			"published": 1,
			"fid":       1,
			"url":       1,
		},
	})
	if err != nil {
		return nil, err
	}
	var posts = make([]*Post, 0)
	err = cur.All(ctx, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *Post) InsertMany(ctx context.Context, posts []interface{}) error {
	_, err := DB(m).InsertMany(ctx, posts)
	return err
}

func (m *Post) FindPostByID(ctx context.Context) error {
	return DB(m).FindOne(ctx, bson.M{
		"_id": m.ID,
	}).Decode(m)
}
