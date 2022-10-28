package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	conn *mongo.Client

	teams         *mongo.Collection
	projects      *mongo.Collection
	mergeRequests *mongo.Collection
}

func New(conn *mongo.Client, databaseName string) (*Repository, error) {
	database := conn.Database(databaseName)

	return &Repository{
		conn:          conn,
		teams:         database.Collection("teams"),
		projects:      database.Collection("projects"),
		mergeRequests: database.Collection("merge_requests"),
	}, nil
}
