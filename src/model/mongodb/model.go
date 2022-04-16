package mongodb

import (
	"RssReader/dao/mongoDao"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model interface {
	CollectionName() string
}

func DB(m Model) *mongo.Collection {
	return mongoDao.Client(m.CollectionName())
}
