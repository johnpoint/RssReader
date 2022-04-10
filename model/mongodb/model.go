package mongodb

type Model interface {
	CollectionName() string
}
