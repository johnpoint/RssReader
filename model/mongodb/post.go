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
	Published   string `json:"published" bson:"published"`
}

func (m *Post) CollectionName() string {
	return "post"
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
		Sort:  bson.M{"Published": -1},
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
