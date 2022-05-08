package mongodb

import (
	"RssReader/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type User struct {
	ID        string   `json:"_id" bson:"_id"`
	Mail      string   `json:"mail" bson:"mail"`
	Password  string   `json:"password" bson:"password"`
	CreatedAt int64    `json:"created_at" bson:"created_at"`
	SubFeeds  []string `json:"sub_feeds" bson:"sub_feeds"`
}

func (m *User) CollectionName() string {
	return "user"
}

func (m *User) InsertOne(ctx context.Context) error {
	m.CreatedAt = time.Now().UnixMilli()
	m.ID = utils.RandomString()
	_, err := DB(m).InsertOne(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *User) FindFeedByID(ctx context.Context) error {
	return DB(m).FindOne(ctx, bson.M{
		"_id": m.ID,
	}, &options.FindOneOptions{Projection: bson.M{
		"sub_feeds": 1,
		"_id":       1,
	}}).Decode(m)
}

func (m *User) FindInfoByID(ctx context.Context) error {
	return DB(m).FindOne(ctx, bson.M{
		"_id": m.ID,
	}, &options.FindOneOptions{Projection: bson.M{
		"mail": 1,
		"_id":  1,
	}}).Decode(m)
}

func (m *User) UpdatePassword(ctx context.Context, password string) error {
	_, err := DB(m).UpdateOne(ctx, bson.M{
		"_id": m.ID,
	}, bson.M{
		"$set": bson.M{
			"password": password,
		},
	})
	return err
}

func (m *User) FindOne(ctx context.Context, mail, password string) error {
	return DB(m).FindOne(ctx, bson.M{
		"mail":     mail,
		"password": password,
	}).Decode(m)
}

func (m *User) SubscribeFeeds(ctx context.Context, feedID ...string) error {
	if len(feedID) == 0 {
		return nil
	}
	_, err := DB(m).UpdateOne(ctx, bson.M{
		"_id": m.ID,
	}, bson.M{
		"$addToSet": bson.M{
			"sub_feeds": bson.M{
				"$each": feedID,
			},
		},
	})
	return err
}

func (m *User) UnSubscribeFeed(ctx context.Context, feedID string) error {
	_, err := DB(m).UpdateOne(ctx, bson.M{
		"_id": m.ID,
	}, bson.M{
		"$pull": bson.M{
			"sub_feeds": feedID,
		},
	})
	return err
}
