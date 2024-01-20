package store

import "go.mongodb.org/mongo-driver/mongo"

type MongoDBStore struct {
	*mongo.Collection
}
