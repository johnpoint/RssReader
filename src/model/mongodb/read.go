package mongodb

import (
	"RssReader/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Read struct {
	ID       string `json:"_id" bson:"_id"`
	PId      string `json:"pid" bson:"pid"`
	UId      string `json:"uid" bson:"uid"`
	CreateAt int64  `json:"create_at" bson:"create_at"`
}

func (m *Read) CollectionName() string {
	return "read"
}

func (m *Read) FindReadListByUserId(ctx context.Context) ([]string, error) {
	var res []*Read
	var data = make([]string, 0)
	var limit int64 = 200
	cur, err := DB(m).Find(ctx, bson.M{
		"uid": m.UId,
	}, &options.FindOptions{
		Sort: bson.M{
			"create_at": 1,
		},
		Limit: &limit,
	})
	if err != nil {
		return data, err
	}
	err = cur.All(ctx, &res)
	if err != nil {
		return data, err
	}
	for _, v := range res {
		data = append(data, v.PId)
	}
	return data, nil
}

func (m *Read) MarkAsRead(ctx context.Context, read []*Read) error {
	var many []any
	for _, v := range read {
		v.ID = utils.Sha256(v.UId + v.PId)
		v.CreateAt = time.Now().UnixMilli()
		many = append(many, v)
	}

	_, err := DB(m).InsertMany(ctx, many)
	if err != nil {
		return err
	}
	return nil
}

func (m *Read) MarkAsUnRead(ctx context.Context, read []*Read) error {
	var delId []string
	for _, v := range read {
		delId = append(delId, utils.Sha256(v.UId+v.PId))
	}

	_, err := DB(m).DeleteMany(ctx, bson.M{
		"_id": bson.M{
			"$in": delId,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
