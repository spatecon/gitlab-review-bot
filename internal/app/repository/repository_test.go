//go:build mongodb

package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// dto is the default timeout used in the tests.
	dto = 5 * time.Second
)

func connHelper(t *testing.T, ctx context.Context) *mongo.Client {
	_ = godotenv.Load()

	host := os.Getenv("MONGO_HOST")

	if host == "" {
		t.Fatal("env MONGO_HOST is empty")
	}

	URI := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASS"),
		host,
		os.Getenv("MONGO_PORT"))

	conn, err := mongo.NewClient(options.Client().ApplyURI(URI).SetTimeout(dto))

	require.NoError(t, conn.Connect(ctx), "failed to connect to mongo")

	require.NoError(t, err, "failed to create mongo client")

	ctx, cancel := context.WithTimeout(ctx, dto)
	defer cancel()

	require.NoError(t, conn.Ping(ctx, nil), "failed to ping mongo client")

	return conn
}

// repositoryHelper is a helper function that returns a repository.
// It also creates a new database for a test run and deletes it after tests end.
func repositoryHelper(t *testing.T) *Repository {
	ctx := context.Background()
	conn := connHelper(t, ctx)

	dbName := "test-db-" + t.Name() + "-" + time.Now().Format("20060102150405")
	dbName = strings.Replace(dbName, "_", "-", -1)
	dbName = strings.ToLower(dbName)

	repository, err := New(ctx, conn, dbName)
	require.NoError(t, err, "failed to create repository")

	t.Cleanup(func() {
		require.NoError(t, conn.Database(dbName).Drop(ctx), "failed to drop test database in cleanup")
		require.NoError(t, conn.Disconnect(ctx), "failed to disconnect mongo client in cleanup")
	})

	return repository
}
